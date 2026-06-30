package httpapi_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/httpapi"
	"agent-platform-eino/internal/einoapp/product"
)

func TestCurrentViewReturnsWorkbenchProjectionEnvelope(t *testing.T) {
	server := httptest.NewServer(httpapi.NewRouter(emptyProjectionDeps()))
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/workspaces/ws_123/views/current")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	if body["schema_version"] != "eino_workbench_view.v1" {
		t.Fatalf("schema_version = %v", body["schema_version"])
	}
	if body["workspace_id"] != "ws_123" {
		t.Fatalf("workspace_id = %v", body["workspace_id"])
	}
	if body["status"] != "created" {
		t.Fatalf("status = %v", body["status"])
	}
	if _, ok := body["timeline"].([]any); !ok {
		t.Fatalf("timeline must be an array")
	}
	if _, ok := body["inspector"].(map[string]any); !ok {
		t.Fatalf("inspector must be an object")
	}
}

func TestPostAgentActionReturnsAcceptedActionResult(t *testing.T) {
	server := httptest.NewServer(httpapi.NewRouter(emptyProjectionDeps()))
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/workspaces/ws_123/agent/actions", "application/json", strings.NewReader(`{
		"schema_version": "eino_action_request.v1",
		"action_id": "action-demo",
		"client_request_id": "client-1",
		"input": {"text": "hello"}
	}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	if body["schema_version"] != "eino_action_result.v1" {
		t.Fatalf("schema_version = %v", body["schema_version"])
	}
	if body["workspace_id"] != "ws_123" {
		t.Fatalf("workspace_id = %v", body["workspace_id"])
	}
	if body["status"] != "accepted" {
		t.Fatalf("status = %v", body["status"])
	}
	if _, ok := body["result_cards"].([]any); !ok {
		t.Fatalf("result_cards must be an array")
	}
	if _, ok := body["audit_refs"].([]any); !ok {
		t.Fatalf("audit_refs must be an array")
	}
}

func TestStreamEndpointUsesSSEWithStableEventID(t *testing.T) {
	server := httptest.NewServer(httpapi.NewRouter(emptyProjectionDeps()))
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/workspaces/ws_123/runs/run_123/stream")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if contentType := resp.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "text/event-stream") {
		t.Fatalf("Content-Type = %q, want text/event-stream", contentType)
	}
}

func TestLegacyWorkbenchRoutesAreNotMounted(t *testing.T) {
	server := httptest.NewServer(httpapi.NewRouter(emptyProjectionDeps()))
	defer server.Close()

	for _, legacyPath := range []string{"/workbench/", "/chat", "/v2/workbench/chat"} {
		resp, err := http.Get(server.URL + legacyPath)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("%s status = %d, want 404", legacyPath, resp.StatusCode)
		}
	}
}

func TestUnknownRouteReturnsUnifiedErrorEnvelope(t *testing.T) {
	server := httptest.NewServer(httpapi.NewRouter(emptyProjectionDeps()))
	defer server.Close()

	resp, err := http.Get(server.URL + "/unknown")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", resp.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["schema_version"] != "eino_error_envelope.v1" {
		t.Fatalf("schema_version = %v", body["schema_version"])
	}
	if body["request_id"] == "" {
		t.Fatalf("request_id is empty")
	}
}

func TestRouterUsesInjectedProjectionForCurrentView(t *testing.T) {
	projection := recordingProjection{
		view: product.WorkbenchView{
			SchemaVersion: "eino_workbench_view.v1",
			WorkspaceID:   "ws_injected",
			RunID:         "run-injected",
			Status:        "created",
			Timeline:      []product.TimelineItem{},
			Inspector:     product.Inspector{Tabs: []string{"evidence"}},
		},
	}
	server := httptest.NewServer(httpapi.NewRouter(httpapi.Dependencies{Projection: &projection, Commands: execution.NewStaticCommands()}))
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/workspaces/ws_injected/views/current")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["run_id"] != "run-injected" {
		t.Fatalf("run_id = %v, want injected projection result", body["run_id"])
	}
	if projection.workspaceID != "ws_injected" {
		t.Fatalf("projection workspaceID = %q", projection.workspaceID)
	}
}

func TestDefaultDependenciesProvideProjection(t *testing.T) {
	deps := httpapi.DefaultDependencies()
	server := httptest.NewServer(httpapi.NewRouter(deps))
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/workspaces/ws_default/views/current")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
}

func TestDefaultDependenciesShareFactsBetweenCommandsAndProjection(t *testing.T) {
	server := httptest.NewServer(httpapi.NewRouter(httpapi.DefaultDependencies()))
	defer server.Close()

	messageResp, err := http.Post(server.URL+"/api/workspaces/ws_shared/messages", "application/json", strings.NewReader(`{
		"schema_version": "eino_workbench_message_request.v1",
		"message": "hello",
		"client_request_id": "client-shared-1"
	}`))
	if err != nil {
		t.Fatal(err)
	}
	defer messageResp.Body.Close()
	if messageResp.StatusCode != http.StatusOK {
		t.Fatalf("message status = %d, want 200", messageResp.StatusCode)
	}
	var messageBody map[string]any
	if err := json.NewDecoder(messageResp.Body).Decode(&messageBody); err != nil {
		t.Fatal(err)
	}

	viewResp, err := http.Get(server.URL + "/api/workspaces/ws_shared/views/current")
	if err != nil {
		t.Fatal(err)
	}
	defer viewResp.Body.Close()
	var viewBody map[string]any
	if err := json.NewDecoder(viewResp.Body).Decode(&viewBody); err != nil {
		t.Fatal(err)
	}

	if viewBody["run_id"] != messageBody["run_id"] {
		t.Fatalf("view run_id = %v, want message run_id %v", viewBody["run_id"], messageBody["run_id"])
	}
}

func TestEmptyDependenciesUseSharedDefaultFacts(t *testing.T) {
	server := httptest.NewServer(httpapi.NewRouter(httpapi.Dependencies{}))
	defer server.Close()

	messageResp, err := http.Post(server.URL+"/api/workspaces/ws_empty/messages", "application/json", strings.NewReader(`{
		"schema_version": "eino_workbench_message_request.v1",
		"message": "hello",
		"client_request_id": "client-empty-1"
	}`))
	if err != nil {
		t.Fatal(err)
	}
	defer messageResp.Body.Close()
	var messageBody map[string]any
	if err := json.NewDecoder(messageResp.Body).Decode(&messageBody); err != nil {
		t.Fatal(err)
	}

	viewResp, err := http.Get(server.URL + "/api/workspaces/ws_empty/views/current")
	if err != nil {
		t.Fatal(err)
	}
	defer viewResp.Body.Close()
	var viewBody map[string]any
	if err := json.NewDecoder(viewResp.Body).Decode(&viewBody); err != nil {
		t.Fatal(err)
	}

	if viewBody["run_id"] != messageBody["run_id"] {
		t.Fatalf("view run_id = %v, want message run_id %v", viewBody["run_id"], messageBody["run_id"])
	}
}

func TestPartialDependenciesPanicToAvoidSplitFacts(t *testing.T) {
	t.Run("projection only", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic")
			}
		}()
		_ = httpapi.NewRouter(httpapi.Dependencies{Projection: product.NewEmptyProjection()})
	})

	t.Run("commands only", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Fatal("expected panic")
			}
		}()
		_ = httpapi.NewRouter(httpapi.Dependencies{Commands: &recordingCommands{}})
	})
}

func TestActionRequestParsesActionIDBeforeProjection(t *testing.T) {
	projection := recordingProjection{}
	commands := recordingCommands{acceptedRunID: "run-from-command"}
	server := httptest.NewServer(httpapi.NewRouter(httpapi.Dependencies{Projection: &projection, Commands: &commands}))
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/workspaces/ws_123/agent/actions", "application/json", strings.NewReader(`{
		"schema_version": "eino_action_request.v1",
		"action_id": "action-from-request",
		"client_request_id": "client-1",
		"capability_hint": "asset.lookup",
		"input": {
			"text": "hello",
			"attachments": [{"attachment_ref": "attachment-safe-1", "media_type": "text/plain", "safe_name": "query.txt"}]
		},
		"context": {"timezone": "Asia/Shanghai", "locale": "zh-CN", "safe_user_label": "analyst"}
	}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if projection.actionID != "action-from-request" {
		t.Fatalf("projection actionID = %q", projection.actionID)
	}
	if commands.action.ActionID != "action-from-request" {
		t.Fatalf("command actionID = %q", commands.action.ActionID)
	}
	if commands.action.CapabilityHint != "asset.lookup" {
		t.Fatalf("command capability hint = %q", commands.action.CapabilityHint)
	}
	if len(commands.action.Attachments) != 1 || commands.action.Attachments[0].AttachmentRef != "attachment-safe-1" {
		t.Fatalf("command attachments not populated: %+v", commands.action.Attachments)
	}
	if commands.action.Context.Locale != "zh-CN" {
		t.Fatalf("command context not populated: %+v", commands.action.Context)
	}
	if projection.runID != "run-from-command" {
		t.Fatalf("projection runID = %q", projection.runID)
	}
}

func TestMessageRequestParsesIntoExecutionCommand(t *testing.T) {
	commands := recordingCommands{acceptedRunID: "run-message-command"}
	server := httptest.NewServer(httpapi.NewRouter(httpapi.Dependencies{Projection: product.NewEmptyProjection(), Commands: &commands}))
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/workspaces/ws_123/messages", "application/json", strings.NewReader(`{
		"schema_version": "eino_workbench_message_request.v1",
		"message": "hello",
		"client_request_id": "client-message-1"
	}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["run_id"] != "run-message-command" {
		t.Fatalf("run_id = %v", body["run_id"])
	}
	if commands.message.Message != "hello" || commands.message.ClientRequestID != "client-message-1" {
		t.Fatalf("message command not populated: %+v", commands.message)
	}
}

func TestInvalidMessageRequestReturnsUnifiedErrorEnvelope(t *testing.T) {
	server := httptest.NewServer(httpapi.NewRouter(emptyProjectionDeps()))
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/workspaces/ws_123/messages", "application/json", strings.NewReader(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assertErrorEnvelope(t, resp, http.StatusBadRequest, "invalid_request")
}

func TestInvalidActionRequestReturnsUnifiedErrorEnvelope(t *testing.T) {
	server := httptest.NewServer(httpapi.NewRouter(emptyProjectionDeps()))
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/workspaces/ws_123/agent/actions", "application/json", strings.NewReader(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assertErrorEnvelope(t, resp, http.StatusBadRequest, "invalid_request")
}

func TestResumeRequestParsesIntoExecutionCommand(t *testing.T) {
	projection := recordingProjection{}
	commands := recordingCommands{acceptedRunID: "run-resume-command"}
	server := httptest.NewServer(httpapi.NewRouter(httpapi.Dependencies{Projection: &projection, Commands: &commands}))
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/workspaces/ws_123/runs/run_existing/resume", "application/json", strings.NewReader(`{
		"schema_version": "eino_workbench_resume_request.v1",
		"resume_ref": "resume_ref_1",
		"client_request_id": "client-resume-1",
		"decision": "approve",
		"comment": "ok"
	}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if commands.resume.RunID != "run_existing" || commands.resume.ResumeRef != "resume_ref_1" {
		t.Fatalf("resume command not populated: %+v", commands.resume)
	}
	if commands.resume.Decision != "approve" || commands.resume.Comment != "ok" {
		t.Fatalf("resume decision not populated: %+v", commands.resume)
	}
	if projection.runID != "run-resume-command" {
		t.Fatalf("projection runID = %q", projection.runID)
	}
	if projection.resumeCalls != 1 {
		t.Fatalf("resume projection calls = %d", projection.resumeCalls)
	}
}

func TestResumeMissingRunReturnsUnifiedErrorEnvelope(t *testing.T) {
	projection := recordingProjection{}
	commands := recordingCommands{resumeErr: execution.ErrRunNotFound}
	server := httptest.NewServer(httpapi.NewRouter(httpapi.Dependencies{Projection: &projection, Commands: &commands}))
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/workspaces/ws_123/runs/run_missing/resume", "application/json", strings.NewReader(`{
		"schema_version": "eino_workbench_resume_request.v1",
		"resume_ref": "resume_ref_1",
		"client_request_id": "client-resume-1",
		"decision": "approve"
	}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assertErrorEnvelope(t, resp, http.StatusNotFound, "run_not_found")
	if projection.resumeCalls != 0 {
		t.Fatalf("resume projection calls = %d", projection.resumeCalls)
	}
}

func TestInvalidResumeRequestReturnsUnifiedErrorEnvelope(t *testing.T) {
	server := httptest.NewServer(httpapi.NewRouter(emptyProjectionDeps()))
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/workspaces/ws_123/runs/run_existing/resume", "application/json", strings.NewReader(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assertErrorEnvelope(t, resp, http.StatusBadRequest, "invalid_request")
}

func TestTrailingJSONReturnsUnifiedErrorEnvelope(t *testing.T) {
	server := httptest.NewServer(httpapi.NewRouter(emptyProjectionDeps()))
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/workspaces/ws_123/messages", "application/json", strings.NewReader(`{
		"schema_version": "eino_workbench_message_request.v1",
		"message": "hello",
		"client_request_id": "client-message-1"
	}{}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	assertErrorEnvelope(t, resp, http.StatusBadRequest, "invalid_request")
}

func assertErrorEnvelope(t *testing.T, resp *http.Response, wantStatus int, wantCode string) {
	t.Helper()

	if resp.StatusCode != wantStatus {
		t.Fatalf("status = %d, want %d", resp.StatusCode, wantStatus)
	}
	var body struct {
		SchemaVersion string `json:"schema_version"`
		Error         struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.SchemaVersion != "eino_error_envelope.v1" {
		t.Fatalf("schema_version = %q", body.SchemaVersion)
	}
	if body.Error.Code != wantCode {
		t.Fatalf("error.code = %q, want %q", body.Error.Code, wantCode)
	}
}

type recordingProjection struct {
	view        product.WorkbenchView
	workspaceID string
	actionID    string
	runID       string
	resumeCalls int
}

func (projection *recordingProjection) WorkbenchView(_ context.Context, workspaceID string) (product.WorkbenchView, error) {
	projection.workspaceID = workspaceID
	if projection.view.SchemaVersion != "" {
		return projection.view, nil
	}
	return product.NewEmptyProjection().WorkbenchView(context.Background(), workspaceID)
}

func (projection *recordingProjection) RunSnapshot(_ context.Context, workspaceID string, runID string) (product.WorkbenchView, error) {
	return product.NewEmptyProjection().RunSnapshot(context.Background(), workspaceID, runID)
}

func (projection *recordingProjection) ActionResult(_ context.Context, workspaceID string, actionID string, runID string) (product.ActionResult, error) {
	projection.workspaceID = workspaceID
	projection.actionID = actionID
	projection.runID = runID
	return product.NewEmptyProjection().ActionResult(context.Background(), workspaceID, actionID, runID)
}

func (projection *recordingProjection) ResumeResult(_ context.Context, workspaceID string, runID string) (product.ActionResult, error) {
	projection.workspaceID = workspaceID
	projection.runID = runID
	projection.resumeCalls++
	return product.NewEmptyProjection().ResumeResult(context.Background(), workspaceID, runID)
}

func (projection *recordingProjection) ReplayView(_ context.Context, workspaceID string, runID string) (product.ReplayView, error) {
	return product.NewEmptyProjection().ReplayView(context.Background(), workspaceID, runID)
}

type recordingCommands struct {
	acceptedRunID string
	message       execution.MessageCommand
	action        execution.ActionCommand
	resume        execution.ResumeCommand
	resumeErr     error
}

func (commands *recordingCommands) StartMessage(_ context.Context, command execution.MessageCommand) (execution.AcceptedRun, error) {
	commands.message = command
	return execution.AcceptedRun{RunID: commands.acceptedRunID, Status: "accepted"}, nil
}

func (commands *recordingCommands) StartAction(_ context.Context, command execution.ActionCommand) (execution.AcceptedRun, error) {
	commands.action = command
	return execution.AcceptedRun{RunID: commands.acceptedRunID, Status: "accepted"}, nil
}

func (commands *recordingCommands) Resume(_ context.Context, command execution.ResumeCommand) (execution.AcceptedRun, error) {
	commands.resume = command
	if commands.resumeErr != nil {
		return execution.AcceptedRun{}, commands.resumeErr
	}
	return execution.AcceptedRun{RunID: commands.acceptedRunID, Status: "accepted"}, nil
}

func emptyProjectionDeps() httpapi.Dependencies {
	return httpapi.Dependencies{
		Projection: product.NewEmptyProjection(),
		Commands:   execution.NewStaticCommands(),
	}
}
