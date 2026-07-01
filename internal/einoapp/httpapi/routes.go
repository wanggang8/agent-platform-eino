package httpapi

import (
	"net/http"

	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

// Dependencies 是 HTTP 层的依赖集合，必须同时提供命令和投影以保持同源事实。
type Dependencies struct {
	Projection product.Projection
	Commands   execution.Commands
}

// api 保存 HTTP handler 需要的边界依赖，不直接依赖 provider、LLM 或 Eino event。
type api struct {
	projection product.Projection
	commands   execution.Commands
}

// DefaultDependencies 创建内存 facts 版本的默认依赖，供本地 smoke 和早期 Phase 使用。
func DefaultDependencies() Dependencies {
	repository := facts.NewMemoryRepository()
	return Dependencies{
		Projection: product.NewFactsProjection(repository),
		Commands:   execution.NewFactCommands(repository),
	}
}

// NewRouter 构建 HTTP 路由；部分依赖缺失时直接 panic，避免半注入导致双事实来源。
func NewRouter(deps Dependencies) http.Handler {
	if deps.Projection == nil && deps.Commands == nil {
		deps = DefaultDependencies()
	}
	if (deps.Projection == nil) != (deps.Commands == nil) {
		panic("httpapi dependencies must provide projection and commands together")
	}
	projection := deps.Projection
	if projection == nil {
		projection = product.NewEmptyProjection()
	}
	commands := deps.Commands
	if commands == nil {
		commands = execution.NewStaticCommands()
	}
	api := api{projection: projection, commands: commands}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", api.handleHealth)
	mux.HandleFunc("GET /api/workspaces/{workspace_id}/views/current", api.handleCurrentView)
	mux.HandleFunc("POST /api/workspaces/{workspace_id}/messages", api.handleMessage)
	mux.HandleFunc("GET /api/workspaces/{workspace_id}/runs/{run_id}/stream", api.handleRunStream)
	mux.HandleFunc("GET /api/workspaces/{workspace_id}/runs/{run_id}", api.handleRunSnapshot)
	mux.HandleFunc("GET /api/workspaces/{workspace_id}/runs/{run_id}/replay", api.handleReplay)
	mux.HandleFunc("POST /api/workspaces/{workspace_id}/runs/{run_id}/resume", api.handleResume)
	mux.HandleFunc("POST /api/workspaces/{workspace_id}/agent/actions", api.handleAgentAction)
	mux.HandleFunc("/", api.handleNotFound)
	return mux
}

// handleNotFound 使用统一错误 envelope 返回 404。
func (api api) handleNotFound(w http.ResponseWriter, r *http.Request) {
	WriteError(w, r, http.StatusNotFound, product.NewSafeError("not_found", "请求的资源不存在", false))
}
