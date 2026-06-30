package httpapi

import (
	"fmt"
	"net/http"
	"time"

	"agent-platform-eino/internal/einoapp/product"
)

type messageResponse struct {
	SchemaVersion string `json:"schema_version"`
	RunID         string `json:"run_id"`
	Status        string `json:"status"`
	StreamURL     string `json:"stream_url"`
	SnapshotURL   string `json:"snapshot_url,omitempty"`
}

type streamEvent struct {
	SchemaVersion string    `json:"schema_version"`
	EventID       string    `json:"event_id"`
	RunID         string    `json:"run_id"`
	Type          string    `json:"type"`
	Sequence      int       `json:"sequence"`
	CreatedAt     time.Time `json:"created_at"`
	Run           runPatch  `json:"run"`
}

type runPatch struct {
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

var projection product.Projection = product.NewEmptyProjection()

func handleHealth(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, r, http.StatusOK, map[string]string{
		"schema_version": "eino_workbench_health.v1",
		"status":         "ok",
	})
}

func handleCurrentView(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	view, err := projection.WorkbenchView(r.Context(), workspaceID)
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("projection_failed", "工作台视图暂不可用", true))
		return
	}
	WriteJSON(w, r, http.StatusOK, view)
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	runID := "run_" + workspaceID + "_accepted"
	WriteJSON(w, r, http.StatusOK, messageResponse{
		SchemaVersion: "eino_workbench_message_response.v1",
		RunID:         runID,
		Status:        "accepted",
		StreamURL:     fmt.Sprintf("/api/workspaces/%s/runs/%s/stream", workspaceID, runID),
		SnapshotURL:   fmt.Sprintf("/api/workspaces/%s/runs/%s", workspaceID, runID),
	})
}

func handleAgentAction(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	result, err := projection.ActionResult(r.Context(), workspaceID, "action_initial", "run_action_initial")
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("projection_failed", "操作结果暂不可用", true))
		return
	}
	WriteJSON(w, r, http.StatusOK, result)
}

func handleResume(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	runID := r.PathValue("run_id")
	result, err := projection.ActionResult(r.Context(), workspaceID, "resume_initial", runID)
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("projection_failed", "恢复结果暂不可用", true))
		return
	}
	WriteJSON(w, r, http.StatusOK, result)
}

func handleRunSnapshot(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	runID := r.PathValue("run_id")
	view, err := projection.RunSnapshot(r.Context(), workspaceID, runID)
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("projection_failed", "运行快照暂不可用", true))
		return
	}
	WriteJSON(w, r, http.StatusOK, view)
}

func handleReplay(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	runID := r.PathValue("run_id")
	replay, err := projection.ReplayView(r.Context(), workspaceID, runID)
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("projection_failed", "回放暂不可用", true))
		return
	}
	WriteJSON(w, r, http.StatusOK, replay)
}

func handleRunStream(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()
	runID := r.PathValue("run_id")
	event := streamEvent{
		SchemaVersion: "eino_workbench_stream_event.v1",
		EventID:       runID + ":000001",
		RunID:         runID,
		Type:          "run.updated",
		Sequence:      1,
		CreatedAt:     now,
		Run: runPatch{
			Status:    "created",
			UpdatedAt: now,
		},
	}

	w.Header().Set(requestIDHeader, requestID(r))
	if err := WriteSSEEvent(w, SSEEvent{
		ID:    event.EventID,
		Event: event.Type,
		Data:  event,
	}); err != nil {
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("sse_encode_failed", "事件编码失败", true))
		return
	}
}
