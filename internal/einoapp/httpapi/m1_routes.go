package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"

	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

type M1Dependencies struct {
	Projection product.M1ProjectionPort
	Executor   execution.M1MessageExecutor
	NewRunID   func() string
	Readiness  func(context.Context) (bool, string)
}

type m1MessageRequest struct {
	ConversationID string `json:"conversation_id"`
	ActorID        string `json:"actor_id"`
	Content        string `json:"content"`
}

var m1IdentifierPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)

// NewM1Router 只注册 Story 1.2 的 read-only 边界；Action/Resume/Replay/History 不存在。
func NewM1Router(dependencies M1Dependencies) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, r, http.StatusOK, map[string]any{"status": "ok"})
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		if dependencies.Readiness == nil {
			WriteError(w, r, http.StatusServiceUnavailable, product.NewSafeError("not_ready", "服务尚未就绪", true))
			return
		}
		ready, reason := dependencies.Readiness(r.Context())
		if !ready {
			WriteError(w, r, http.StatusServiceUnavailable, product.NewSafeError("not_ready", safeReadinessReason(reason), true))
			return
		}
		WriteJSON(w, r, http.StatusOK, map[string]any{"status": "ready"})
	})
	mux.HandleFunc("GET /api/workspaces/{workspace_id}/views/current", func(w http.ResponseWriter, r *http.Request) {
		workspaceID := r.PathValue("workspace_id")
		if !validM1Identifier(workspaceID) || dependencies.Projection == nil {
			writeM1BoundaryError(w, r, facts.ErrNotFound)
			return
		}
		view, err := dependencies.Projection.WorkbenchView(r.Context(), workspaceID)
		if err != nil {
			writeM1BoundaryError(w, r, err)
			return
		}
		WriteJSON(w, r, http.StatusOK, view)
	})
	mux.HandleFunc("POST /api/workspaces/{workspace_id}/messages", func(w http.ResponseWriter, r *http.Request) {
		workspaceID := r.PathValue("workspace_id")
		request, err := decodeM1MessageRequest(r)
		if err != nil || !validM1Identifier(workspaceID) || !validM1Identifier(request.ConversationID) || !validM1Identifier(request.ActorID) || strings.TrimSpace(request.Content) == "" || len(request.Content) > 2000 || facts.ContainsUnsafeMaterial(request.Content) {
			WriteError(w, r, http.StatusBadRequest, product.NewSafeError("invalid_request", "请求格式无效", false))
			return
		}
		if dependencies.Executor == nil || dependencies.Projection == nil {
			WriteError(w, r, http.StatusServiceUnavailable, product.NewSafeError("not_ready", "查询服务尚未就绪", true))
			return
		}
		runID := newM1RunID(dependencies.NewRunID)
		if !validM1Identifier(runID) {
			WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("internal_error", "无法创建查询", false))
			return
		}
		command := execution.M1MessageCommand{
			WorkspaceID: workspaceID, ConversationID: request.ConversationID, ActorID: request.ActorID,
			RunID: runID, ToolCallID: "call_" + runID, Content: request.Content,
		}
		if err := dependencies.Executor.Execute(r.Context(), command); err != nil {
			writeM1BoundaryError(w, r, err)
			return
		}
		view, err := dependencies.Projection.RunSnapshot(r.Context(), workspaceID, runID)
		if err != nil {
			writeM1BoundaryError(w, r, err)
			return
		}
		WriteJSON(w, r, http.StatusOK, view)
	})
	mux.HandleFunc("GET /api/workspaces/{workspace_id}/runs/{run_id}", func(w http.ResponseWriter, r *http.Request) {
		workspaceID, runID := r.PathValue("workspace_id"), r.PathValue("run_id")
		if !validM1Identifier(workspaceID) || !validM1Identifier(runID) || dependencies.Projection == nil {
			writeM1BoundaryError(w, r, facts.ErrNotFound)
			return
		}
		view, err := dependencies.Projection.RunSnapshot(r.Context(), workspaceID, runID)
		if err != nil {
			writeM1BoundaryError(w, r, err)
			return
		}
		WriteJSON(w, r, http.StatusOK, view)
	})
	mux.HandleFunc("GET /api/workspaces/{workspace_id}/runs/{run_id}/stream", func(w http.ResponseWriter, r *http.Request) {
		workspaceID, runID := r.PathValue("workspace_id"), r.PathValue("run_id")
		if !validM1Identifier(workspaceID) || !validM1Identifier(runID) || dependencies.Projection == nil {
			writeM1BoundaryError(w, r, facts.ErrNotFound)
			return
		}
		events, err := dependencies.Projection.StreamEvents(r.Context(), workspaceID, runID)
		if err != nil {
			writeM1BoundaryError(w, r, err)
			return
		}
		for key, values := range SSEHeaders() {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(http.StatusOK)
		for _, event := range events {
			if err := EncodeSSEEvent(w, SSEEvent{ID: event.EventID, Event: event.Type, Data: event}); err != nil {
				return
			}
		}
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, r, http.StatusNotFound, product.NewSafeError("not_found", "资源不存在", false))
	})
	return mux
}

func decodeM1MessageRequest(r *http.Request) (m1MessageRequest, error) {
	defer r.Body.Close()
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	var request m1MessageRequest
	if err := decoder.Decode(&request); err != nil {
		return m1MessageRequest{}, err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return m1MessageRequest{}, errors.New("multiple JSON values")
	}
	return request, nil
}

func writeM1BoundaryError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, facts.ErrNotFound):
		WriteError(w, r, http.StatusNotFound, product.NewSafeError("not_found", "查询结果不存在", false))
	case errors.Is(err, facts.ErrUnsafeFactMaterial):
		WriteError(w, r, http.StatusBadRequest, product.NewSafeError("unsafe_material", "请求未通过安全校验", false))
	default:
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("query_failed", "查询执行失败", false))
	}
}

func validM1Identifier(value string) bool { return m1IdentifierPattern.MatchString(value) }

func newM1RunID(generator func() string) string {
	if generator != nil {
		return generator()
	}
	var raw [12]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return ""
	}
	return "run_" + hex.EncodeToString(raw[:])
}

func safeReadinessReason(reason string) string {
	if strings.TrimSpace(reason) == "" || facts.ContainsUnsafeMaterial(reason) {
		return "服务尚未就绪"
	}
	return reason
}
