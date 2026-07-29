package execution

import (
	"context"
	"errors"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/facts"
)

// M1MessageCommand 是 HTTP 边界传入 execution 的封闭命令对象。
type M1MessageCommand struct {
	WorkspaceID, ConversationID, ActorID, RunID, ToolCallID, Content string
}

type M1MessageExecutor interface {
	Execute(ctx context.Context, command M1MessageCommand) error
}

type M1QueryServiceConfig struct {
	Repository facts.QueryRepository
	Registry   *capabilities.Registry
	Source     capabilities.NewVulnerabilitySource
	Clock      facts.Clock
	Freshness  facts.FreshnessPolicy
}

// M1QueryService 为每条消息创建隔离 runner，并从当前事实分配查询序号。
type M1QueryService struct{ config M1QueryServiceConfig }

func NewM1QueryService(config M1QueryServiceConfig) (M1QueryService, error) {
	if config.Repository == nil || config.Registry == nil || config.Source == nil || config.Clock == nil {
		return M1QueryService{}, errors.New("invalid M1 query service dependencies")
	}
	return M1QueryService{config: config}, nil
}

func (service M1QueryService) Execute(ctx context.Context, command M1MessageCommand) error {
	sequence := int64(1)
	latest, err := service.config.Repository.LatestQuery(ctx, command.WorkspaceID)
	if err == nil {
		sequence = latest.QuerySequence() + 1
	} else if !errors.Is(err, facts.ErrNotFound) {
		return err
	}
	runner, err := NewM1QueryRunner(M1QueryRunnerConfig{
		Repository: service.config.Repository, Registry: service.config.Registry, Source: service.config.Source,
		WorkspaceID: command.WorkspaceID, ConversationID: command.ConversationID, ActorID: command.ActorID,
		RunID: command.RunID, ToolCallID: command.ToolCallID, QuerySequence: sequence,
		Clock: service.config.Clock, Freshness: service.config.Freshness,
	})
	if err != nil {
		return err
	}
	return runner.Run(ctx, command.Content)
}
