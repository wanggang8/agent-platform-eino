package httpapi

import (
	"net/http"

	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

type Dependencies struct {
	Projection product.Projection
	Commands   execution.Commands
}

type api struct {
	projection product.Projection
	commands   execution.Commands
}

func DefaultDependencies() Dependencies {
	repository := facts.NewMemoryRepository()
	return Dependencies{
		Projection: product.NewFactsProjection(repository),
		Commands:   execution.NewFactCommands(repository),
	}
}

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

func (api api) handleNotFound(w http.ResponseWriter, r *http.Request) {
	WriteError(w, r, http.StatusNotFound, product.NewSafeError("not_found", "请求的资源不存在", false))
}
