package sqlite

import (
	"context"
	"database/sql"
)

// schemaVersion 是当前内置 SQLite schema 版本。
const schemaVersion = 1

// migrations 按顺序声明 Phase 3 Product Facts 所需表结构。
var migrations = []string{
	`CREATE TABLE IF NOT EXISTS schema_versions (
		version INTEGER PRIMARY KEY,
		applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`,
	`CREATE TABLE IF NOT EXISTS runs (
		run_id TEXT PRIMARY KEY,
		workspace_id TEXT NOT NULL,
		status TEXT NOT NULL,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		model_label TEXT NOT NULL DEFAULT '',
		safe_error TEXT NOT NULL DEFAULT ''
	);`,
	`CREATE INDEX IF NOT EXISTS idx_runs_workspace_updated ON runs(workspace_id, updated_at);`,
	`CREATE TABLE IF NOT EXISTS turns (
		turn_id TEXT PRIMARY KEY,
		run_id TEXT NOT NULL,
		role TEXT NOT NULL,
		content TEXT NOT NULL,
		sequence INTEGER NOT NULL,
		created_at TEXT NOT NULL,
		FOREIGN KEY(run_id) REFERENCES runs(run_id)
	);`,
	`CREATE UNIQUE INDEX IF NOT EXISTS idx_turns_run_sequence ON turns(run_id, sequence);`,
	`CREATE TABLE IF NOT EXISTS tool_calls (
		tool_call_id TEXT PRIMARY KEY,
		run_id TEXT NOT NULL,
		tool_id TEXT NOT NULL,
		display_name TEXT NOT NULL,
		status TEXT NOT NULL,
		args_hash TEXT NOT NULL DEFAULT '',
		args_preview TEXT NOT NULL DEFAULT '',
		created_at TEXT NOT NULL,
		ended_at TEXT NOT NULL DEFAULT '',
		FOREIGN KEY(run_id) REFERENCES runs(run_id)
	);`,
	`CREATE TABLE IF NOT EXISTS tool_results (
		result_id TEXT PRIMARY KEY,
		tool_call_id TEXT NOT NULL,
		status TEXT NOT NULL,
		structured_schema_version TEXT NOT NULL,
		result_ref TEXT NOT NULL,
		safe_summary TEXT NOT NULL,
		FOREIGN KEY(tool_call_id) REFERENCES tool_calls(tool_call_id)
	);`,
	`CREATE TABLE IF NOT EXISTS pending_interactions (
		pending_id TEXT PRIMARY KEY,
		run_id TEXT NOT NULL,
		kind TEXT NOT NULL,
		status TEXT NOT NULL,
		resume_ref TEXT NOT NULL UNIQUE,
		checkpoint_ref TEXT NOT NULL,
		question TEXT NOT NULL DEFAULT '',
		risk_summary TEXT NOT NULL DEFAULT '',
		input_mode TEXT NOT NULL DEFAULT '',
		candidates_json TEXT NOT NULL DEFAULT '[]',
		expires_at TEXT NOT NULL DEFAULT '',
		FOREIGN KEY(run_id) REFERENCES runs(run_id)
	);`,
	`CREATE TABLE IF NOT EXISTS audit_events (
		audit_id TEXT PRIMARY KEY,
		run_id TEXT NOT NULL,
		event_type TEXT NOT NULL,
		safe_summary TEXT NOT NULL,
		actor TEXT NOT NULL,
		created_at TEXT NOT NULL,
		FOREIGN KEY(run_id) REFERENCES runs(run_id)
	);`,
	`CREATE TABLE IF NOT EXISTS context_snapshots (
		snapshot_id TEXT PRIMARY KEY,
		run_id TEXT NOT NULL,
		safe_summary TEXT NOT NULL,
		created_at TEXT NOT NULL,
		FOREIGN KEY(run_id) REFERENCES runs(run_id)
	);`,
	`CREATE TABLE IF NOT EXISTS idempotency_records (
		scope TEXT NOT NULL,
		key TEXT NOT NULL,
		run_id TEXT NOT NULL,
		resource_ref TEXT NOT NULL DEFAULT '',
		status TEXT NOT NULL,
		created_at TEXT NOT NULL,
		PRIMARY KEY(scope, key)
	);`,
}

// migrate 执行内置 migration，并记录当前 schema 版本。
func migrate(ctx context.Context, db *sql.DB) error {
	for _, stmt := range migrations {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	if err := addColumnIfMissing(ctx, db, "pending_interactions", "input_mode", "TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if err := addColumnIfMissing(ctx, db, "pending_interactions", "candidates_json", "TEXT NOT NULL DEFAULT '[]'"); err != nil {
		return err
	}
	_, err := db.ExecContext(ctx, `INSERT OR IGNORE INTO schema_versions(version) VALUES (?)`, schemaVersion)
	return err
}

// addColumnIfMissing 兼容已有本地 SQLite 文件，避免 Product Facts 新字段只在全新库中存在。
func addColumnIfMissing(ctx context.Context, db *sql.DB, tableName string, columnName string, columnDDL string) error {
	rows, err := db.QueryContext(ctx, `PRAGMA table_info(`+tableName+`)`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var name string
		var columnType string
		var notNull int
		var defaultValue sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &columnType, &notNull, &defaultValue, &pk); err != nil {
			return err
		}
		if name == columnName {
			return rows.Err()
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, `ALTER TABLE `+tableName+` ADD COLUMN `+columnName+` `+columnDDL)
	return err
}
