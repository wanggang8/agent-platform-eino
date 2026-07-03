package facts

import "time"

// RunStatus 是产品层 run 生命周期状态，必须与 docs/facts-contract.md 保持一致。
type RunStatus string

// ToolResultStatus 表示工具结果是否形成了安全结构化结果。
type ToolResultStatus string

// TurnRole 是 Workbench/replay 可展示的消息角色集合。
type TurnRole string

// ToolCallStatus 表示工具调用在 Product Facts 中的产品状态。
type ToolCallStatus string

// PendingKind 区分需要用户介入的审批和澄清场景。
type PendingKind string

// PendingStatus 表示 approval/clarification 的可恢复状态。
type PendingStatus string

// PendingInputMode 表示 clarification 需要用户提供输入的交互形态。
type PendingInputMode string

// IdempotencyScope 区分 resume 和 mutation 的幂等键空间。
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

	PendingInputModeSingleChoice PendingInputMode = "single_choice"
	PendingInputModeMultiChoice  PendingInputMode = "multi_choice"
	PendingInputModeFreeText     PendingInputMode = "free_text"
	PendingInputModeMixed        PendingInputMode = "mixed"

	IdempotencyScopeResume   IdempotencyScope = "resume"
	IdempotencyScopeMutation IdempotencyScope = "mutation"
)

// Valid 校验 run 状态是否属于产品契约允许集合。
func (status RunStatus) Valid() bool {
	switch status {
	case RunStatusCreated, RunStatusRunning, RunStatusWaiting, RunStatusSucceeded, RunStatusFailed, RunStatusCancelled, RunStatusStopped:
		return true
	default:
		return false
	}
}

// Terminal 判断 run 状态是否已经不可继续执行。
func (status RunStatus) Terminal() bool {
	switch status {
	case RunStatusSucceeded, RunStatusFailed, RunStatusCancelled, RunStatusStopped:
		return true
	default:
		return false
	}
}

// Valid 校验 pending input mode 是否属于 JSON Schema 允许集合。
func (mode PendingInputMode) Valid() bool {
	switch mode {
	case PendingInputModeSingleChoice, PendingInputModeMultiChoice, PendingInputModeFreeText, PendingInputModeMixed:
		return true
	default:
		return false
	}
}

// Run 是 Product Facts 的运行根事实，Workbench、Action API 和 replay 均从这里投影状态。
type Run struct {
	RunID       string
	WorkspaceID string
	Status      RunStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ModelLabel  string
	SafeError   string
}

// Turn 是用户、assistant 或系统提示在同一 run 内的有序消息事实。
type Turn struct {
	TurnID    string
	RunID     string
	Role      TurnRole
	Content   string
	Sequence  int64
	CreatedAt time.Time
}

// ToolCall 记录模型选择的能力调用，不包含 raw provider 参数。
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

// ToolResult 只引用 StructuredResult 安全材料，不保存 provider 原始响应。
type ToolResult struct {
	ResultID         string
	ToolCallID       string
	Status           ToolResultStatus
	StructuredResult StructuredResultRef
}

// StructuredResultRef 是工具结果事实的安全引用和摘要。
type StructuredResultRef struct {
	SchemaVersion string
	ResultRef     string
	SafeSummary   string
}

// PendingCandidateField 是 clarification 候选可展示的安全字段，不包含 provider 原始字段名和值。
type PendingCandidateField struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// PendingCandidate 是 Product Facts 中可恢复的澄清候选，仅保存安全引用和安全展示字段。
type PendingCandidate struct {
	CandidateRef string                  `json:"candidate_ref"`
	Label        string                  `json:"label"`
	Description  string                  `json:"description,omitempty"`
	EntityType   string                  `json:"entity_type"`
	SafeFields   []PendingCandidateField `json:"safe_fields,omitempty"`
}

// PendingInteraction 是 approval/clarification 的产品级等待事实。
// CheckpointRef 只能是内部安全引用，不能是 raw Eino checkpoint id。
type PendingInteraction struct {
	PendingID     string
	RunID         string
	Kind          PendingKind
	Status        PendingStatus
	ResumeRef     string
	CheckpointRef string
	Question      string
	OperationName string
	RiskSummary   string
	TargetSummary string
	InputMode     PendingInputMode
	Candidates    []PendingCandidate
	ExpiresAt     time.Time
}

// AuditEvent 是可回放、可审计的安全摘要事件。
type AuditEvent struct {
	AuditID     string
	RunID       string
	EventType   string
	SafeSummary string
	Actor       string
	CreatedAt   time.Time
}

// ContextSnapshot 记录模型调用前的安全上下文摘要，不保存 raw prompt 或 provider payload。
type ContextSnapshot struct {
	SnapshotID  string
	RunID       string
	SafeSummary string
	CreatedAt   time.Time
}

// Snapshot 是按 run 读取的完整 Product Facts 视图，供投影层同源消费。
type Snapshot struct {
	Run                 Run
	Turns               []Turn
	ToolCalls           []ToolCall
	ToolResults         []ToolResult
	PendingInteractions []PendingInteraction
	AuditEvents         []AuditEvent
	ContextSnapshots    []ContextSnapshot
}

// IdempotencyRecord 记录 resume/mutation 请求幂等性，ResourceRef 用于防止同 key 误绑其他资源。
type IdempotencyRecord struct {
	Scope       IdempotencyScope
	Key         string
	RunID       string
	ResourceRef string
	Status      string
	CreatedAt   time.Time
}

// LifecycleTransition 是 run 生命周期的原子事实迁移请求。
type LifecycleTransition struct {
	RunID               string
	ExpectedRunStatuses []RunStatus
	Status              RunStatus
	SafeError           string
	UpdatedAt           time.Time
	PendingIDs          []string
	PendingStatus       PendingStatus
	ToolCallIDs         []string
	ToolStatus          ToolCallStatus
	AuditEvent          AuditEvent
}

// RetryRunTransition 是 retry 新 run 创建的原子事实迁移请求。
type RetryRunTransition struct {
	OldRunID             string
	ExpectedOldRunStatus RunStatus
	ExpectedOldSafeError string
	NewRun               Run
	UserTurn             Turn
	OldAudit             AuditEvent
	NewAudit             AuditEvent
	Idempotency          IdempotencyRecord
}
