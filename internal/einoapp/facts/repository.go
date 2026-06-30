package facts

import "context"

type Repository interface {
	CreateRun(ctx context.Context, run Run) error
	GetRun(ctx context.Context, runID string) (Run, error)
	AppendTurn(ctx context.Context, turn Turn) error
	AppendToolCall(ctx context.Context, call ToolCall) error
	AppendToolResult(ctx context.Context, result ToolResult) error
	AppendPendingInteraction(ctx context.Context, pending PendingInteraction) error
	AppendAuditEvent(ctx context.Context, event AuditEvent) error
	SaveContextSnapshot(ctx context.Context, snapshot ContextSnapshot) error
}
