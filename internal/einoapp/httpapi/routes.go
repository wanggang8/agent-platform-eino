package httpapi

import (
	"net/http"

	"agent-platform-eino/internal/einoapp/product"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealth)
	mux.HandleFunc("GET /api/workspaces/{workspace_id}/views/current", handleCurrentView)
	mux.HandleFunc("POST /api/workspaces/{workspace_id}/messages", handleMessage)
	mux.HandleFunc("GET /api/workspaces/{workspace_id}/runs/{run_id}/stream", handleRunStream)
	mux.HandleFunc("GET /api/workspaces/{workspace_id}/runs/{run_id}", handleRunSnapshot)
	mux.HandleFunc("GET /api/workspaces/{workspace_id}/runs/{run_id}/replay", handleReplay)
	mux.HandleFunc("POST /api/workspaces/{workspace_id}/runs/{run_id}/resume", handleResume)
	mux.HandleFunc("POST /api/workspaces/{workspace_id}/agent/actions", handleAgentAction)
	mux.HandleFunc("/", handleNotFound)
	return mux
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	WriteError(w, r, http.StatusNotFound, product.NewSafeError("not_found", "请求的资源不存在", false))
}
