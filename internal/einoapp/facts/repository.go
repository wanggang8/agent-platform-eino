package facts

import (
	"context"
	"time"
)

type Repository interface {
	CreateRun(ctx context.Context, run Run) error
	GetRun(ctx context.Context, runID string) (Run, error)
	LatestRun(ctx context.Context, workspaceID string) (Run, error)
	GetSnapshot(ctx context.Context, runID string) (Snapshot, error)
	UpdateRunStatus(ctx context.Context, runID string, status RunStatus, safeError string, updatedAt time.Time) error
	RecordIdempotency(ctx context.Context, record IdempotencyRecord) (IdempotencyRecord, bool, error)
	AppendTurn(ctx context.Context, turn Turn) error
	AppendToolCall(ctx context.Context, call ToolCall) error
	AppendToolResult(ctx context.Context, result ToolResult) error
	AppendPendingInteraction(ctx context.Context, pending PendingInteraction) error
	ConsumeResumeRef(ctx context.Context, resumeRef string, submittedStatus PendingStatus) (PendingInteraction, error)
	ConsumeResumeRefWithIdempotency(ctx context.Context, resumeRef string, submittedStatus PendingStatus, record IdempotencyRecord) (PendingInteraction, IdempotencyRecord, bool, error)
	AppendAuditEvent(ctx context.Context, event AuditEvent) error
	SaveContextSnapshot(ctx context.Context, snapshot ContextSnapshot) error
}
