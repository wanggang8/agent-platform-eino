package capabilities

import "strings"

// PolicyReason 是 provider policy decision 的稳定原因码。
type PolicyReason string

const (
	PolicyReasonAllowed                       = "allowed"
	PolicyReasonPermissionDenied              = "permission_denied"
	PolicyReasonApprovalRequired              = "approval_required"
	PolicyReasonCredentialMissing             = "credential_missing"
	PolicyReasonCredentialScopeDenied         = "credential_scope_denied"
	PolicyReasonConnectorAuthFailure          = "connector_auth_failure"
	PolicyReasonConnectorTimeout              = "connector_timeout"
	PolicyReasonConnectorTransportUnavailable = "connector_transport_unavailable"
	PolicyReasonConnectorExecutionFailed      = "connector_execution_failed"
	PolicyReasonResourceConflict              = "resource_conflict"
	PolicyReasonSchemaMismatch                = "schema_mismatch"
)

// PolicyContext 是 policy checker 的安全输入，不包含真实凭据或 provider raw config。
type PolicyContext struct {
	WorkspaceID       string
	CallerID          string
	CredentialBinding CredentialBinding
	ConnectorStatus   ConnectorStatus
}

// PolicyDecision 是能力选择后的安全策略结果，后续会写入 Product Facts/audit。
type PolicyDecision struct {
	SchemaVersion     string             `json:"schema_version"`
	Allowed           bool               `json:"allowed"`
	ReasonCode        PolicyReason       `json:"reason_code"`
	PolicyRef         string             `json:"policy_ref"`
	AuditRef          string             `json:"audit_ref"`
	ApprovalRequired  bool               `json:"approval_required"`
	SafeSummary       string             `json:"safe_summary,omitempty"`
	CredentialBinding *CredentialBinding `json:"credential_binding,omitempty"`
	RequiresApproval  bool               `json:"-"`
}

// EvaluatePolicy 执行 provider policy 门禁：凭据、workspace scope、connector 状态和写域审批。
func EvaluatePolicy(capability Capability, contexts ...PolicyContext) PolicyDecision {
	ctx := PolicyContext{ConnectorStatus: ConnectorStatusAvailable}
	if len(contexts) > 0 {
		ctx = contexts[0]
		if ctx.ConnectorStatus == ConnectorStatusUnknown {
			ctx.ConnectorStatus = ConnectorStatusAvailable
		}
	}
	binding := SafeCredentialBinding(ctx.CredentialBinding)
	decision := newPolicyDecision(capability, binding)

	if requiresApproval(capability) {
		return decision.block(PolicyReasonApprovalRequired, "能力需要审批后执行", true)
	}
	if requiresCredential(capability) {
		if ctx.WorkspaceID != "" && binding.WorkspaceID != "" && binding.WorkspaceID != ctx.WorkspaceID {
			return decision.block(PolicyReasonCredentialScopeDenied, "凭据绑定不属于当前工作区", false)
		}
		if binding.Status == CredentialStatusMissing || binding.Status == CredentialStatusUnbound {
			return decision.block(PolicyReasonCredentialMissing, "能力缺少可用凭据绑定", false)
		}
	}
	if capability.ConnectorID != "" && ctx.ConnectorStatus == ConnectorStatusUnavailable {
		return decision.block(PolicyReasonConnectorTransportUnavailable, "连接器当前不可用", false)
	}
	decision.SafeSummary = "policy allowed"
	return decision
}

// newPolicyDecision 构造基础决策对象，确保 policy_ref、audit_ref 和凭据摘要稳定。
func newPolicyDecision(capability Capability, binding CredentialBinding) PolicyDecision {
	return PolicyDecision{
		SchemaVersion:     "eino.provider_policy_decision.v1",
		Allowed:           true,
		ReasonCode:        PolicyReasonAllowed,
		PolicyRef:         policyRef(capability),
		AuditRef:          "audit:policy:" + safeAuditPart(capability.ID),
		ApprovalRequired:  false,
		SafeSummary:       "policy allowed",
		CredentialBinding: credentialBindingPointer(capability, binding),
	}
}

// block 返回安全拒绝决策，不暴露 provider 或凭据内部细节。
func (decision PolicyDecision) block(reason PolicyReason, summary string, approval bool) PolicyDecision {
	decision.Allowed = false
	decision.ReasonCode = reason
	decision.ApprovalRequired = approval
	decision.RequiresApproval = approval
	decision.SafeSummary = summary
	return decision
}

// requiresCredential 判断能力是否需要凭据绑定。
func requiresCredential(capability Capability) bool {
	return capability.CredentialBindingPolicy == CredentialBindingRequired
}

// requiresApproval 判断能力是否必须进入 HITL approval。
func requiresApproval(capability Capability) bool {
	return capability.ApprovalRequired || capability.RiskLevel == RiskWrite || capability.SideEffect == SideEffectWriteExternal
}

// policyRef 返回显式策略引用；旧配置未声明时按 provider/risk 派生安全引用。
func policyRef(capability Capability) string {
	if validPolicyRef(capability.PolicyRef) {
		return capability.PolicyRef
	}
	risk := string(capability.RiskLevel)
	if risk == "" {
		risk = "default"
	}
	return "policy:" + safeAuditPart(capability.ProviderID) + ":" + safeAuditPart(risk) + ":v1"
}

// validPolicyRef 只允许 policy:<safe-part...> 形态进入产品出口，避免 URL/DSN 被当作策略引用展示。
func validPolicyRef(value string) bool {
	value = strings.TrimSpace(value)
	parts := strings.Split(value, ":")
	if len(parts) < 2 || parts[0] != "policy" || unsafeCredentialText(value) {
		return false
	}
	for _, part := range parts[1:] {
		if safeAuditPart(part) != part {
			return false
		}
	}
	return true
}

// credentialBindingPointer 只有需要凭据语义的能力才把安全绑定写入 decision。
func credentialBindingPointer(capability Capability, binding CredentialBinding) *CredentialBinding {
	if (capability.CredentialBindingPolicy == "" || capability.CredentialBindingPolicy == CredentialBindingNone) &&
		binding.Status == CredentialStatusMissing {
		return nil
	}
	return &binding
}

// safeAuditPart 将 audit 引用片段限制为不含敏感材料的稳定标识。
func safeAuditPart(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" || unsafeCredentialText(value) {
		return "unknown"
	}
	var builder strings.Builder
	for _, char := range value {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '_' || char == '-' || char == '.' {
			builder.WriteRune(char)
		}
	}
	if builder.Len() == 0 {
		return "unknown"
	}
	return builder.String()
}
