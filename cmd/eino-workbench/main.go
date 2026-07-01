package main

import (
	"context"
	"flag"
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
	if cfg.LLM.Provider != "mock" {
		log.Fatalf("unsupported llm provider %q before Phase 4", cfg.LLM.Provider)
	}
	runner := execution.NewChatModelRunner(repository, llm.NewMockProvider("已收到请求。"), llm.Config{
		Provider: cfg.LLM.Provider,
		Model:    cfg.LLM.Model,
	}, execution.ChatModelRunnerConfig{})
	registry, err := capabilityRegistryFromConfig(cfg.Capabilities)
	if err != nil {
		log.Fatal(err)
	}

	// 服务路径使用 SQLite Product Facts，确保 Workbench、Action API、Replay 和 SSE 同源。
	deps := httpapi.Dependencies{
		Projection: product.NewFactsProjection(repository),
		Commands:   execution.NewRunnerCommandsWithRegistry(repository, runner, registry),
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

// capabilityRegistryFromConfig 从配置文件注册能力，避免在启动路径硬编码工具名或业务 provider。
func capabilityRegistryFromConfig(configs []bootstrap.CapabilityConfig) (*capabilities.Registry, error) {
	registry := capabilities.NewRegistry()
	for _, cfg := range configs {
		if err := registry.Register(capabilities.Capability{
			ID:               cfg.ID,
			ProviderID:       cfg.ProviderID,
			ToolName:         cfg.ToolName,
			DisplayName:      cfg.DisplayName,
			Description:      cfg.Description,
			ResultSchema:     cfg.ResultSchema,
			RiskLevel:        capabilities.RiskLevel(cfg.RiskLevel),
			ApprovalRequired: cfg.ApprovalRequired,
			Timeout:          cfg.Timeout,
		}); err != nil {
			return nil, err
		}
	}
	return registry, nil
}
