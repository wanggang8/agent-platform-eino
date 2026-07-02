package capabilities

// MCPProtocolVersion20250618 是本阶段 mock adapter 固定复核过的 MCP 协议版本。
const MCPProtocolVersion20250618 = "2025-06-18"

// MCPJSONSchema 是 Phase 4.5 支持的 MCP tool schema 子集，只承载可映射到 Capability 的字段。
type MCPJSONSchema struct {
	Type       string
	Properties map[string]string
	Required   []string
}

// MCPToolAnnotations 是 MCP tool annotations 的最小子集；这些提示不可信，必须与项目策略合并。
type MCPToolAnnotations struct {
	ReadOnlyHint    bool
	DestructiveHint bool
	IdempotentHint  bool
}

// MCPProjectPolicy 是项目对 MCP tool 的可信策略覆盖，优先级高于 MCP annotations。
type MCPProjectPolicy struct {
	RiskLevel           RiskLevel
	SideEffect          SideEffect
	PolicyRef           string
	PermissionScope     PermissionScope
	ApprovalRequired    bool
	IdempotencyRequired bool
}

// MCPToolDefinition 描述 mock MCP server 暴露的工具定义。
type MCPToolDefinition struct {
	Name          string
	Title         string
	Description   string
	InputSchema   MCPJSONSchema
	OutputSchema  MCPJSONSchema
	Annotations   MCPToolAnnotations
	ProjectPolicy MCPProjectPolicy
}

// MCPContent 是 MCP tool result 的内容块子集；Phase 4.5 只把它作为错误摘要来源。
type MCPContent struct {
	Type string
	Text string
}

// MCPToolResult 是 mock MCP tools/call 返回结果，structuredContent 仍需转 StructuredResult candidate。
type MCPToolResult struct {
	Content           []MCPContent
	StructuredContent map[string]any
	IsError           bool
}
