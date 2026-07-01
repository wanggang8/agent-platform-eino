package capabilities

import (
	"strings"
	"time"
)

// CredentialStatus 是产品可展示的凭据绑定状态，不包含真实 token 或 credential_ref。
type CredentialStatus string

const (
	CredentialStatusConfigured CredentialStatus = "configured"
	CredentialStatusMissing    CredentialStatus = "missing"
	CredentialStatusUnbound    CredentialStatus = "unbound"
	CredentialStatusBound      CredentialStatus = "bound"
)

// CredentialBinding 是 provider policy 可返回给产品层的安全凭据摘要。
type CredentialBinding struct {
	SchemaVersion string           `json:"schema_version"`
	WorkspaceID   string           `json:"workspace_id"`
	System        string           `json:"system"`
	Status        CredentialStatus `json:"status"`
	DisplayRef    string           `json:"display_ref"`
	OwnerScope    PermissionScope  `json:"owner_scope,omitempty"`
	UpdatedAt     string           `json:"updated_at,omitempty"`
	AuditRef      string           `json:"audit_ref,omitempty"`
}

// ConnectorStatus 表示 connector 当前是否可用；只进入策略判断，不替代业务读取工具。
type ConnectorStatus string

const (
	ConnectorStatusUnknown     ConnectorStatus = ""
	ConnectorStatusAvailable   ConnectorStatus = "available"
	ConnectorStatusUnavailable ConnectorStatus = "unavailable"
)

// SafeCredentialBinding 归一化凭据展示字段，避免 credential_ref、token、secret 进入 JSON 输出。
func SafeCredentialBinding(binding CredentialBinding) CredentialBinding {
	binding.SchemaVersion = "eino.provider_credential_binding.v1"
	if unsafeCredentialText(binding.WorkspaceID) || !validSafeIdentifier(binding.WorkspaceID) {
		binding.WorkspaceID = "workspace"
	}
	binding.Status = safeCredentialStatus(binding.Status)
	binding.System = safeCredentialSystem(binding.System)
	if binding.OwnerScope != PermissionScopeWorkspace &&
		binding.OwnerScope != PermissionScopeCaller &&
		binding.OwnerScope != PermissionScopeSystem {
		binding.OwnerScope = PermissionScopeWorkspace
	}
	if unsafeCredentialText(binding.UpdatedAt) || !validCredentialTimestamp(binding.UpdatedAt) {
		binding.UpdatedAt = ""
	}
	binding.DisplayRef = safeCredentialDisplayRef(binding.System, binding.Status, binding.DisplayRef)
	if unsafeCredentialText(binding.AuditRef) || !validCredentialAuditRef(binding.AuditRef) {
		binding.AuditRef = "audit:credential"
	}
	return binding
}

// safeCredentialStatus 将外部输入归一到凭据状态枚举，避免状态字段夹带敏感文本。
func safeCredentialStatus(status CredentialStatus) CredentialStatus {
	switch status {
	case CredentialStatusConfigured, CredentialStatusMissing, CredentialStatusUnbound, CredentialStatusBound:
		return status
	default:
		return CredentialStatusMissing
	}
}

// safeCredentialSystem 限制 system 到 schema 允许的 provider 标识。
func safeCredentialSystem(system string) string {
	if unsafeCredentialText(system) {
		return "local"
	}
	switch system {
	case "fobrain", "llm", "mcp", "local":
		return system
	default:
		return "local"
	}
}

// safeCredentialDisplayRef 只保留 schema 允许的安全 display_ref。
func safeCredentialDisplayRef(system string, status CredentialStatus, displayRef string) string {
	if !unsafeCredentialText(displayRef) && validCredentialDisplayRef(displayRef) {
		return displayRef
	}
	switch status {
	case CredentialStatusBound:
		return "bound:" + safeRefPart(system) + ":local"
	case CredentialStatusConfigured:
		return "configured"
	case CredentialStatusUnbound:
		return "unbound"
	default:
		return "missing"
	}
}

// validCredentialDisplayRef 匹配 provider_credential_binding schema 的安全引用格式。
func validCredentialDisplayRef(value string) bool {
	if value == "configured" || value == "missing" || value == "unbound" {
		return true
	}
	parts := strings.Split(value, ":")
	if len(parts) != 3 || parts[0] != "bound" {
		return false
	}
	return safeRefPart(parts[1]) == parts[1] && safeRefPart(parts[2]) == parts[2]
}

// validCredentialTimestamp 只允许 schema 声明的 date-time，空值表示不输出该可选字段。
func validCredentialTimestamp(value string) bool {
	if strings.TrimSpace(value) == "" {
		return true
	}
	_, err := time.Parse(time.RFC3339, value)
	return err == nil
}

// validCredentialAuditRef 只允许 audit 命名空间下的安全引用，避免 URL/DSN/raw payload 混入摘要。
func validCredentialAuditRef(value string) bool {
	if strings.TrimSpace(value) == "" {
		return true
	}
	parts := strings.Split(value, ":")
	if len(parts) < 2 || parts[0] != "audit" {
		return false
	}
	for _, part := range parts[1:] {
		if !validSafeIdentifier(part) {
			return false
		}
	}
	return true
}

// validSafeIdentifier 限制可展示标识为简单 ASCII 标识，避免把路径、URL 或连接串当作摘要输出。
func validSafeIdentifier(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	for _, char := range value {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' || char == '-' {
			continue
		}
		return false
	}
	return true
}

// safeRefPart 将系统名或绑定名限制为安全小写标识。
func safeRefPart(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return "local"
	}
	var builder strings.Builder
	for _, char := range value {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '_' || char == '-' {
			builder.WriteRune(char)
		}
	}
	if builder.Len() == 0 {
		return "local"
	}
	return builder.String()
}

// unsafeCredentialText 检查凭据摘要里禁止出现的敏感标记。
func unsafeCredentialText(value string) bool {
	lower := strings.ToLower(value)
	for _, marker := range []string{"authorization:", "bearer ", "api_key", "apikey", "token", "password", "credential", "secret", "raw provider", "raw body", "raw error", "raw payload", "raw_payload", "provider_payload", "raw config", "raw_config", "raw-config", "checkpoint-raw", "interrupt-raw"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}
