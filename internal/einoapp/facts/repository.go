package facts

import (
	"context"
	"time"
)

type Repository interface {
	// CreateRun 写入 run 根事实，作为后续 turn、tool、audit 的归属锚点。
	CreateRun(ctx context.Context, run Run) error
	// GetRun 按 run_id 读取根事实；找不到时返回 ErrNotFound。
	GetRun(ctx context.Context, runID string) (Run, error)
	// LatestRun 返回 workspace 下最新 run，用于 Workbench current view。
	LatestRun(ctx context.Context, workspaceID string) (Run, error)
	// GetSnapshot 聚合同一 run 的安全事实，不暴露 provider raw payload。
	GetSnapshot(ctx context.Context, runID string) (Snapshot, error)
	// UpdateRunStatus 只更新 run 生命周期状态和脱敏错误摘要。
	UpdateRunStatus(ctx context.Context, runID string, status RunStatus, safeError string, updatedAt time.Time) error
	// ApplyLifecycleTransition 原子迁移 run/pending/tool/audit 生命周期事实。
	ApplyLifecycleTransition(ctx context.Context, transition LifecycleTransition) error
	// CreateRetryRun 原子写入 retry 幂等记录、新 run 和 retry audit。
	CreateRetryRun(ctx context.Context, transition RetryRunTransition) (IdempotencyRecord, bool, error)
	// RecordIdempotency 保存请求幂等记录，返回是否命中已有记录。
	RecordIdempotency(ctx context.Context, record IdempotencyRecord) (IdempotencyRecord, bool, error)
	// AppendTurn 追加对话 turn，sequence 由调用方保证单 run 内稳定递增。
	AppendTurn(ctx context.Context, turn Turn) error
	// AppendToolCall 追加工具调用事实；工具选择来源必须是 capability registry。
	AppendToolCall(ctx context.Context, call ToolCall) error
	// AppendToolResult 追加 StructuredResult 引用，禁止写入 provider raw payload。
	AppendToolResult(ctx context.Context, result ToolResult) error
	// AppendPendingInteraction 记录 approval/clarification 等待态和安全 resume ref。
	AppendPendingInteraction(ctx context.Context, pending PendingInteraction) error
	// GetPendingByResumeRef 只读取等待交互，不消费 resume_ref。
	GetPendingByResumeRef(ctx context.Context, resumeRef string) (PendingInteraction, error)
	// UpdatePendingStatus 更新等待交互终态，供 cancel/timeout 阻断旧 resume。
	UpdatePendingStatus(ctx context.Context, pendingID string, status PendingStatus) error
	// ConsumeResumeRef 消费 resume ref，避免同一审批或澄清被重复提交。
	ConsumeResumeRef(ctx context.Context, resumeRef string, submittedStatus PendingStatus) (PendingInteraction, error)
	// ConsumeResumeRefWithIdempotency 在同一事务语义下处理 resume 和幂等记录。
	ConsumeResumeRefWithIdempotency(ctx context.Context, resumeRef string, submittedStatus PendingStatus, record IdempotencyRecord) (PendingInteraction, IdempotencyRecord, bool, error)
	// AppendAuditEvent 追加审计事件，供 approval、replay 和验收追踪复用。
	AppendAuditEvent(ctx context.Context, event AuditEvent) error
	// SaveContextSnapshot 保存安全上下文摘要，证明 LLM 输入只来自投影事实。
	SaveContextSnapshot(ctx context.Context, snapshot ContextSnapshot) error
}
