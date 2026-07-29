package product_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

type projectionClock struct{ now time.Time }

func (clock projectionClock) Now() time.Time { return clock.now }

func TestM1ProjectionUsesOneFactsAggregateForViewSnapshotAndSSE(t *testing.T) {
	repository := facts.NewQueryMemoryRepository()
	aggregate := projectionSucceededFacts(t, facts.QueryResultResolved, "发现 2 条新增漏洞")
	if err := repository.CommitQuery(context.Background(), aggregate); err != nil {
		t.Fatal(err)
	}
	projection := product.NewM1Projection(repository)
	current, err := projection.WorkbenchView(context.Background(), "ws-1")
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := projection.RunSnapshot(context.Background(), "ws-1", "run-1")
	if err != nil {
		t.Fatal(err)
	}
	events, err := projection.StreamEvents(context.Background(), "ws-1", "run-1")
	if err != nil {
		t.Fatal(err)
	}
	if current.Result == nil || snapshot.Result == nil || len(events) != 1 {
		t.Fatal("projection omitted M1 result")
	}
	if *current.Result != *snapshot.Result || *snapshot.Result != *events[0].View.Result {
		t.Fatalf("product exits drifted: current=%+v snapshot=%+v event=%+v", current.Result, snapshot.Result, events[0].View.Result)
	}
	if events[0].EventID != aggregate.Events()[0].EventID() || events[0].Sequence != aggregate.Events()[0].Sequence() {
		t.Fatal("SSE cursor was rebuilt instead of using persisted FactEvent")
	}
	encoded, _ := json.Marshal(events)
	for _, forbidden := range [][]byte{[]byte("result_ref"), []byte("snapshot_item_ref"), []byte("structured_result"), []byte("opaque")} {
		if bytes.Contains(bytes.ToLower(encoded), forbidden) {
			t.Fatalf("product projection leaked %q: %s", forbidden, encoded)
		}
	}
}

func TestM1ProjectionSeparatesEmptyAndFailedCopy(t *testing.T) {
	tests := []struct {
		name, wantStatus, wantMessage string
		aggregate                     facts.QueryFacts
	}{
		{name: "empty", wantStatus: "empty", wantMessage: "没有待派发漏洞", aggregate: projectionSucceededFacts(t, facts.QueryResultEmpty, "没有待派发漏洞")},
		{name: "failed", wantStatus: "failed", wantMessage: facts.FailedQueryMessage, aggregate: projectionFailedFacts(t)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := facts.NewQueryMemoryRepository()
			if err := repository.CommitQuery(context.Background(), test.aggregate); err != nil {
				t.Fatal(err)
			}
			view, err := product.NewM1Projection(repository).RunSnapshot(context.Background(), "ws-1", test.aggregate.RunID())
			if err != nil {
				t.Fatal(err)
			}
			if view.Status != test.wantStatus || view.Messages[len(view.Messages)-1].Content != test.wantMessage {
				t.Fatalf("view = %+v", view)
			}
			if test.wantStatus == "failed" && (view.Result != nil || view.SafeError != facts.FailedQueryMessage) {
				t.Fatal("failed projection created pseudo empty card")
			}
		})
	}
}

func projectionSucceededFacts(t *testing.T, status facts.QueryResultStatus, summary string) facts.QueryFacts {
	t.Helper()
	now := time.Date(2026, 7, 27, 10, 0, 0, 0, time.UTC)
	items := []facts.SnapshotItem{}
	if status == facts.QueryResultResolved {
		items = []facts.SnapshotItem{
			facts.MustSnapshotItem("item_alpha", "新增漏洞 A", now),
			facts.MustSnapshotItem("item_beta", "新增漏洞 B", now.Add(-time.Minute)),
		}
	}
	result, err := facts.NewStructuredResult(facts.StructuredResultInput{InternalResultRef: "result_ref:run-1", Status: status, Summary: summary, ObservedAt: now, Items: items})
	if err != nil {
		t.Fatal(err)
	}
	policy, _ := facts.NewFreshnessPolicy("m1.v1", 30*time.Minute)
	snapshot, err := facts.NewQueryResultSnapshot(facts.QueryResultSnapshotInput{
		SnapshotID: "snapshot-1", WorkspaceID: "ws-1", ConversationID: "conversation-1", ActorID: "actor-1", RunID: "run-1", ToolCallID: "call-1",
		Result: result, Clock: projectionClock{now: now}, Freshness: policy,
	})
	if err != nil {
		t.Fatal(err)
	}
	aggregate, err := facts.NewSucceededQueryFacts(facts.SucceededQueryFactsInput{
		WorkspaceID: "ws-1", ConversationID: "conversation-1", ActorID: "actor-1", RunID: "run-1", QuerySequence: 1, CreatedAt: now, Result: result, Snapshot: snapshot,
	})
	if err != nil {
		t.Fatal(err)
	}
	return aggregate
}

func projectionFailedFacts(t *testing.T) facts.QueryFacts {
	t.Helper()
	aggregate, err := facts.NewFailedQueryFacts(facts.FailedQueryFactsInput{
		WorkspaceID: "ws-1", ConversationID: "conversation-1", ActorID: "actor-1", RunID: "run-failed", QuerySequence: 1,
		CreatedAt: time.Date(2026, 7, 27, 10, 0, 0, 0, time.UTC), Error: facts.SafeQueryError{Code: facts.ErrorCodeFixtureReadFailed, Message: facts.FailedQueryMessage},
	})
	if err != nil {
		t.Fatal(err)
	}
	return aggregate
}
