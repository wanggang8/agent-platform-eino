package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"agent-platform-eino/internal/einoapp/facts"

	_ "modernc.org/sqlite"
)

// Repository 是 Product Facts 的 SQLite 持久化实现。
// 该层只处理事实存取和事务，不依赖 HTTP、provider、LLM 或 Eino 事件。
type Repository struct {
	db *sql.DB
}

var _ facts.Repository = (*Repository)(nil)

// Open 打开 SQLite 数据库并执行项目内 migration。
func Open(ctx context.Context, path string) (*Repository, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.ExecContext(ctx, `PRAGMA foreign_keys = ON`); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := migrate(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Repository{db: db}, nil
}

// Close 关闭底层数据库连接。
func (repo *Repository) Close() error {
	return repo.db.Close()
}

// CreateRun 创建 run 根事实，调用方负责提供稳定 id 和时间。
func (repo *Repository) CreateRun(ctx context.Context, run facts.Run) error {
	if facts.ContainsUnsafeMaterial(run.SafeError) {
		return facts.ErrUnsafeFactMaterial
	}
	_, err := repo.db.ExecContext(ctx, `
		INSERT INTO runs(run_id, workspace_id, status, created_at, updated_at, model_label, safe_error)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		run.RunID,
		run.WorkspaceID,
		string(run.Status),
		formatTime(run.CreatedAt),
		formatTime(run.UpdatedAt),
		run.ModelLabel,
		run.SafeError,
	)
	return err
}

// GetRun 按 run_id 读取 run 根事实。
func (repo *Repository) GetRun(ctx context.Context, runID string) (facts.Run, error) {
	row := repo.db.QueryRowContext(ctx, `
		SELECT run_id, workspace_id, status, created_at, updated_at, model_label, safe_error
		FROM runs
		WHERE run_id = ?`, runID)
	return scanRun(row)
}

// LatestRun 读取 workspace 内最近更新的 run，用于 Workbench 默认视图。
func (repo *Repository) LatestRun(ctx context.Context, workspaceID string) (facts.Run, error) {
	row := repo.db.QueryRowContext(ctx, `
		SELECT run_id, workspace_id, status, created_at, updated_at, model_label, safe_error
		FROM runs
		WHERE workspace_id = ?
		ORDER BY updated_at DESC, rowid DESC
		LIMIT 1`, workspaceID)
	return scanRun(row)
}

// GetSnapshot 按 run 读取完整 Product Facts，保证 Workbench、Action API、replay 同源投影。
func (repo *Repository) GetSnapshot(ctx context.Context, runID string) (facts.Snapshot, error) {
	run, err := repo.GetRun(ctx, runID)
	if err != nil {
		return facts.Snapshot{}, err
	}
	snapshot := facts.Snapshot{Run: run}

	if snapshot.Turns, err = repo.listTurns(ctx, runID); err != nil {
		return facts.Snapshot{}, err
	}
	if snapshot.ToolCalls, err = repo.listToolCalls(ctx, runID); err != nil {
		return facts.Snapshot{}, err
	}
	if snapshot.ToolResults, err = repo.listToolResults(ctx, runID); err != nil {
		return facts.Snapshot{}, err
	}
	if snapshot.PendingInteractions, err = repo.listPendingInteractions(ctx, runID); err != nil {
		return facts.Snapshot{}, err
	}
	if snapshot.AuditEvents, err = repo.listAuditEvents(ctx, runID); err != nil {
		return facts.Snapshot{}, err
	}
	if snapshot.ContextSnapshots, err = repo.listContextSnapshots(ctx, runID); err != nil {
		return facts.Snapshot{}, err
	}
	return snapshot, nil
}

// UpdateRunStatus 迁移 run 生命周期，并拒绝把明显不安全材料写入 safe_error。
func (repo *Repository) UpdateRunStatus(ctx context.Context, runID string, status facts.RunStatus, safeError string, updatedAt time.Time) error {
	if facts.ContainsUnsafeMaterial(safeError) {
		return facts.ErrUnsafeFactMaterial
	}
	result, err := repo.db.ExecContext(ctx, `
		UPDATE runs
		SET status = ?, safe_error = ?, updated_at = ?
		WHERE run_id = ?`,
		string(status),
		safeError,
		formatTime(updatedAt),
		runID,
	)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return facts.ErrNotFound
	}
	return nil
}

// RecordIdempotency 记录幂等请求；相同 key 再次提交返回原记录，不覆盖事实。
func (repo *Repository) RecordIdempotency(ctx context.Context, record facts.IdempotencyRecord) (facts.IdempotencyRecord, bool, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return facts.IdempotencyRecord{}, false, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	existing, err := getIdempotencyRecord(ctx, tx, record.Scope, record.Key)
	if err == nil {
		if existing.ResourceRef != "" && record.ResourceRef != "" && existing.ResourceRef != record.ResourceRef {
			return facts.IdempotencyRecord{}, false, facts.ErrIdempotencyConflict
		}
		if err := tx.Commit(); err != nil {
			return facts.IdempotencyRecord{}, false, err
		}
		return existing, true, nil
	}
	if !errors.Is(err, facts.ErrNotFound) {
		return facts.IdempotencyRecord{}, false, err
	}

	if err := insertIdempotencyRecord(ctx, tx, record); err != nil {
		return facts.IdempotencyRecord{}, false, err
	}
	if err := tx.Commit(); err != nil {
		return facts.IdempotencyRecord{}, false, err
	}
	return record, false, nil
}

// AppendTurn 追加消息事实；content 必须已经是安全投影后的文本。
func (repo *Repository) AppendTurn(ctx context.Context, turn facts.Turn) error {
	if facts.ContainsUnsafeMaterial(turn.Content) {
		return facts.ErrUnsafeFactMaterial
	}
	_, err := repo.db.ExecContext(ctx, `
		INSERT INTO turns(turn_id, run_id, role, content, sequence, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		turn.TurnID,
		turn.RunID,
		string(turn.Role),
		turn.Content,
		turn.Sequence,
		formatTime(turn.CreatedAt),
	)
	return err
}

// AppendToolCall 追加工具调用事实；args preview 只能是 allowlist 后的安全摘要。
func (repo *Repository) AppendToolCall(ctx context.Context, call facts.ToolCall) error {
	if facts.ContainsUnsafeMaterial(call.ArgsPreview) {
		return facts.ErrUnsafeFactMaterial
	}
	_, err := repo.db.ExecContext(ctx, `
		INSERT INTO tool_calls(tool_call_id, run_id, tool_id, display_name, status, args_hash, args_preview, created_at, ended_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		call.ToolCallID,
		call.RunID,
		call.ToolID,
		call.DisplayName,
		string(call.Status),
		call.ArgsHash,
		call.ArgsPreview,
		formatTime(call.CreatedAt),
		formatTime(call.EndedAt),
	)
	return err
}

// AppendToolResult 追加工具结果事实，只保存 StructuredResult 引用和安全摘要。
func (repo *Repository) AppendToolResult(ctx context.Context, result facts.ToolResult) error {
	if facts.UnsafeStructuredResultRef(result.StructuredResult) {
		return facts.ErrUnsafeFactMaterial
	}
	_, err := repo.db.ExecContext(ctx, `
		INSERT INTO tool_results(result_id, tool_call_id, status, structured_schema_version, result_ref, safe_summary)
		VALUES (?, ?, ?, ?, ?, ?)`,
		result.ResultID,
		result.ToolCallID,
		string(result.Status),
		result.StructuredResult.SchemaVersion,
		result.StructuredResult.ResultRef,
		result.StructuredResult.SafeSummary,
	)
	return err
}

// AppendPendingInteraction 追加审批/澄清等待事实；resume/checkpoint 都必须是安全引用。
func (repo *Repository) AppendPendingInteraction(ctx context.Context, pending facts.PendingInteraction) error {
	if facts.ContainsUnsafeMaterial(pending.ResumeRef) ||
		facts.ContainsUnsafeMaterial(pending.CheckpointRef) ||
		facts.ContainsUnsafeMaterial(pending.Question) ||
		facts.ContainsUnsafeMaterial(pending.RiskSummary) {
		return facts.ErrUnsafeFactMaterial
	}
	_, err := repo.db.ExecContext(ctx, `
		INSERT INTO pending_interactions(pending_id, run_id, kind, status, resume_ref, checkpoint_ref, question, risk_summary, expires_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		pending.PendingID,
		pending.RunID,
		string(pending.Kind),
		string(pending.Status),
		pending.ResumeRef,
		pending.CheckpointRef,
		pending.Question,
		pending.RiskSummary,
		formatTime(pending.ExpiresAt),
	)
	return err
}

// ConsumeResumeRef 消费一次性 resume 引用；重复消费返回已消费 pending 和幂等错误。
func (repo *Repository) ConsumeResumeRef(ctx context.Context, resumeRef string, submittedStatus facts.PendingStatus) (facts.PendingInteraction, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return facts.PendingInteraction{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	pending, err := consumeResumeRefInTx(ctx, tx, resumeRef, submittedStatus)
	if err != nil {
		if errors.Is(err, facts.ErrResumeAlreadyConsumed) {
			if commitErr := tx.Commit(); commitErr != nil {
				return facts.PendingInteraction{}, commitErr
			}
			return pending, err
		}
		return facts.PendingInteraction{}, err
	}
	if err := tx.Commit(); err != nil {
		return facts.PendingInteraction{}, err
	}
	return pending, nil
}

// ConsumeResumeRefWithIdempotency 在同一事务内完成 resume 消费和幂等记录写入。
// ResourceRef 必须与 resumeRef 一致，防止同一幂等键误绑其他 pending。
func (repo *Repository) ConsumeResumeRefWithIdempotency(ctx context.Context, resumeRef string, submittedStatus facts.PendingStatus, record facts.IdempotencyRecord) (facts.PendingInteraction, facts.IdempotencyRecord, bool, error) {
	if record.ResourceRef != "" && record.ResourceRef != resumeRef {
		return facts.PendingInteraction{}, facts.IdempotencyRecord{}, false, facts.ErrIdempotencyConflict
	}
	if record.ResourceRef == "" {
		record.ResourceRef = resumeRef
	}

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return facts.PendingInteraction{}, facts.IdempotencyRecord{}, false, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	existing, err := getIdempotencyRecord(ctx, tx, record.Scope, record.Key)
	if err == nil {
		if existing.ResourceRef != "" && record.ResourceRef != "" && existing.ResourceRef != record.ResourceRef {
			return facts.PendingInteraction{}, facts.IdempotencyRecord{}, false, facts.ErrIdempotencyConflict
		}
		pending, pendingErr := getPendingForResume(ctx, tx, resumeRef)
		if pendingErr != nil {
			return facts.PendingInteraction{}, facts.IdempotencyRecord{}, false, pendingErr
		}
		if err := tx.Commit(); err != nil {
			return facts.PendingInteraction{}, facts.IdempotencyRecord{}, false, err
		}
		return pending, existing, true, nil
	}
	if !errors.Is(err, facts.ErrNotFound) {
		return facts.PendingInteraction{}, facts.IdempotencyRecord{}, false, err
	}

	pending, err := consumeResumeRefInTx(ctx, tx, resumeRef, submittedStatus)
	if err != nil {
		if errors.Is(err, facts.ErrResumeAlreadyConsumed) {
			if commitErr := tx.Commit(); commitErr != nil {
				return facts.PendingInteraction{}, facts.IdempotencyRecord{}, false, commitErr
			}
			return pending, facts.IdempotencyRecord{}, false, err
		}
		return facts.PendingInteraction{}, facts.IdempotencyRecord{}, false, err
	}
	if err := insertIdempotencyRecord(ctx, tx, record); err != nil {
		return facts.PendingInteraction{}, facts.IdempotencyRecord{}, false, err
	}
	if err := tx.Commit(); err != nil {
		return facts.PendingInteraction{}, facts.IdempotencyRecord{}, false, err
	}
	return pending, record, false, nil
}

// AppendAuditEvent 追加审计事件；只允许安全摘要进入 audit。
func (repo *Repository) AppendAuditEvent(ctx context.Context, event facts.AuditEvent) error {
	if facts.ContainsUnsafeMaterial(event.SafeSummary) {
		return facts.ErrUnsafeFactMaterial
	}
	_, err := repo.db.ExecContext(ctx, `
		INSERT INTO audit_events(audit_id, run_id, event_type, safe_summary, actor, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		event.AuditID,
		event.RunID,
		event.EventType,
		event.SafeSummary,
		event.Actor,
		formatTime(event.CreatedAt),
	)
	return err
}

// SaveContextSnapshot 保存模型调用前的安全上下文摘要。
func (repo *Repository) SaveContextSnapshot(ctx context.Context, snapshot facts.ContextSnapshot) error {
	if facts.ContainsUnsafeMaterial(snapshot.SafeSummary) {
		return facts.ErrUnsafeFactMaterial
	}
	_, err := repo.db.ExecContext(ctx, `
		INSERT INTO context_snapshots(snapshot_id, run_id, safe_summary, created_at)
		VALUES (?, ?, ?, ?)`,
		snapshot.SnapshotID,
		snapshot.RunID,
		snapshot.SafeSummary,
		formatTime(snapshot.CreatedAt),
	)
	return err
}

type rowScanner interface {
	Scan(dest ...any) error
}

// scanRun 将数据库行转换为 Run，并把 sql.ErrNoRows 归一为 facts.ErrNotFound。
func scanRun(row rowScanner) (facts.Run, error) {
	var run facts.Run
	var status string
	var createdAt string
	var updatedAt string
	err := row.Scan(
		&run.RunID,
		&run.WorkspaceID,
		&status,
		&createdAt,
		&updatedAt,
		&run.ModelLabel,
		&run.SafeError,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return facts.Run{}, facts.ErrNotFound
	}
	if err != nil {
		return facts.Run{}, err
	}
	run.Status = facts.RunStatus(status)
	run.CreatedAt, err = parseTime(createdAt)
	if err != nil {
		return facts.Run{}, err
	}
	run.UpdatedAt, err = parseTime(updatedAt)
	if err != nil {
		return facts.Run{}, err
	}
	return run, nil
}

// getPendingForResume 在事务内读取 resume_ref 对应的 pending。
func getPendingForResume(ctx context.Context, tx *sql.Tx, resumeRef string) (facts.PendingInteraction, error) {
	row := tx.QueryRowContext(ctx, `
		SELECT pending_id, run_id, kind, status, resume_ref, checkpoint_ref, question, risk_summary, expires_at
		FROM pending_interactions
		WHERE resume_ref = ?`, resumeRef)
	var pending facts.PendingInteraction
	var kind string
	var status string
	var expiresAt string
	err := row.Scan(
		&pending.PendingID,
		&pending.RunID,
		&kind,
		&status,
		&pending.ResumeRef,
		&pending.CheckpointRef,
		&pending.Question,
		&pending.RiskSummary,
		&expiresAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return facts.PendingInteraction{}, facts.ErrNotFound
	}
	if err != nil {
		return facts.PendingInteraction{}, err
	}
	pending.Kind = facts.PendingKind(kind)
	pending.Status = facts.PendingStatus(status)
	pending.ExpiresAt, err = parseTime(expiresAt)
	if err != nil {
		return facts.PendingInteraction{}, err
	}
	return pending, nil
}

// consumeResumeRefInTx 在已有事务中迁移 pending 为 consumed。
func consumeResumeRefInTx(ctx context.Context, tx *sql.Tx, resumeRef string, submittedStatus facts.PendingStatus) (facts.PendingInteraction, error) {
	pending, err := getPendingForResume(ctx, tx, resumeRef)
	if err != nil {
		return facts.PendingInteraction{}, err
	}
	if pending.Status == facts.PendingStatusConsumed {
		return pending, facts.ErrResumeAlreadyConsumed
	}
	if pending.Status != facts.PendingStatusWaiting && pending.Status != submittedStatus {
		return pending, facts.ErrResumeAlreadyConsumed
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE pending_interactions
		SET status = ?
		WHERE resume_ref = ? AND status IN (?, ?)`,
		string(facts.PendingStatusConsumed),
		resumeRef,
		string(facts.PendingStatusWaiting),
		string(submittedStatus),
	); err != nil {
		return facts.PendingInteraction{}, err
	}
	pending.Status = facts.PendingStatusConsumed
	return pending, nil
}

// getIdempotencyRecord 在事务内读取幂等记录。
func getIdempotencyRecord(ctx context.Context, tx *sql.Tx, scope facts.IdempotencyScope, key string) (facts.IdempotencyRecord, error) {
	row := tx.QueryRowContext(ctx, `
		SELECT scope, key, run_id, resource_ref, status, created_at
		FROM idempotency_records
		WHERE scope = ? AND key = ?`,
		string(scope),
		key,
	)
	return scanIdempotencyRecord(row)
}

// insertIdempotencyRecord 在事务内写入幂等记录。
func insertIdempotencyRecord(ctx context.Context, tx *sql.Tx, record facts.IdempotencyRecord) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO idempotency_records(scope, key, run_id, resource_ref, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		string(record.Scope),
		record.Key,
		record.RunID,
		record.ResourceRef,
		record.Status,
		formatTime(record.CreatedAt),
	)
	return err
}

func (repo *Repository) listTurns(ctx context.Context, runID string) ([]facts.Turn, error) {
	rows, err := repo.db.QueryContext(ctx, `
		SELECT turn_id, run_id, role, content, sequence, created_at
		FROM turns
		WHERE run_id = ?
		ORDER BY sequence ASC`, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var turns []facts.Turn
	for rows.Next() {
		var turn facts.Turn
		var role string
		var createdAt string
		if err := rows.Scan(&turn.TurnID, &turn.RunID, &role, &turn.Content, &turn.Sequence, &createdAt); err != nil {
			return nil, err
		}
		turn.Role = facts.TurnRole(role)
		parsed, err := parseTime(createdAt)
		if err != nil {
			return nil, err
		}
		turn.CreatedAt = parsed
		turns = append(turns, turn)
	}
	return turns, rows.Err()
}

func (repo *Repository) listToolCalls(ctx context.Context, runID string) ([]facts.ToolCall, error) {
	rows, err := repo.db.QueryContext(ctx, `
		SELECT tool_call_id, run_id, tool_id, display_name, status, args_hash, args_preview, created_at, ended_at
		FROM tool_calls
		WHERE run_id = ?
		ORDER BY created_at ASC, rowid ASC`, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calls []facts.ToolCall
	for rows.Next() {
		var call facts.ToolCall
		var status string
		var createdAt string
		var endedAt string
		if err := rows.Scan(&call.ToolCallID, &call.RunID, &call.ToolID, &call.DisplayName, &status, &call.ArgsHash, &call.ArgsPreview, &createdAt, &endedAt); err != nil {
			return nil, err
		}
		call.Status = facts.ToolCallStatus(status)
		var err error
		call.CreatedAt, err = parseTime(createdAt)
		if err != nil {
			return nil, err
		}
		call.EndedAt, err = parseTime(endedAt)
		if err != nil {
			return nil, err
		}
		calls = append(calls, call)
	}
	return calls, rows.Err()
}

func (repo *Repository) listToolResults(ctx context.Context, runID string) ([]facts.ToolResult, error) {
	rows, err := repo.db.QueryContext(ctx, `
		SELECT tr.result_id, tr.tool_call_id, tr.status, tr.structured_schema_version, tr.result_ref, tr.safe_summary
		FROM tool_results tr
		INNER JOIN tool_calls tc ON tc.tool_call_id = tr.tool_call_id
		WHERE tc.run_id = ?
		ORDER BY tc.created_at ASC, tr.rowid ASC`, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []facts.ToolResult
	for rows.Next() {
		var result facts.ToolResult
		var status string
		if err := rows.Scan(&result.ResultID, &result.ToolCallID, &status, &result.StructuredResult.SchemaVersion, &result.StructuredResult.ResultRef, &result.StructuredResult.SafeSummary); err != nil {
			return nil, err
		}
		result.Status = facts.ToolResultStatus(status)
		results = append(results, result)
	}
	return results, rows.Err()
}

func (repo *Repository) listPendingInteractions(ctx context.Context, runID string) ([]facts.PendingInteraction, error) {
	rows, err := repo.db.QueryContext(ctx, `
		SELECT pending_id, run_id, kind, status, resume_ref, checkpoint_ref, question, risk_summary, expires_at
		FROM pending_interactions
		WHERE run_id = ?
		ORDER BY rowid ASC`, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pending []facts.PendingInteraction
	for rows.Next() {
		item, err := scanPending(rows)
		if err != nil {
			return nil, err
		}
		pending = append(pending, item)
	}
	return pending, rows.Err()
}

func (repo *Repository) listAuditEvents(ctx context.Context, runID string) ([]facts.AuditEvent, error) {
	rows, err := repo.db.QueryContext(ctx, `
		SELECT audit_id, run_id, event_type, safe_summary, actor, created_at
		FROM audit_events
		WHERE run_id = ?
		ORDER BY created_at ASC, rowid ASC`, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []facts.AuditEvent
	for rows.Next() {
		var event facts.AuditEvent
		var createdAt string
		if err := rows.Scan(&event.AuditID, &event.RunID, &event.EventType, &event.SafeSummary, &event.Actor, &createdAt); err != nil {
			return nil, err
		}
		parsed, err := parseTime(createdAt)
		if err != nil {
			return nil, err
		}
		event.CreatedAt = parsed
		events = append(events, event)
	}
	return events, rows.Err()
}

func (repo *Repository) listContextSnapshots(ctx context.Context, runID string) ([]facts.ContextSnapshot, error) {
	rows, err := repo.db.QueryContext(ctx, `
		SELECT snapshot_id, run_id, safe_summary, created_at
		FROM context_snapshots
		WHERE run_id = ?
		ORDER BY created_at ASC, rowid ASC`, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []facts.ContextSnapshot
	for rows.Next() {
		var snapshot facts.ContextSnapshot
		var createdAt string
		if err := rows.Scan(&snapshot.SnapshotID, &snapshot.RunID, &snapshot.SafeSummary, &createdAt); err != nil {
			return nil, err
		}
		parsed, err := parseTime(createdAt)
		if err != nil {
			return nil, err
		}
		snapshot.CreatedAt = parsed
		snapshots = append(snapshots, snapshot)
	}
	return snapshots, rows.Err()
}

func scanPending(row rowScanner) (facts.PendingInteraction, error) {
	var pending facts.PendingInteraction
	var kind string
	var status string
	var expiresAt string
	if err := row.Scan(&pending.PendingID, &pending.RunID, &kind, &status, &pending.ResumeRef, &pending.CheckpointRef, &pending.Question, &pending.RiskSummary, &expiresAt); err != nil {
		return facts.PendingInteraction{}, err
	}
	pending.Kind = facts.PendingKind(kind)
	pending.Status = facts.PendingStatus(status)
	parsed, err := parseTime(expiresAt)
	if err != nil {
		return facts.PendingInteraction{}, err
	}
	pending.ExpiresAt = parsed
	return pending, nil
}

func scanIdempotencyRecord(row rowScanner) (facts.IdempotencyRecord, error) {
	var record facts.IdempotencyRecord
	var scope string
	var createdAt string
	err := row.Scan(&scope, &record.Key, &record.RunID, &record.ResourceRef, &record.Status, &createdAt)
	if errors.Is(err, sql.ErrNoRows) {
		return facts.IdempotencyRecord{}, facts.ErrNotFound
	}
	if err != nil {
		return facts.IdempotencyRecord{}, err
	}
	record.Scope = facts.IdempotencyScope(scope)
	record.CreatedAt, err = parseTime(createdAt)
	if err != nil {
		return facts.IdempotencyRecord{}, err
	}
	return record, nil
}

func formatTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}

func parseTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, err
	}
	return parsed, nil
}
