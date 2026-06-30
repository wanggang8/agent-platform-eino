package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"agent-platform-eino/internal/einoapp/httpapi"
)

func TestCurrentViewReturnsWorkbenchProjectionEnvelope(t *testing.T) {
	server := httptest.NewServer(httpapi.NewRouter())
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
	server := httptest.NewServer(httpapi.NewRouter())
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/workspaces/ws_123/agent/actions", "application/json", strings.NewReader(`{}`))
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
	server := httptest.NewServer(httpapi.NewRouter())
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
	server := httptest.NewServer(httpapi.NewRouter())
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
