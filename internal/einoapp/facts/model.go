package facts

import "time"

type RunStatus string
type ToolResultStatus string

const (
	RunStatusCreated RunStatus = "created"
	RunStatusRunning RunStatus = "running"
	RunStatusWaiting RunStatus = "waiting"
	RunStatusFailed  RunStatus = "failed"
	RunStatusDone    RunStatus = "done"

	ToolResultSucceeded ToolResultStatus = "succeeded"
	ToolResultFailed    ToolResultStatus = "failed"
)

type Run struct {
	RunID       string
	WorkspaceID string
	Status      RunStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Turn struct {
	TurnID    string
	RunID     string
	Role      string
	Content   string
	CreatedAt time.Time
}

type ToolCall struct {
	ToolCallID string
	RunID      string
	ToolName   string
	Status     string
	CreatedAt  time.Time
}

type ToolResult struct {
	ResultID         string
	ToolCallID       string
	Status           ToolResultStatus
	StructuredResult StructuredResultRef
}

type StructuredResultRef struct {
	SchemaVersion string
	ResultRef     string
	SafeSummary   string
}

type PendingInteraction struct {
	PendingID string
	RunID     string
	Kind      string
	Status    string
}

type AuditEvent struct {
	AuditID     string
	RunID       string
	EventType   string
	SafeSummary string
}

type ContextSnapshot struct {
	SnapshotID  string
	RunID       string
	SafeSummary string
}
