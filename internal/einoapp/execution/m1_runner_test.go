package execution_test

import (
	"context"
	"errors"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
)

type runnerClock struct{ now time.Time }

func (clock runnerClock) Now() time.Time { return clock.now }

type staticCandidateSource struct {
	candidate capabilities.NewVulnerabilityCandidate
}

func (source staticCandidateSource) Read(context.Context) (capabilities.NewVulnerabilityCandidate, error) {
	return source.candidate, nil
}

func TestM1RunnerUsesOneEinoToolLoopForAllFixtureCases(t *testing.T) {
	for _, fixtureCase := range []capabilities.FixtureCase{capabilities.FixtureResolved, capabilities.FixtureEmpty, capabilities.FixtureFailed} {
		t.Run(string(fixtureCase), func(t *testing.T) {
			repository := facts.NewQueryMemoryRepository()
			registry := capabilities.NewRegistry()
			capability := capabilities.NewVulnerabilityReadCapability(5 * time.Second)
			if err := registry.Register(capability); err != nil {
				t.Fatal(err)
			}
			source, err := capabilities.LoadNewVulnerabilityFixture(fixturePath(t), fixtureCase)
			if err != nil {
				t.Fatal(err)
			}
			now := time.Date(2026, 7, 27, 9, 30, 0, 0, time.UTC)
			freshness, _ := facts.NewFreshnessPolicy("m1.v1", 30*time.Minute)
			runner, err := execution.NewM1QueryRunner(execution.M1QueryRunnerConfig{
				Repository: repository, Registry: registry, Source: source,
				WorkspaceID: "ws-1", ConversationID: "conversation-1", ActorID: "actor-1",
				RunID: "run-" + string(fixtureCase), ToolCallID: "call-1", QuerySequence: 1,
				Clock: runnerClock{now: now}, Freshness: freshness,
			})
			if err != nil {
				t.Fatal(err)
			}
			// 查询文本刻意不包含关键词，证明能力选择来自 Eino tool call 与唯一 registry，而非字符串分支。
			if err := runner.Run(context.Background(), "请执行当前允许的只读检查"); err != nil {
				t.Fatal(err)
			}
			aggregate, err := repository.GetQuery(context.Background(), "ws-1", "run-"+string(fixtureCase))
			if err != nil {
				t.Fatal(err)
			}
			switch fixtureCase {
			case capabilities.FixtureResolved:
				if aggregate.Status() != facts.QueryFactsSucceeded || aggregate.Result().Count() != 3 || aggregate.Result().Items()[0].Ref() != "item_alpha" {
					t.Fatal("resolved fixture did not preserve stable result")
				}
			case capabilities.FixtureEmpty:
				if aggregate.Status() != facts.QueryFactsSucceeded || aggregate.Result().Status() != facts.QueryResultEmpty || aggregate.Result().Summary() != "没有待派发漏洞" {
					t.Fatal("empty fixture was not a successful empty result")
				}
			case capabilities.FixtureFailed:
				if aggregate.Status() != facts.QueryFactsFailed || aggregate.HasStructuredResult() || aggregate.SafeError().Message != facts.FailedQueryMessage {
					t.Fatal("failed fixture was confused with empty")
				}
			}
			counters := runner.Counters()
			if counters.MockModelCalls != 1 || counters.ToolCalls != 1 || counters.RealLLMCalls != 0 || counters.RealFOBrainCalls != 0 || counters.ExternalMutations != 0 {
				t.Fatalf("call counters = %+v", counters)
			}
			if len(registry.List()) != 1 || registry.List()[0].ID != facts.NewVulnerabilityCapabilityID {
				t.Fatalf("registry = %+v", registry.List())
			}
		})
	}
}

func TestM1SafetyGateRejectsMalformedSecretAndProviderLocatorBeforeFacts(t *testing.T) {
	now := time.Date(2026, 7, 27, 9, 30, 0, 0, time.UTC)
	validItems := []capabilities.NewVulnerabilityItem{{SnapshotItemRef: "item_alpha", DisplayLabel: "新增漏洞 A", DiscoveredAt: now}}
	tests := []struct {
		name      string
		candidate capabilities.NewVulnerabilityCandidate
	}{
		{name: "malformed", candidate: capabilities.NewVulnerabilityCandidate{Status: "unknown", Summary: "查询完成", ObservedAt: now, Items: validItems}},
		{name: "secret", candidate: capabilities.NewVulnerabilityCandidate{Status: "resolved", Summary: "Authorization: Bearer hidden", ObservedAt: now, Items: validItems}},
		{name: "provider locator", candidate: capabilities.NewVulnerabilityCandidate{Status: "resolved", Summary: "查询完成", ObservedAt: now, Items: validItems, ProviderLocator: "https://internal.example.invalid"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := facts.NewQueryMemoryRepository()
			registry := capabilities.NewRegistry()
			if err := registry.Register(capabilities.NewVulnerabilityReadCapability(time.Second)); err != nil {
				t.Fatal(err)
			}
			freshness, _ := facts.NewFreshnessPolicy("m1.v1", 30*time.Minute)
			runner, err := execution.NewM1QueryRunner(execution.M1QueryRunnerConfig{
				Repository: repository, Registry: registry, Source: staticCandidateSource{candidate: test.candidate},
				WorkspaceID: "ws-1", ConversationID: "conversation-1", ActorID: "actor-1", RunID: "run-rejected", ToolCallID: "call-1", QuerySequence: 1,
				Clock: runnerClock{now: now}, Freshness: freshness,
			})
			if err != nil {
				t.Fatal(err)
			}
			if err := runner.Run(context.Background(), "执行只读检查"); !errors.Is(err, facts.ErrUnsafeFactMaterial) {
				t.Fatalf("Run err = %v", err)
			}
			if _, err := repository.GetQuery(context.Background(), "ws-1", "run-rejected"); !errors.Is(err, facts.ErrNotFound) {
				t.Fatalf("unsafe candidate reached facts: %v", err)
			}
		})
	}
}

func fixturePath(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller unavailable")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
	return filepath.Join(root, "docs", "fixtures", "new-vulnerability-walking-skeleton.v1.json")
}
