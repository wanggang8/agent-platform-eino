package facts

import "time"

type RunStatus string
type ToolResultStatus string
type TurnRole string
type ToolCallStatus string
type PendingKind string
type PendingStatus string
type IdempotencyScope string

const (
	RunStatusCreated   RunStatus = "created"
	RunStatusRunning   RunStatus = "running"
	RunStatusWaiting   RunStatus = "waiting"
	RunStatusSucceeded RunStatus = "succeeded"
	RunStatusFailed    RunStatus = "failed"
	RunStatusCancelled RunStatus = "cancelled"
	RunStatusStopped   RunStatus = "stopped"

	TurnRoleUser         TurnRole = "user"
	TurnRoleAssistant    TurnRole = "assistant"
	TurnRoleSystemNotice TurnRole = "system_notice"

	ToolCallQueued    ToolCallStatus = "queued"
	ToolCallRunning   ToolCallStatus = "running"
	ToolCallSucceeded ToolCallStatus = "succeeded"
	ToolCallFailed    ToolCallStatus = "failed"
	ToolCallCancelled ToolCallStatus = "cancelled"

	ToolResultSucceeded ToolResultStatus = "succeeded"
	ToolResultFailed    ToolResultStatus = "failed"

	PendingKindApproval      PendingKind = "approval"
	PendingKindClarification PendingKind = "clarification"

	PendingStatusWaiting   PendingStatus = "waiting"
	PendingStatusSubmitted PendingStatus = "submitted"
	PendingStatusApproved  PendingStatus = "approved"
	PendingStatusRejected  PendingStatus = "rejected"
	PendingStatusCancelled PendingStatus = "cancelled"
	PendingStatusExpired   PendingStatus = "expired"
	PendingStatusConsumed  PendingStatus = "consumed"

	IdempotencyScopeResume   IdempotencyScope = "resume"
	IdempotencyScopeMutation IdempotencyScope = "mutation"
)

func (status RunStatus) Valid() bool {
	switch status {
	case RunStatusCreated, RunStatusRunning, RunStatusWaiting, RunStatusSucceeded, RunStatusFailed, RunStatusCancelled, RunStatusStopped:
		return true
	default:
		return false
	}
}

func (status RunStatus) Terminal() bool {
	switch status {
	case RunStatusSucceeded, RunStatusFailed, RunStatusCancelled, RunStatusStopped:
		return true
	default:
		return false
	}
}

type Run struct {
	RunID       string
	WorkspaceID string
	Status      RunStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ModelLabel  string
	SafeError   string
}

type Turn struct {
	TurnID    string
	RunID     string
	Role      TurnRole
	Content   string
	Sequence  int64
	CreatedAt time.Time
}

type ToolCall struct {
	ToolCallID  string
	RunID       string
	ToolID      string
	DisplayName string
	Status      ToolCallStatus
	ArgsHash    string
	ArgsPreview string
	CreatedAt   time.Time
	EndedAt     time.Time
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
	PendingID     string
	RunID         string
	Kind          PendingKind
	Status        PendingStatus
	ResumeRef     string
	CheckpointRef string
	Question      string
	RiskSummary   string
	ExpiresAt     time.Time
}

type AuditEvent struct {
	AuditID     string
	RunID       string
	EventType   string
	SafeSummary string
	Actor       string
	CreatedAt   time.Time
}

type ContextSnapshot struct {
	SnapshotID  string
	RunID       string
	SafeSummary string
	CreatedAt   time.Time
}

type Snapshot struct {
	Run                 Run
	Turns               []Turn
	ToolCalls           []ToolCall
	ToolResults         []ToolResult
	PendingInteractions []PendingInteraction
	AuditEvents         []AuditEvent
	ContextSnapshots    []ContextSnapshot
}

type IdempotencyRecord struct {
	Scope       IdempotencyScope
	Key         string
	RunID       string
	ResourceRef string
	Status      string
	CreatedAt   time.Time
}
