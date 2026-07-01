package execution

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
)

// ErrRunNotFound 表示 resume 请求引用了不存在的 run。
var ErrRunNotFound = errors.New("run not found")

// ErrCapabilityNotRegistered 表示 Action API 请求的 capability hint 未通过注册表门禁。
var ErrCapabilityNotRegistered = errors.New("capability not registered")

// AcceptedRun 是 HTTP/API 边界可以返回的最小接收结果，不包含执行事件。
type AcceptedRun struct {
	RunID  string
	Status string
}

// MessageCommand 表示 Workbench message 入口的执行命令。
type MessageCommand struct {
	WorkspaceID     string
	Message         string
	ClientRequestID string
	RunID           string
}

// ActionCommand 表示外部 Action API 入口的执行命令。
type ActionCommand struct {
	WorkspaceID     string
	ActionID        string
	ClientRequestID string
	CapabilityHint  string
	InputText       string
	Attachments     []Attachment
	Context         RequestContext
}

// ResumeCommand 表示 approval/clarification 的恢复请求。
type ResumeCommand struct {
	WorkspaceID     string
	RunID           string
	ResumeRef       string
	ClientRequestID string
	Decision        string
	SelectedRefs    []string
	FreeText        string
	Comment         string
}

// Attachment 是 Action API 的安全附件摘要，不包含文件原文。
type Attachment struct {
	AttachmentRef string
	MediaType     string
	SafeName      string
}

// RequestContext 是请求侧安全上下文摘要，可进入后续 Product Facts。
type RequestContext struct {
	Timezone      string
	Locale        string
	SafeUserLabel string
}

// Commands 是 HTTP 层调用 execution 的唯一接口。
type Commands interface {
	StartMessage(context.Context, MessageCommand) (AcceptedRun, error)
	StartAction(context.Context, ActionCommand) (AcceptedRun, error)
	Resume(context.Context, ResumeCommand) (AcceptedRun, error)
}

// Runner 是命令层触发执行的最小接口，生产实现由 ChatModelRunner 提供。
type Runner interface {
	Run(context.Context, string) error
}

// StaticCommands 是 Phase 1/3 的最小命令实现，负责创建基础 run fact。
type StaticCommands struct {
	repository facts.Repository
	newRunID   func() (string, error)
	runner     Runner
	registry   *capabilities.Registry
}

// NewStaticCommands 创建不落 Product Facts 的静态命令替身。
func NewStaticCommands() StaticCommands {
	return StaticCommands{newRunID: randomRunID}
}

// NewFactCommands 创建会写入 Product Facts 的命令实现。
func NewFactCommands(repository facts.Repository) StaticCommands {
	return StaticCommands{
		repository: repository,
		newRunID:   randomRunID,
	}
}

// NewRunnerCommands 创建会写入 Product Facts 并触发 ChatModelRunner 的命令实现。
func NewRunnerCommands(repository facts.Repository, runner Runner) StaticCommands {
	return StaticCommands{
		repository: repository,
		newRunID:   randomRunID,
		runner:     runner,
	}
}

// NewRunnerCommandsWithRegistry 创建带 capability registry 的命令实现，用于 Action API 选择门禁。
func NewRunnerCommandsWithRegistry(repository facts.Repository, runner Runner, registry *capabilities.Registry) StaticCommands {
	return StaticCommands{
		repository: repository,
		newRunID:   randomRunID,
		runner:     runner,
		registry:   registry,
	}
}

// StartMessage 接收 Workbench 消息并创建 run。
func (commands StaticCommands) StartMessage(ctx context.Context, command MessageCommand) (AcceptedRun, error) {
	accepted, err := commands.acceptRun(ctx, command.WorkspaceID, command.RunID, command.Message)
	if err != nil {
		return AcceptedRun{}, err
	}
	return commands.runIfConfigured(ctx, accepted)
}

// StartAction 接收 Action API 请求并创建 run。
func (commands StaticCommands) StartAction(ctx context.Context, command ActionCommand) (AcceptedRun, error) {
	selection := SelectCapability(commands.registry, SelectionRequest{
		InputText:      command.InputText,
		CapabilityHint: command.CapabilityHint,
		ProductAction:  command.ActionID,
	})
	if selection.Mode == SelectionModeRejected {
		return AcceptedRun{}, ErrCapabilityNotRegistered
	}
	accepted, err := commands.acceptRun(ctx, command.WorkspaceID, "", command.InputText)
	if err != nil {
		return AcceptedRun{}, err
	}
	if err := commands.recordSelection(ctx, accepted.RunID, selection); err != nil {
		return AcceptedRun{}, err
	}
	if selection.RequiresApproval {
		return accepted, nil
	}
	return commands.runIfConfigured(ctx, accepted)
}

// Resume 校验 run 存在后接收恢复请求；完整 HITL 执行在后续 Phase 接入。
func (commands StaticCommands) Resume(ctx context.Context, command ResumeCommand) (AcceptedRun, error) {
	if commands.repository != nil {
		if _, err := commands.repository.GetRun(ctx, command.RunID); err != nil {
			if errors.Is(err, facts.ErrNotFound) {
				return AcceptedRun{}, ErrRunNotFound
			}
			return AcceptedRun{}, err
		}
	}
	return AcceptedRun{RunID: command.RunID, Status: "accepted"}, nil
}

// acceptRun 统一创建 run，保证 Message 和 Action 入口共享 Product Facts。
func (commands StaticCommands) acceptRun(ctx context.Context, workspaceID string, requestedRunID string, inputText string) (AcceptedRun, error) {
	runID := strings.TrimSpace(requestedRunID)
	if runID == "" {
		generated, err := commands.newRunID()
		if err != nil {
			return AcceptedRun{}, err
		}
		runID = generated
	}
	if commands.repository != nil {
		now := time.Now().UTC()
		if err := commands.repository.CreateRun(ctx, facts.Run{
			RunID:       runID,
			WorkspaceID: workspaceID,
			Status:      facts.RunStatusCreated,
			CreatedAt:   now,
			UpdatedAt:   now,
		}); err != nil {
			return AcceptedRun{}, err
		}
		if strings.TrimSpace(inputText) != "" {
			if err := commands.repository.AppendTurn(ctx, facts.Turn{
				TurnID:    runID + ":user:1",
				RunID:     runID,
				Role:      facts.TurnRoleUser,
				Content:   inputText,
				Sequence:  1,
				CreatedAt: now,
			}); err != nil {
				return AcceptedRun{}, err
			}
			if err := commands.repository.AppendAuditEvent(ctx, facts.AuditEvent{
				AuditID:     runID + ":audit:message:1",
				RunID:       runID,
				EventType:   "message",
				SafeSummary: "user message accepted",
				Actor:       "user",
				CreatedAt:   now,
			}); err != nil {
				return AcceptedRun{}, err
			}
		}
	}
	return AcceptedRun{RunID: runID, Status: "accepted"}, nil
}

// recordSelection 将 capability 选择结果写成安全审计事实，后续 tool execution 只能读取该事实链路。
func (commands StaticCommands) recordSelection(ctx context.Context, runID string, selection SelectionResult) error {
	if commands.repository == nil || selection.Mode != SelectionModeCapability {
		return nil
	}
	return commands.repository.AppendAuditEvent(ctx, facts.AuditEvent{
		AuditID:     runID + ":audit:capability:1",
		RunID:       runID,
		EventType:   "tool",
		SafeSummary: "capability selected: " + selection.CapabilityID + "; policy: " + selection.PolicyReason,
		Actor:       "system",
		CreatedAt:   time.Now().UTC(),
	})
}

// runIfConfigured 在生产路径触发 runner；测试替身可不配置 runner。
func (commands StaticCommands) runIfConfigured(ctx context.Context, accepted AcceptedRun) (AcceptedRun, error) {
	if commands.runner == nil {
		return accepted, nil
	}
	if err := commands.runner.Run(ctx, accepted.RunID); err != nil {
		return AcceptedRun{}, err
	}
	return accepted, nil
}

// randomRunID 生成不含 workspace/action 信息的 opaque run id。
func randomRunID() (string, error) {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return "run_" + hex.EncodeToString(bytes[:]), nil
}
