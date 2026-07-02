package capabilities

import (
	"context"
	"errors"
	"strconv"
)

// ErrMCPSessionNotInitialized 表示 tools/list 在 MCP initialize 前被调用。
var ErrMCPSessionNotInitialized = errors.New("mcp session not initialized")

// MCPSessionState 记录 mock MCP initialize 后协商出的能力。
type MCPSessionState struct {
	ProtocolVersion  string
	Initialized      bool
	ToolsListChanged bool
}

// MCPListToolsRequest 对应 MCP tools/list 的分页请求子集。
type MCPListToolsRequest struct {
	Cursor   string
	PageSize int
}

// MCPListToolsResponse 对应 MCP tools/list 的分页响应子集。
type MCPListToolsResponse struct {
	Tools      []MCPToolDefinition
	NextCursor string
}

// Initialize 模拟 MCP initialize + initialized notification 完成后的 session 状态。
func (provider *MockMCPProvider) Initialize(_ context.Context) (MCPSessionState, error) {
	provider.session = MCPSessionState{
		ProtocolVersion:  MCPProtocolVersion20250618,
		Initialized:      true,
		ToolsListChanged: true,
	}
	return provider.session, nil
}

// ListTools 返回 mock MCP tools/list 分页结果，要求先完成 initialize。
func (provider *MockMCPProvider) ListTools(_ context.Context, request MCPListToolsRequest) (MCPListToolsResponse, error) {
	if !provider.session.Initialized {
		return MCPListToolsResponse{}, ErrMCPSessionNotInitialized
	}
	start, err := parseMCPCursor(request.Cursor)
	if err != nil {
		return MCPListToolsResponse{}, err
	}
	if request.PageSize <= 0 || request.PageSize > len(provider.tools) {
		request.PageSize = len(provider.tools)
	}
	if start >= len(provider.tools) {
		return MCPListToolsResponse{}, nil
	}
	end := start + request.PageSize
	if end > len(provider.tools) {
		end = len(provider.tools)
	}
	response := MCPListToolsResponse{
		Tools: append([]MCPToolDefinition(nil), provider.tools[start:end]...),
	}
	if end < len(provider.tools) {
		response.NextCursor = strconv.Itoa(end)
	}
	return response, nil
}

// parseMCPCursor 把 mock cursor 限制为数组下标，避免在分页实现中引入不透明 token。
func parseMCPCursor(cursor string) (int, error) {
	if cursor == "" {
		return 0, nil
	}
	start, err := strconv.Atoi(cursor)
	if err != nil || start < 0 {
		return 0, errors.New("invalid mcp cursor")
	}
	return start, nil
}
