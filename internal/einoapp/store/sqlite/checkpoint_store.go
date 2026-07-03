package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"agent-platform-eino/internal/einoapp/facts"

	"github.com/cloudwego/eino/compose"
)

// CheckpointBinding 绑定 Product Facts 安全 checkpoint_ref 与内部 Eino checkpoint id。
// CheckpointID 是内部恢复键，不能进入 Workbench、Action API、audit 或 replay。
type CheckpointBinding struct {
	CheckpointRef string
	CheckpointID  string
	RunID         string
	PendingID     string
	CreatedAt     time.Time
}

// CheckpointStore 实现 Eino CheckPointStore，并提供安全 checkpoint_ref 解析。
type CheckpointStore struct {
	db *sql.DB
}

var _ compose.CheckPointStore = (*CheckpointStore)(nil)

// OpenCheckpointStore 打开 checkpoint store，并创建内部 checkpoint 表。
func OpenCheckpointStore(ctx context.Context, path string) (*CheckpointStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.ExecContext(ctx, `PRAGMA foreign_keys = ON`); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := migrateCheckpointStore(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &CheckpointStore{db: db}, nil
}

// Close 关闭 checkpoint store 的数据库连接。
func (store *CheckpointStore) Close() error {
	return store.db.Close()
}

// Set 写入 Eino 序列化后的 checkpoint bytes。
func (store *CheckpointStore) Set(ctx context.Context, checkpointID string, checkpoint []byte) error {
	if strings.TrimSpace(checkpointID) == "" {
		return facts.ErrUnsafeFactMaterial
	}
	now := formatTime(time.Now().UTC())
	_, err := store.db.ExecContext(ctx, `
		INSERT INTO checkpoint_values(checkpoint_id, checkpoint_blob, updated_at)
		VALUES (?, ?, ?)
		ON CONFLICT(checkpoint_id) DO UPDATE SET checkpoint_blob = excluded.checkpoint_blob, updated_at = excluded.updated_at`,
		checkpointID,
		checkpoint,
		now,
	)
	return err
}

// Get 按内部 Eino checkpoint id 读取 checkpoint bytes。
func (store *CheckpointStore) Get(ctx context.Context, checkpointID string) ([]byte, bool, error) {
	row := store.db.QueryRowContext(ctx, `
		SELECT checkpoint_blob
		FROM checkpoint_values
		WHERE checkpoint_id = ?`, checkpointID)
	var payload []byte
	err := row.Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	return append([]byte(nil), payload...), true, nil
}

// BindCheckpointRef 写入安全 checkpoint_ref 到内部 Eino checkpoint id 的映射。
func (store *CheckpointStore) BindCheckpointRef(ctx context.Context, binding CheckpointBinding) error {
	if strings.TrimSpace(binding.CheckpointRef) == "" ||
		strings.TrimSpace(binding.CheckpointID) == "" ||
		strings.TrimSpace(binding.RunID) == "" ||
		strings.TrimSpace(binding.PendingID) == "" ||
		!facts.SafeCheckpointRef(binding.CheckpointRef) ||
		binding.CheckpointRef == binding.CheckpointID ||
		facts.ContainsUnsafeMaterial(binding.RunID) ||
		facts.ContainsUnsafeMaterial(binding.PendingID) {
		return facts.ErrUnsafeFactMaterial
	}
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var existingCheckpointID string
	var existingRunID string
	var existingPendingID string
	err = tx.QueryRowContext(ctx, `
		SELECT checkpoint_id, run_id, pending_id
		FROM checkpoint_refs
		WHERE checkpoint_ref = ?`, binding.CheckpointRef).Scan(&existingCheckpointID, &existingRunID, &existingPendingID)
	if err == nil {
		if existingCheckpointID == binding.CheckpointID && existingRunID == binding.RunID && existingPendingID == binding.PendingID {
			return tx.Commit()
		}
		return facts.ErrIdempotencyConflict
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	createdAt := binding.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO checkpoint_refs(checkpoint_ref, checkpoint_id, run_id, pending_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		binding.CheckpointRef,
		binding.CheckpointID,
		binding.RunID,
		binding.PendingID,
		formatTime(createdAt),
		formatTime(time.Now().UTC()),
	); err != nil {
		return err
	}
	return tx.Commit()
}

// ResolveCheckpointID 通过安全 checkpoint_ref 和 run/pending 绑定找到内部 checkpoint id。
func (store *CheckpointStore) ResolveCheckpointID(ctx context.Context, checkpointRef string, runID string, pendingID string) (string, bool, error) {
	if !facts.SafeCheckpointRef(checkpointRef) ||
		strings.TrimSpace(runID) == "" ||
		strings.TrimSpace(pendingID) == "" ||
		facts.ContainsUnsafeMaterial(runID) ||
		facts.ContainsUnsafeMaterial(pendingID) {
		return "", false, facts.ErrUnsafeFactMaterial
	}
	row := store.db.QueryRowContext(ctx, `
		SELECT cr.checkpoint_id
		FROM checkpoint_refs cr
		JOIN checkpoint_values cv ON cv.checkpoint_id = cr.checkpoint_id
		WHERE cr.checkpoint_ref = ? AND cr.run_id = ? AND cr.pending_id = ?`,
		checkpointRef,
		runID,
		pendingID,
	)
	var checkpointID string
	err := row.Scan(&checkpointID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return checkpointID, true, nil
}

func migrateCheckpointStore(ctx context.Context, db *sql.DB) error {
	for _, stmt := range []string{
		`CREATE TABLE IF NOT EXISTS checkpoint_values (
			checkpoint_id TEXT PRIMARY KEY,
			checkpoint_blob BLOB NOT NULL,
			updated_at TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS checkpoint_refs (
			checkpoint_ref TEXT PRIMARY KEY,
			checkpoint_id TEXT NOT NULL,
			run_id TEXT NOT NULL,
			pending_id TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY(checkpoint_id) REFERENCES checkpoint_values(checkpoint_id)
		);`,
		`CREATE INDEX IF NOT EXISTS idx_checkpoint_refs_run_pending ON checkpoint_refs(run_id, pending_id);`,
	} {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}
