package httpapi_test

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/httpapi"
	"agent-platform-eino/internal/einoapp/product"
)

type apiClock struct{ now time.Time }

func (clock apiClock) Now() time.Time { return clock.now }

type apiFixtureExecutor struct {
	repository facts.QueryRepository
}

// Execute 在 HTTP 边界测试中写入固定 Product Facts；真实 Eino/fixture 链由 execution 集成测试覆盖。
func (executor apiFixtureExecutor) Execute(ctx context.Context, command execution.M1MessageCommand) error {
	now := time.Date(2026, 7, 27, 11, 0, 0, 0, time.UTC)
	items := []facts.SnapshotItem{
		facts.MustSnapshotItem("item_3", "CVE-2026-0003 · 10.0.0.3", now.Add(-time.Minute)),
		facts.MustSnapshotItem("item_2", "CVE-2026-0002 · 10.0.0.2", now.Add(-2*time.Minute)),
		facts.MustSnapshotItem("item_1", "CVE-2026-0001 · 10.0.0.1", now.Add(-3*time.Minute)),
	}
	result, err := facts.NewStructuredResult(facts.StructuredResultInput{
		InternalResultRef: "result_ref:" + command.RunID, Status: facts.QueryResultResolved,
		Summary: "发现 3 条新增漏洞。", ObservedAt: now, Items: items,
	})
	if err != nil {
		return err
	}
	freshness, err := facts.NewFreshnessPolicy("m1.v1", 30*time.Minute)
	if err != nil {
		return err
	}
	snapshot, err := facts.NewQueryResultSnapshot(facts.QueryResultSnapshotInput{
		SnapshotID: "snapshot_" + command.RunID, WorkspaceID: command.WorkspaceID,
		ConversationID: command.ConversationID, ActorID: command.ActorID, RunID: command.RunID,
		ToolCallID: command.ToolCallID, Result: result, Clock: apiClock{now: now}, Freshness: freshness,
	})
	if err != nil {
		return err
	}
	aggregate, err := facts.NewSucceededQueryFacts(facts.SucceededQueryFactsInput{
		WorkspaceID: command.WorkspaceID, ConversationID: command.ConversationID, ActorID: command.ActorID,
		RunID: command.RunID, QuerySequence: 1, CreatedAt: now, Result: result, Snapshot: snapshot,
	})
	if err != nil {
		return err
	}
	return executor.repository.CommitQuery(ctx, aggregate)
}

func TestM1RouterRunsMessageAndProjectsSameSnapshotAndSSE(t *testing.T) {
	repository := facts.NewQueryMemoryRepository()
	projection := product.NewM1Projection(repository)
	router := httpapi.NewM1Router(httpapi.M1Dependencies{
		Projection: projection, Executor: apiFixtureExecutor{repository: repository}, NewRunID: func() string { return "run_api_1" },
		Readiness: func(context.Context) (bool, string) { return true, "ready" },
	})

	response := performM1Request(router, http.MethodPost, "/api/workspaces/ws-1/messages", `{"conversation_id":"conversation-1","actor_id":"actor-1","content":"查询新增的漏洞"}`)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("message status = %d", response.StatusCode)
	}
	var messageView product.M1WorkbenchView
	if err := json.NewDecoder(response.Body).Decode(&messageView); err != nil {
		t.Fatal(err)
	}
	if messageView.RunID != "run_api_1" || messageView.Result == nil || messageView.Result.Count != 3 {
		t.Fatalf("message view = %+v", messageView)
	}

	snapshotResponse := performM1Request(router, http.MethodGet, "/api/workspaces/ws-1/runs/run_api_1", "")
	var snapshotView product.M1WorkbenchView
	if err := json.NewDecoder(snapshotResponse.Body).Decode(&snapshotView); err != nil {
		t.Fatal(err)
	}
	if *snapshotView.Result != *messageView.Result {
		t.Fatalf("HTTP snapshot drifted: %+v vs %+v", snapshotView.Result, messageView.Result)
	}

	streamResponse := performM1Request(router, http.MethodGet, "/api/workspaces/ws-1/runs/run_api_1/stream", "")
	if streamResponse.Header.Get("Content-Type") != "text/event-stream; charset=utf-8" {
		t.Fatalf("stream content type = %q", streamResponse.Header.Get("Content-Type"))
	}
	scanner := bufio.NewScanner(streamResponse.Body)
	var dataLine string
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "data: ") {
			dataLine = strings.TrimPrefix(scanner.Text(), "data: ")
		}
	}
	var event product.M1StreamEvent
	if err := json.Unmarshal([]byte(dataLine), &event); err != nil {
		t.Fatal(err)
	}
	if event.View.Result == nil || *event.View.Result != *messageView.Result || event.EventID != "run_api_1:1" {
		t.Fatalf("SSE event drifted: %+v", event)
	}
}

func TestM1RouterDoesNotMountActionHistoryResumeOrReplay(t *testing.T) {
	repository := facts.NewQueryMemoryRepository()
	router := httpapi.NewM1Router(httpapi.M1Dependencies{
		Projection: product.NewM1Projection(repository),
		Readiness:  func(context.Context) (bool, string) { return true, "ready" },
	})
	for _, target := range []string{
		"/api/workspaces/ws-1/agent/actions",
		"/api/workspaces/ws-1/runs/run-1/resume",
		"/api/workspaces/ws-1/runs/run-1/replay",
		"/api/workspaces/ws-1/history",
	} {
		response := performM1Request(router, http.MethodGet, target, "")
		if response.StatusCode != http.StatusNotFound {
			t.Fatalf("%s status = %d", target, response.StatusCode)
		}
	}
}

func TestM1RouterReadinessAndStrictInput(t *testing.T) {
	repository := facts.NewQueryMemoryRepository()
	router := httpapi.NewM1Router(httpapi.M1Dependencies{
		Projection: product.NewM1Projection(repository),
		Readiness:  func(context.Context) (bool, string) { return false, "database_not_ready" },
	})
	readyResponse := performM1Request(router, http.MethodGet, "/readyz", "")
	if readyResponse.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("ready status = %d", readyResponse.StatusCode)
	}
	invalidResponse := performM1Request(router, http.MethodPost, "/api/workspaces/ws-1/messages", `{"conversation_id":"c","actor_id":"a","content":"query","unexpected":true}`)
	if invalidResponse.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid input status = %d", invalidResponse.StatusCode)
	}
}

func performM1Request(handler http.Handler, method, target, body string) *http.Response {
	request := httptest.NewRequest(method, target, strings.NewReader(body))
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder.Result()
}
