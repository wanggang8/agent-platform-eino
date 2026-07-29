package execution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync/atomic"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// M1CallCounters 为验收证明 mock loop 只调用一个工具且真实外部调用始终为零。
type M1CallCounters struct {
	MockModelCalls    int64
	ToolCalls         int64
	RealLLMCalls      int64
	RealFOBrainCalls  int64
	ExternalMutations int64
}

type m1CounterState struct {
	mockModel atomic.Int64
	tool      atomic.Int64
}

type M1QueryRunnerConfig struct {
	Repository                                              facts.QueryRepository
	Registry                                                *capabilities.Registry
	Source                                                  capabilities.NewVulnerabilitySource
	WorkspaceID, ConversationID, ActorID, RunID, ToolCallID string
	QuerySequence                                           int64
	Clock                                                   facts.Clock
	Freshness                                               facts.FreshnessPolicy
}

// M1QueryRunner 只装配 mock model、唯一 fixture capability 和 Product Facts 端口。
type M1QueryRunner struct {
	config   M1QueryRunnerConfig
	counters *m1CounterState
}

func NewM1QueryRunner(config M1QueryRunnerConfig) (*M1QueryRunner, error) {
	if config.Repository == nil || config.Registry == nil || config.Source == nil || config.Clock == nil || config.QuerySequence < 1 {
		return nil, errors.New("invalid M1 query runner dependencies")
	}
	capabilityList := config.Registry.List()
	if len(capabilityList) != 1 {
		return nil, errors.New("M1 registry must contain exactly one capability")
	}
	capability := capabilityList[0]
	if capability.ID != facts.NewVulnerabilityCapabilityID || capability.RiskLevel != capabilities.RiskNone || capability.SideEffect != capabilities.SideEffectNone || capability.ApprovalRequired {
		return nil, errors.New("M1 capability metadata violates read-only boundary")
	}
	for _, value := range []string{config.WorkspaceID, config.ConversationID, config.ActorID, config.RunID, config.ToolCallID} {
		if strings.TrimSpace(value) == "" || facts.ContainsUnsafeMaterial(value) {
			return nil, errors.New("invalid M1 query identity")
		}
	}
	return &M1QueryRunner{config: config, counters: &m1CounterState{}}, nil
}

// Run 让 Eino ChatModelAgent 绑定 registry 工具并执行一次确定性 tool call。
func (runner *M1QueryRunner) Run(ctx context.Context, query string) error {
	if strings.TrimSpace(query) == "" || facts.ContainsUnsafeMaterial(query) {
		return facts.ErrUnsafeFactMaterial
	}
	capability := runner.config.Registry.List()[0]
	tool := &m1QueryTool{config: runner.config, capability: capability, counters: runner.counters}
	mockModel := &m1ToolCallingModel{counters: runner.counters}
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name: "m1-query-agent", Description: "M1 fixture-only read agent", Model: mockModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{Tools: []einotool.BaseTool{tool}},
			ReturnDirectly:  map[string]bool{capability.ToolName: true},
		},
		MaxIterations: 2,
	})
	if err != nil {
		return err
	}
	iterator := adk.NewRunner(ctx, adk.RunnerConfig{Agent: agent, EnableStreaming: false}).Query(ctx, query)
	for {
		event, ok := iterator.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			return event.Err
		}
	}
	return nil
}

func (runner *M1QueryRunner) Counters() M1CallCounters {
	return M1CallCounters{MockModelCalls: runner.counters.mockModel.Load(), ToolCalls: runner.counters.tool.Load()}
}

type m1ToolCallingModel struct {
	tools    []*schema.ToolInfo
	counters *m1CounterState
}

var _ model.ToolCallingChatModel = (*m1ToolCallingModel)(nil)

func (chatModel *m1ToolCallingModel) WithTools(tools []*schema.ToolInfo) (model.ToolCallingChatModel, error) {
	if len(tools) != 1 || tools[0] == nil || strings.TrimSpace(tools[0].Name) == "" {
		return nil, errors.New("M1 mock model requires exactly one tool")
	}
	return &m1ToolCallingModel{tools: append([]*schema.ToolInfo(nil), tools...), counters: chatModel.counters}, nil
}

func (chatModel *m1ToolCallingModel) Generate(_ context.Context, _ []*schema.Message, options ...model.Option) (*schema.Message, error) {
	chatModel.counters.mockModel.Add(1)
	tools := model.GetCommonOptions(nil, options...).Tools
	if len(tools) == 0 {
		tools = chatModel.tools
	}
	if len(tools) != 1 {
		return nil, errors.New("M1 mock model has no bound tool")
	}
	return schema.AssistantMessage("", []schema.ToolCall{{
		ID: "m1-tool-call", Type: "function",
		Function: schema.FunctionCall{Name: tools[0].Name, Arguments: `{}`},
	}}), nil
}

func (chatModel *m1ToolCallingModel) Stream(context.Context, []*schema.Message, ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	return nil, errors.New("M1 mock model streaming is disabled")
}

type m1QueryTool struct {
	config     M1QueryRunnerConfig
	capability capabilities.Capability
	counters   *m1CounterState
}

var _ einotool.InvokableTool = (*m1QueryTool)(nil)

func (tool *m1QueryTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return capabilities.ToolInfoFromCapability(ctx, tool.capability)
}

func (tool *m1QueryTool) InvokableRun(ctx context.Context, arguments string, _ ...einotool.Option) (string, error) {
	tool.counters.tool.Add(1)
	var input map[string]any
	decoder := json.NewDecoder(strings.NewReader(arguments))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil || len(input) != 0 {
		return "", facts.ErrUnsafeFactMaterial
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return "", facts.ErrUnsafeFactMaterial
	}
	candidate, err := tool.config.Source.Read(ctx)
	if err != nil {
		return "", err
	}
	aggregate, output, err := tool.approveAndBuild(candidate)
	if err != nil {
		return "", err
	}
	if err := tool.config.Repository.CommitQuery(ctx, aggregate); err != nil {
		return "", err
	}
	return output, nil
}

func (tool *m1QueryTool) approveAndBuild(candidate capabilities.NewVulnerabilityCandidate) (facts.QueryFacts, string, error) {
	now := tool.config.Clock.Now().UTC()
	if now.IsZero() || candidate.ProviderLocator != "" || len(candidate.RawPayload) != 0 {
		return facts.QueryFacts{}, "", facts.ErrUnsafeFactMaterial
	}
	if candidate.FailureCode != "" || candidate.FailureMessage != "" {
		if candidate.FailureCode != facts.ErrorCodeFixtureReadFailed || candidate.FailureMessage != facts.FailedQueryMessage || candidate.Status != "" || len(candidate.Items) != 0 {
			return facts.QueryFacts{}, "", facts.ErrUnsafeFactMaterial
		}
		aggregate, err := facts.NewFailedQueryFacts(facts.FailedQueryFactsInput{
			WorkspaceID: tool.config.WorkspaceID, ConversationID: tool.config.ConversationID, ActorID: tool.config.ActorID,
			RunID: tool.config.RunID, QuerySequence: tool.config.QuerySequence, CreatedAt: now,
			Error: facts.SafeQueryError{Code: candidate.FailureCode, Message: candidate.FailureMessage},
		})
		return aggregate, `{"code":"fixture_read_failed","message":"无法读取漏洞事实。此次请求不是空结果。"}`, err
	}
	status := facts.QueryResultStatus(candidate.Status)
	if status != facts.QueryResultResolved && status != facts.QueryResultEmpty {
		return facts.QueryFacts{}, "", facts.ErrUnsafeFactMaterial
	}
	items := make([]facts.SnapshotItem, 0, len(candidate.Items))
	for _, candidateItem := range candidate.Items {
		item, err := facts.NewSnapshotItem(candidateItem.SnapshotItemRef, candidateItem.DisplayLabel, candidateItem.DiscoveredAt)
		if err != nil {
			return facts.QueryFacts{}, "", facts.ErrUnsafeFactMaterial
		}
		items = append(items, item)
	}
	result, err := facts.NewStructuredResult(facts.StructuredResultInput{
		InternalResultRef: "result_ref:" + tool.config.RunID, Status: status,
		Summary: candidate.Summary, ObservedAt: candidate.ObservedAt, Items: items,
	})
	if err != nil {
		return facts.QueryFacts{}, "", facts.ErrUnsafeFactMaterial
	}
	snapshot, err := facts.NewQueryResultSnapshot(facts.QueryResultSnapshotInput{
		SnapshotID: "snapshot-" + tool.config.RunID, WorkspaceID: tool.config.WorkspaceID,
		ConversationID: tool.config.ConversationID, ActorID: tool.config.ActorID, RunID: tool.config.RunID,
		ToolCallID: tool.config.ToolCallID, Result: result, Clock: tool.config.Clock, Freshness: tool.config.Freshness,
	})
	if err != nil {
		return facts.QueryFacts{}, "", err
	}
	aggregate, err := facts.NewSucceededQueryFacts(facts.SucceededQueryFactsInput{
		WorkspaceID: tool.config.WorkspaceID, ConversationID: tool.config.ConversationID, ActorID: tool.config.ActorID,
		RunID: tool.config.RunID, QuerySequence: tool.config.QuerySequence, CreatedAt: now, Result: result, Snapshot: snapshot,
	})
	if err != nil {
		return facts.QueryFacts{}, "", err
	}
	return aggregate, string(result.SafeJSON()), nil
}

func (counters M1CallCounters) String() string {
	return fmt.Sprintf("mock_model=%d tool=%d real_llm=%d real_fobrain=%d mutation=%d", counters.MockModelCalls, counters.ToolCalls, counters.RealLLMCalls, counters.RealFOBrainCalls, counters.ExternalMutations)
}
