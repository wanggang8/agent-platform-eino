package capabilities

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

// ErrMCPToolExecutionFailed 表示 MCP tools/call 返回 isError=true。
var ErrMCPToolExecutionFailed = errors.New("mcp tool execution failed")

// ErrMCPInvalidArguments 表示 MCP tools/call 参数不满足 inputSchema 子集。
var ErrMCPInvalidArguments = errors.New("mcp invalid arguments")

// MCPMockConfig 定义可运行 mock MCP provider 的 server catalog 和固定工具结果。
type MCPMockConfig struct {
	ServerID string
	Tools    []MCPToolDefinition
	Results  map[string]MCPToolResult
}

// MockMCPProvider 是 Phase 4.5 的可运行 MCP mock provider，不连接生产 MCP server。
type MockMCPProvider struct {
	serverID     string
	tools        []MCPToolDefinition
	results      map[string]MCPToolResult
	session      MCPSessionState
	capabilities map[string]Capability
	toolByCapID  map[string]MCPToolDefinition
}

var _ Provider = (*MockMCPProvider)(nil)
var _ Invoker = (*MockMCPProvider)(nil)

// NewMockMCPProvider 创建 mock MCP provider，并把 MCP tool catalog 映射为项目 Capability。
func NewMockMCPProvider(config MCPMockConfig) *MockMCPProvider {
	serverID := strings.TrimSpace(config.ServerID)
	if serverID == "" {
		serverID = "mock"
	}
	tools := append([]MCPToolDefinition(nil), config.Tools...)
	slices.SortFunc(tools, func(a, b MCPToolDefinition) int {
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})

	provider := &MockMCPProvider{
		serverID:     serverID,
		tools:        tools,
		results:      cloneMCPResults(config.Results),
		capabilities: map[string]Capability{},
		toolByCapID:  map[string]MCPToolDefinition{},
	}
	for _, tool := range tools {
		capability := provider.capabilityFromTool(tool)
		provider.capabilities[capability.ID] = capability
		provider.toolByCapID[capability.ID] = tool
	}
	return provider
}

// ID 返回 MCP provider id，保持和其他 capability provider 一致。
func (provider *MockMCPProvider) ID() string {
	return "mcp:" + safeAuditPart(provider.serverID)
}

// ListCapabilities 返回由 MCP tools/list catalog 映射出的项目 capability。
func (provider *MockMCPProvider) ListCapabilities() ([]Capability, error) {
	out := make([]Capability, 0, len(provider.capabilities))
	for _, capability := range provider.capabilities {
		out = append(out, capability)
	}
	slices.SortFunc(out, func(a, b Capability) int {
		if a.ID < b.ID {
			return -1
		}
		if a.ID > b.ID {
			return 1
		}
		return 0
	})
	return out, nil
}

// Invoke 模拟 MCP tools/call，并把 structuredContent 收敛为安全 StructuredResult candidate。
func (provider *MockMCPProvider) Invoke(_ context.Context, request InvocationRequest) (product.StructuredResultCandidate, error) {
	if !provider.session.Initialized {
		return product.StructuredResultCandidate{}, ErrMCPSessionNotInitialized
	}
	tool, ok := provider.toolByCapID[request.CapabilityID]
	if !ok {
		return product.StructuredResultCandidate{}, ErrCapabilityInvocationUnavailable
	}
	if err := validateMCPArguments(tool.InputSchema, request.Arguments); err != nil {
		return product.StructuredResultCandidate{}, err
	}
	result, ok := provider.results[tool.Name]
	if !ok {
		return product.StructuredResultCandidate{}, ErrCapabilityInvocationUnavailable
	}
	if result.IsError {
		return product.StructuredResultCandidate{}, ErrMCPToolExecutionFailed
	}
	candidate, err := structuredResultCandidateFromMCP(tool, result)
	if err != nil {
		return product.StructuredResultCandidate{}, err
	}
	if _, err := product.NewStructuredResultSafetyGate().Approve(candidate); err != nil {
		return product.StructuredResultCandidate{}, err
	}
	return candidate, nil
}

// validateMCPArguments 执行 mock MCP inputSchema 的最小 required/type 校验。
func validateMCPArguments(schema MCPJSONSchema, arguments map[string]any) error {
	for _, required := range schema.Required {
		if _, ok := arguments[required]; !ok {
			return ErrMCPInvalidArguments
		}
	}
	for name, expectedType := range schema.Properties {
		value, ok := arguments[name]
		if !ok {
			continue
		}
		if !mcpArgumentTypeMatches(expectedType, value) {
			return ErrMCPInvalidArguments
		}
	}
	return nil
}

// mcpArgumentTypeMatches 将 JSON Schema 基础类型映射到 Go 解码后的值类型。
func mcpArgumentTypeMatches(expectedType string, value any) bool {
	switch strings.TrimSpace(expectedType) {
	case "string":
		_, ok := value.(string)
		return ok
	case "number":
		_, ok := value.(float64)
		return ok
	case "integer":
		number, ok := value.(float64)
		return ok && number == float64(int64(number))
	case "boolean":
		_, ok := value.(bool)
		return ok
	default:
		return false
	}
}

// SetToolResultForTest 替换 mock tools/call 结果，只供测试构造安全边界使用。
func (provider *MockMCPProvider) SetToolResultForTest(name string, result MCPToolResult) {
	if provider.results == nil {
		provider.results = map[string]MCPToolResult{}
	}
	provider.results[name] = result
}

// capabilityFromTool 将 MCP tool metadata 转成项目 capability，项目策略优先于 MCP annotations。
func (provider *MockMCPProvider) capabilityFromTool(tool MCPToolDefinition) Capability {
	policy := mergeMCPPolicy(tool.Annotations, tool.ProjectPolicy)
	return Capability{
		ID:                      "mcp." + safeAuditPart(provider.serverID) + ".tool." + safeAuditPart(tool.Name),
		ProviderID:              provider.ID(),
		ToolName:                "mcp_" + safeAuditPart(tool.Name),
		DisplayName:             safeMCPDisplayName(tool),
		Description:             safeMCPDescription(tool),
		InputSchema:             jsonSchemaFromMCP(tool.InputSchema),
		ResultSchema:            facts.StructuredResultSchemaVersion,
		RiskLevel:               policy.RiskLevel,
		SideEffect:              policy.SideEffect,
		PolicyRef:               policy.PolicyRef,
		PermissionScope:         policy.PermissionScope,
		CredentialBindingPolicy: CredentialBindingNone,
		ConnectorID:             provider.ID(),
		ApprovalRequired:        policy.ApprovalRequired,
		IdempotencyRequired:     policy.IdempotencyRequired,
	}
}

// mergeMCPPolicy 合并 MCP annotations 和项目可信策略；项目策略永远优先。
func mergeMCPPolicy(annotations MCPToolAnnotations, project MCPProjectPolicy) MCPProjectPolicy {
	policy := MCPProjectPolicy{
		RiskLevel:       RiskLow,
		SideEffect:      SideEffectReadExternal,
		PolicyRef:       "policy:mcp:read:v1",
		PermissionScope: PermissionScopeWorkspace,
	}
	if annotations.DestructiveHint {
		policy.RiskLevel = RiskHigh
		policy.SideEffect = SideEffectWriteExternal
		policy.PolicyRef = "policy:mcp:write:v1"
		policy.ApprovalRequired = true
		policy.IdempotencyRequired = true
	}
	if project.RiskLevel != "" {
		policy.RiskLevel = project.RiskLevel
	}
	if project.SideEffect != "" {
		policy.SideEffect = project.SideEffect
	}
	if project.PolicyRef != "" {
		policy.PolicyRef = project.PolicyRef
	}
	if project.PermissionScope != "" {
		policy.PermissionScope = project.PermissionScope
	}
	if project.ApprovalRequired {
		policy.ApprovalRequired = true
	}
	if project.IdempotencyRequired {
		policy.IdempotencyRequired = true
	}
	return policy
}

// jsonSchemaFromMCP 保留 Phase 4.5 支持的 object properties 子集。
func jsonSchemaFromMCP(schema MCPJSONSchema) JSONSchema {
	return JSONSchema{
		SchemaVersion: "json_schema.v1",
		Properties:    cloneStringMap(schema.Properties),
		Required:      append([]string(nil), schema.Required...),
	}
}

// structuredResultCandidateFromMCP 只读取 structuredContent 中的安全字段，不透传 MCP raw content。
func structuredResultCandidateFromMCP(tool MCPToolDefinition, result MCPToolResult) (product.StructuredResultCandidate, error) {
	resultRef, _ := result.StructuredContent["result_ref"].(string)
	safeSummary, _ := result.StructuredContent["safe_summary"].(string)
	if strings.TrimSpace(resultRef) == "" {
		resultRef = "result:mcp:" + safeAuditPart(tool.Name)
	}
	if strings.TrimSpace(safeSummary) == "" {
		safeSummary = fmt.Sprintf("%s 完成", safeMCPDisplayName(tool))
	}
	return product.StructuredResultCandidate{
		SchemaVersion: facts.StructuredResultSchemaVersion,
		ResultRef:     resultRef,
		SafeSummary:   safeSummary,
	}, nil
}

// safeMCPDescription 防止 MCP catalog 描述把 raw payload、token 或凭据提示带入模型工具描述。
func safeMCPDescription(tool MCPToolDefinition) string {
	description := strings.TrimSpace(tool.Description)
	if description == "" || unsafeCredentialText(description) {
		return "MCP tool"
	}
	return description
}

// safeMCPDisplayName 选择可展示名称，缺失时回退到安全工具名。
func safeMCPDisplayName(tool MCPToolDefinition) string {
	title := strings.TrimSpace(tool.Title)
	if title != "" && !unsafeCredentialText(title) {
		return title
	}
	return safeAuditPart(tool.Name)
}

// cloneMCPResults 复制测试配置，避免调用方后续修改影响 provider 状态。
func cloneMCPResults(input map[string]MCPToolResult) map[string]MCPToolResult {
	out := make(map[string]MCPToolResult, len(input))
	for key, value := range input {
		out[key] = value
	}
	return out
}

// cloneStringMap 复制 schema properties，保持 Capability metadata 不共享外部 map。
func cloneStringMap(input map[string]string) map[string]string {
	out := make(map[string]string, len(input))
	for key, value := range input {
		out[key] = value
	}
	return out
}
