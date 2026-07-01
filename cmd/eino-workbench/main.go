package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"agent-platform-eino/internal/einoapp/bootstrap"
	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/httpapi"
	"agent-platform-eino/internal/einoapp/llm"
	"agent-platform-eino/internal/einoapp/product"
	"agent-platform-eino/internal/einoapp/store/sqlite"
)

func main() {
	// 服务启动只读取本地配置文件，避免部署参数散落到业务代码或环境变量。
	configPath := flag.String("config", "configs/eino-workbench.local.yaml", "path to eino-workbench config file")
	flag.Parse()

	cfg, err := bootstrap.LoadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Database.Driver != "sqlite" {
		log.Fatalf("unsupported database driver %q", cfg.Database.Driver)
	}
	repository, err := sqlite.Open(context.Background(), cfg.Database.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer repository.Close()
	provider, err := llmProviderFromConfig(cfg.LLM)
	if err != nil {
		log.Fatal(err)
	}
	runner := execution.NewChatModelRunner(repository, provider, llm.Config{
		Provider:          cfg.LLM.Provider,
		BaseURL:           cfg.LLM.BaseURL,
		Model:             cfg.LLM.Model,
		ModelLabel:        cfg.LLM.ModelLabel,
		TimeoutMillis:     cfg.LLM.TimeoutMillis,
		NetworkSafety:     llm.NetworkSafety(cfg.LLM.NetworkSafety),
		CredentialBinding: llm.CredentialBinding(cfg.LLM.CredentialBinding),
	}, execution.ChatModelRunnerConfig{})
	registry, err := capabilityRegistryFromConfig(cfg.Capabilities)
	if err != nil {
		log.Fatal(err)
	}
	toolRunner := execution.NewToolLoopRunner(
		repository,
		registry,
		capabilities.NewMockProvider("config-mock", registry.List()),
		execution.ToolLoopRunnerConfig{},
	)

	// 服务路径使用 SQLite Product Facts，确保 Workbench、Action API、Replay 和 SSE 同源。
	deps := httpapi.Dependencies{
		Projection: product.NewFactsProjection(repository),
		Commands:   execution.NewToolRunnerCommandsWithRegistry(repository, runner, toolRunner, registry),
	}
	server := &http.Server{
		Addr:         cfg.Server.Addr,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		Handler:      httpapi.NewRouter(deps),
	}

	log.Printf("eino-workbench listening on http://%s", cfg.Server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// llmProviderFromConfig 创建模型 provider；真实密钥只传入 provider 私有边界，不进入 llm.Config。
func llmProviderFromConfig(cfg bootstrap.LLMConfig) (llm.Provider, error) {
	switch cfg.Provider {
	case "mock":
		return llm.NewMockProvider("已收到请求。"), nil
	case "openai_compatible":
		return llm.NewOpenAICompatibleProvider(llm.OpenAICompatibleProviderConfig{APIKey: cfg.APIKey}), nil
	default:
		return nil, fmt.Errorf("unsupported llm provider %q", cfg.Provider)
	}
}

// capabilityRegistryFromConfig 从配置文件注册能力，避免在启动路径硬编码工具名或业务 provider。
func capabilityRegistryFromConfig(configs []bootstrap.CapabilityConfig) (*capabilities.Registry, error) {
	registry := capabilities.NewRegistry()
	for _, cfg := range configs {
		if err := registry.Register(capabilities.Capability{
			ID:                      cfg.ID,
			ProviderID:              cfg.ProviderID,
			ToolName:                cfg.ToolName,
			DisplayName:             cfg.DisplayName,
			Description:             cfg.Description,
			ResultSchema:            cfg.ResultSchema,
			RiskLevel:               capabilities.RiskLevel(cfg.RiskLevel),
			SideEffect:              capabilities.SideEffect(cfg.SideEffect),
			PolicyRef:               cfg.PolicyRef,
			PermissionScope:         capabilities.PermissionScope(cfg.PermissionScope),
			CredentialBindingPolicy: capabilities.CredentialBindingPolicy(cfg.CredentialBindingPolicy),
			ConnectorID:             cfg.ConnectorID,
			ApprovalRequired:        cfg.ApprovalRequired,
			IdempotencyRequired:     cfg.IdempotencyRequired,
			Timeout:                 cfg.Timeout,
		}); err != nil {
			return nil, err
		}
	}
	return registry, nil
}
