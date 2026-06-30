package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"agent-platform-eino/internal/einoapp/httpapi"
)

func TestWriteJSONAddsRequestIDHeaderWithoutWrappingSuccessBody(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ok", nil)

	httpapi.WriteJSON(rec, req, http.StatusOK, map[string]any{
		"schema_version": "eino_workbench_view.v1",
		"workspace_id":   "ws-demo",
	})

	if got := rec.Header().Get("X-Request-Id"); got == "" {
		t.Fatal("X-Request-Id header is empty")
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["schema_version"] != "eino_workbench_view.v1" {
		t.Fatalf("schema_version = %v, want business schema direct response", body["schema_version"])
	}
	if _, ok := body["data"]; ok {
		t.Fatalf("success response must not be wrapped in data envelope: %v", body)
	}
}

func TestWriteErrorUsesUnifiedSafeEnvelope(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/bad", nil)

	httpapi.WriteError(rec, req, http.StatusBadRequest, httpapi.APIError{
		Code:       "invalid_request",
		Message:    "请求参数无效",
		SafeDetail: "workspace_id 缺失",
		Retryable:  false,
	})

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}

	var body struct {
		SchemaVersion string `json:"schema_version"`
		RequestID     string `json:"request_id"`
		Error         struct {
			Code       string `json:"code"`
			Message    string `json:"message"`
			SafeDetail string `json:"safe_detail"`
			Retryable  bool   `json:"retryable"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.SchemaVersion != "eino_error_envelope.v1" {
		t.Fatalf("schema_version = %q", body.SchemaVersion)
	}
	if body.RequestID == "" || body.RequestID != rec.Header().Get("X-Request-Id") {
		t.Fatalf("request id body/header mismatch: body=%q header=%q", body.RequestID, rec.Header().Get("X-Request-Id"))
	}
	if body.Error.Code != "invalid_request" || body.Error.SafeDetail != "workspace_id 缺失" {
		t.Fatalf("unexpected error body: %+v", body.Error)
	}
}

func TestWriteJSONGeneratesFreshRequestIDPerRequest(t *testing.T) {
	first := httptest.NewRecorder()
	second := httptest.NewRecorder()

	httpapi.WriteJSON(first, httptest.NewRequest(http.MethodGet, "/first", nil), http.StatusOK, map[string]string{"status": "ok"})
	httpapi.WriteJSON(second, httptest.NewRequest(http.MethodGet, "/second", nil), http.StatusOK, map[string]string{"status": "ok"})

	if first.Header().Get("X-Request-Id") == second.Header().Get("X-Request-Id") {
		t.Fatalf("request ids must be unique, got %q", first.Header().Get("X-Request-Id"))
	}
}
