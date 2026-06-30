package httpapi_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"agent-platform-eino/internal/einoapp/httpapi"
)

func TestEncodeSSEEventWritesIDEventAndJSONData(t *testing.T) {
	var out strings.Builder

	err := httpapi.EncodeSSEEvent(&out, httpapi.SSEEvent{
		ID:    "run-1:000001",
		Event: "run.updated",
		Data: map[string]any{
			"schema_version": "eino_workbench_stream_event.v1",
			"run_id":         "run-1",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	got := out.String()
	for _, want := range []string{
		"id: run-1:000001\n",
		"event: run.updated\n",
		`data: {"run_id":"run-1","schema_version":"eino_workbench_stream_event.v1"}`,
		"\n\n",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("encoded SSE missing %q in %q", want, got)
		}
	}
}

func TestSSEHeadersUseNoCacheStreamDefaults(t *testing.T) {
	header := httpapi.SSEHeaders()

	if got := header.Get("Content-Type"); got != "text/event-stream; charset=utf-8" {
		t.Fatalf("Content-Type = %q", got)
	}
	if got := header.Get("Cache-Control"); got != "no-cache" {
		t.Fatalf("Cache-Control = %q", got)
	}
	if got := header.Get("X-Accel-Buffering"); got != "no" {
		t.Fatalf("X-Accel-Buffering = %q", got)
	}
}

func TestWriteSSEEventSetsHeadersAndFlushes(t *testing.T) {
	writer := newSSETestWriter()

	err := httpapi.WriteSSEEvent(writer, httpapi.SSEEvent{
		ID:    "run-1:000001",
		Event: "run.updated",
		Data:  map[string]string{"schema_version": "eino_workbench_stream_event.v1"},
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := writer.header.Get("Content-Type"); got != "text/event-stream; charset=utf-8" {
		t.Fatalf("Content-Type = %q", got)
	}
	if !writer.flushed {
		t.Fatal("WriteSSEEvent did not flush")
	}
	if !strings.Contains(writer.body.String(), "event: run.updated\n") {
		t.Fatalf("event body missing: %q", writer.body.String())
	}
}

func TestWriteSSEEventReturnsEncodeErrorBeforeWritingStatus(t *testing.T) {
	writer := newSSETestWriter()

	err := httpapi.WriteSSEEvent(writer, httpapi.SSEEvent{
		ID:    "bad",
		Event: "bad",
		Data:  func() {},
	})
	if err == nil {
		t.Fatal("WriteSSEEvent encode error = nil")
	}
	if writer.status != 0 {
		t.Fatalf("status written despite encode error: %d", writer.status)
	}
}

type sseTestWriter struct {
	header  http.Header
	body    strings.Builder
	status  int
	flushed bool
}

func newSSETestWriter() *sseTestWriter {
	return &sseTestWriter{header: make(http.Header)}
}

func (writer *sseTestWriter) Header() http.Header {
	return writer.header
}

func (writer *sseTestWriter) Write(data []byte) (int, error) {
	if writer.status == 0 {
		return 0, errors.New("status must be written before body")
	}
	return writer.body.Write(data)
}

func (writer *sseTestWriter) WriteHeader(statusCode int) {
	writer.status = statusCode
}

func (writer *sseTestWriter) Flush() {
	writer.flushed = true
}
