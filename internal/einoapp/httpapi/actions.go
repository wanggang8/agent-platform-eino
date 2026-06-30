package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type workbenchView struct {
	SchemaVersion string         `json:"schema_version"`
	WorkspaceID   string         `json:"workspace_id"`
	RunID         string         `json:"run_id"`
	Status        string         `json:"status"`
	Timeline      []timelineItem `json:"timeline"`
	Inspector     inspector      `json:"inspector"`
}

type timelineItem struct {
	ItemID  string `json:"item_id"`
	Kind    string `json:"kind"`
	Content string `json:"content,omitempty"`
	Status  string `json:"status,omitempty"`
}

type inspector struct {
	Tabs []string `json:"tabs"`
}

type actionResult struct {
	SchemaVersion string       `json:"schema_version"`
	WorkspaceID   string       `json:"workspace_id"`
	ActionID      string       `json:"action_id"`
	RunID         string       `json:"run_id"`
	Status        string       `json:"status"`
	ResultCards   []resultCard `json:"result_cards"`
	AuditRefs     []string     `json:"audit_refs"`
}

type resultCard struct {
	CardID string `json:"card_id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

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

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"schema_version": "eino_workbench_health.v1",
		"status":         "ok",
	})
}

func handleCurrentView(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	writeJSON(w, http.StatusOK, newWorkbenchView(workspaceID, "run_initial"))
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	runID := "run_" + workspaceID + "_accepted"
	writeJSON(w, http.StatusOK, messageResponse{
		SchemaVersion: "eino_workbench_message_response.v1",
		RunID:         runID,
		Status:        "accepted",
		StreamURL:     fmt.Sprintf("/api/workspaces/%s/runs/%s/stream", workspaceID, runID),
		SnapshotURL:   fmt.Sprintf("/api/workspaces/%s/runs/%s", workspaceID, runID),
	})
}

func handleAgentAction(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	writeJSON(w, http.StatusOK, newActionResult(workspaceID, "action_initial", "run_action_initial", "accepted"))
}

func handleResume(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	runID := r.PathValue("run_id")
	writeJSON(w, http.StatusOK, newActionResult(workspaceID, "resume_initial", runID, "accepted"))
}

func handleRunSnapshot(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	runID := r.PathValue("run_id")
	writeJSON(w, http.StatusOK, newWorkbenchView(workspaceID, runID))
}

func handleReplay(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	runID := r.PathValue("run_id")
	writeJSON(w, http.StatusOK, map[string]any{
		"schema_version": "eino_replay_view.v1",
		"workspace_id":   workspaceID,
		"run_id":         runID,
		"events":         []any{},
		"view":           newWorkbenchView(workspaceID, runID),
	})
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

	data, err := json.Marshal(event)
	if err != nil {
		http.Error(w, "failed to encode stream event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "id: %s\nevent: %s\ndata: %s\n\n", event.EventID, event.Type, data)
}

func newWorkbenchView(workspaceID string, runID string) workbenchView {
	return workbenchView{
		SchemaVersion: "eino_workbench_view.v1",
		WorkspaceID:   workspaceID,
		RunID:         runID,
		Status:        "created",
		Timeline:      []timelineItem{},
		Inspector: inspector{
			Tabs: []string{"evidence", "structured", "runtime", "audit"},
		},
	}
}

func newActionResult(workspaceID string, actionID string, runID string, status string) actionResult {
	return actionResult{
		SchemaVersion: "eino_action_result.v1",
		WorkspaceID:   workspaceID,
		ActionID:      actionID,
		RunID:         runID,
		Status:        status,
		ResultCards:   []resultCard{},
		AuditRefs:     []string{},
	}
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
