package execution

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ErrCapabilityRequiresApproval 表示能力需要 HITL 审批，Phase 4 工具循环不能直接执行。
var ErrCapabilityRequiresApproval = errors.New("capability requires approval")

// ToolLoopRunnerConfig 提供可替换时间源，保证工具事实测试稳定。
type ToolLoopRunnerConfig struct {
	Now func() time.Time
}

// ToolLoopRunner 执行只读 capability 的 Eino tool adapter，并把结果写入 Product Facts。
type ToolLoopRunner struct {
	repository facts.Repository
	registry   *capabilities.Registry
	invoker    capabilities.Invoker
	now        func() time.Time
}

// NewToolLoopRunner 创建 Phase 4 mock tool loop runner。
func NewToolLoopRunner(repository facts.Repository, registry *capabilities.Registry, invoker capabilities.Invoker, config ToolLoopRunnerConfig) ToolLoopRunner {
	now := config.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return ToolLoopRunner{repository: repository, registry: registry, invoker: invoker, now: now}
}

// RunCapability 执行显式选择的只读能力；工具选择本身必须来自 registry/policy。
func (runner ToolLoopRunner) RunCapability(ctx context.Context, runID string, capabilityID string, inputText string) error {
	capability, ok := runner.registry.Get(capabilityID)
	if !ok {
		return ErrCapabilityNotRegistered
	}
	decision := capabilities.EvaluatePolicy(capability)
	if !decision.Allowed || decision.RequiresApproval {
		return ErrCapabilityRequiresApproval
	}
	if err := runner.repository.UpdateRunStatus(ctx, runID, facts.RunStatusRunning, "", runner.now()); err != nil {
		return err
	}

	tool := NewEinoCapabilityTool(runner.repository, capability, runner.invoker, EinoCapabilityToolConfig{
		RunID: runID,
		Now:   runner.now,
	})
	arguments, err := json.Marshal(argumentsFromCapabilityInput(capability, inputText))
	if err != nil {
		return err
	}
	if _, err := tool.InvokableRun(ctx, string(arguments)); err != nil {
		_ = runner.repository.UpdateRunStatus(ctx, runID, facts.RunStatusFailed, "工具执行失败", runner.now())
		return err
	}
	return runner.repository.UpdateRunStatus(ctx, runID, facts.RunStatusSucceeded, "", runner.now())
}

// EinoCapabilityToolConfig 绑定单次 run 的工具执行上下文。
type EinoCapabilityToolConfig struct {
	RunID string
	Now   func() time.Time
}

// EinoCapabilityTool 将 capability provider 包装为 Eino InvokableTool，并负责 Product Facts 写入。
type EinoCapabilityTool struct {
	repository facts.Repository
	capability capabilities.Capability
	invoker    capabilities.Invoker
	runID      string
	now        func() time.Time
}

var _ einotool.InvokableTool = (*EinoCapabilityTool)(nil)

// NewEinoCapabilityTool 创建单 capability 的 Eino 工具适配器。
func NewEinoCapabilityTool(repository facts.Repository, capability capabilities.Capability, invoker capabilities.Invoker, config EinoCapabilityToolConfig) *EinoCapabilityTool {
	now := config.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &EinoCapabilityTool{
		repository: repository,
		capability: capability,
		invoker:    invoker,
		runID:      config.RunID,
		now:        now,
	}
}

// Info 从 capability metadata 生成 Eino ToolInfo。
func (tool *EinoCapabilityTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return capabilities.ToolInfoFromCapability(ctx, tool.capability)
}

// InvokableRun 校验参数、调用 provider、通过 Safety Gate，并写入 ToolCall/ToolResult/audit facts。
func (tool *EinoCapabilityTool) InvokableRun(ctx context.Context, argumentsInJSON string, _ ...einotool.Option) (string, error) {
	arguments, err := decodeToolArguments(argumentsInJSON)
	if err != nil {
		return "", err
	}
	if unsafeArguments(arguments) {
		return "", facts.ErrUnsafeFactMaterial
	}

	candidate, err := tool.invoker.Invoke(ctx, capabilities.InvocationRequest{
		CapabilityID: tool.capability.ID,
		Arguments:    arguments,
	})
	if err != nil {
		return "", err
	}
	callID := tool.toolCallID()
	candidate.ResultRef = "result:" + callID
	structuredResult, err := product.NewStructuredResultSafetyGate().Approve(candidate)
	if err != nil {
		return "", err
	}

	now := tool.now()
	if err := tool.repository.AppendToolCall(ctx, facts.ToolCall{
		ToolCallID:  callID,
		RunID:       tool.runID,
		ToolID:      tool.capability.ID,
		DisplayName: tool.capability.DisplayName,
		Status:      facts.ToolCallSucceeded,
		ArgsHash:    hashArguments(argumentsInJSON),
		ArgsPreview: safeArgsPreview(arguments),
		CreatedAt:   now,
		EndedAt:     now,
	}); err != nil {
		return "", err
	}
	if err := tool.repository.AppendToolResult(ctx, facts.ToolResult{
		ResultID:         callID + ":result",
		ToolCallID:       callID,
		Status:           facts.ToolResultSucceeded,
		StructuredResult: structuredResult,
	}); err != nil {
		return "", err
	}
	if err := tool.repository.AppendAuditEvent(ctx, facts.AuditEvent{
		AuditID:     callID + ":audit",
		RunID:       tool.runID,
		EventType:   "tool",
		SafeSummary: "tool completed: " + tool.capability.ID,
		Actor:       "tool",
		CreatedAt:   now,
	}); err != nil {
		return "", err
	}
	return structuredResult.SafeSummary, nil
}

// toolCallID 生成稳定工具调用 id，便于 replay 和测试复现。
func (tool *EinoCapabilityTool) toolCallID() string {
	return fmt.Sprintf("%s:tool:%s", tool.runID, tool.capability.ID)
}

// decodeToolArguments 只接受 JSON object 参数，避免 provider 直接消费模型原始字符串。
func decodeToolArguments(argumentsInJSON string) (map[string]any, error) {
	var arguments map[string]any
	if err := json.Unmarshal([]byte(argumentsInJSON), &arguments); err != nil {
		return nil, err
	}
	if arguments == nil {
		return nil, errors.New("tool arguments must be an object")
	}
	return arguments, nil
}

// unsafeArguments 检查工具参数中的字符串值是否包含明显不安全材料。
func unsafeArguments(arguments map[string]any) bool {
	return unsafeArgumentValue(arguments)
}

// safeArgsPreview 只暴露字段存在性，不把用户原文或 provider 参数写入 facts。
func safeArgsPreview(arguments map[string]any) string {
	keys := make([]string, 0, len(arguments))
	for key := range arguments {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	if len(keys) == 0 {
		return "args=empty"
	}
	for i, key := range keys {
		keys[i] = key + "=present"
	}
	return strings.Join(keys, ",")
}

// hashArguments 生成参数哈希，供审计定位但不暴露参数明文。
func hashArguments(argumentsInJSON string) string {
	sum := sha256.Sum256([]byte(argumentsInJSON))
	return hex.EncodeToString(sum[:])
}

// argumentsFromCapabilityInput 按 capability input schema 选择参数字段，Phase 4.2 只支持单文本入口。
func argumentsFromCapabilityInput(capability capabilities.Capability, inputText string) map[string]any {
	field := firstInputField(capability)
	if field == "" {
		field = "input"
	}
	return map[string]any{field: inputText}
}

// firstInputField 返回 capability schema 中稳定排序后的第一个字段，避免固定 query 名称。
func firstInputField(capability capabilities.Capability) string {
	required := append([]string(nil), capability.InputSchema.Required...)
	sort.Strings(required)
	for _, key := range required {
		if _, ok := capability.InputSchema.Properties[key]; ok {
			return key
		}
	}
	keys := make([]string, 0, len(capability.InputSchema.Properties))
	for key := range capability.InputSchema.Properties {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	if len(keys) == 0 {
		return ""
	}
	return keys[0]
}

// unsafeArgumentValue 递归检查参数 key/value，避免模型生成的嵌套 raw 材料进入 facts 摘要。
func unsafeArgumentValue(value any) bool {
	switch typed := value.(type) {
	case string:
		return facts.ContainsUnsafeMaterial(typed)
	case map[string]any:
		for key, nested := range typed {
			if facts.ContainsUnsafeMaterial(key) || unsafeArgumentValue(nested) {
				return true
			}
		}
	case []any:
		for _, nested := range typed {
			if unsafeArgumentValue(nested) {
				return true
			}
		}
	}
	return false
}
