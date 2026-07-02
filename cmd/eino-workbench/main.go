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
	"agent-platform-eino/internal/einoapp/providers/fobrain"
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
	registry, capabilityInvoker, err := capabilityRuntimeFromConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}
	policyContexts := policyContextsFromConfig(cfg)
	toolRunner := execution.NewToolLoopRunner(
		repository,
		registry,
		capabilityInvoker,
		execution.ToolLoopRunnerConfig{PolicyContexts: policyContexts},
	)

	// 服务路径使用 SQLite Product Facts，确保 Workbench、Action API、Replay 和 SSE 同源。
	deps := httpapi.Dependencies{
		Projection: product.NewFactsProjection(repository),
		Commands:   execution.NewToolRunnerCommandsWithRegistry(repository, runner, toolRunner, registry).WithPolicyContexts(policyContexts),
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
		if err := registry.Register(capabilityFromConfig(cfg)); err != nil {
			return nil, err
		}
	}
	return registry, nil
}

// capabilityRuntimeFromConfig 构建 capability registry 和 provider invoker 路由。
// 本地配置能力与 mock MCP server 都通过 Provider/Invoker 接口接入，execution 不按工具名分支。
func capabilityRuntimeFromConfig(cfg bootstrap.Config) (*capabilities.Registry, capabilities.Invoker, error) {
	registry := capabilities.NewRegistry()
	mux := capabilities.NewInvokerMux()

	configuredCapabilities, err := capabilitiesFromConfig(cfg.Capabilities)
	if err != nil {
		return nil, nil, err
	}
	if len(configuredCapabilities) > 0 {
		localProvider := capabilities.NewMockProvider("config-mock", configuredCapabilities)
		if err := mux.RegisterProvider(registry, localProvider, localProvider); err != nil {
			return nil, nil, err
		}
	}
	for _, serverConfig := range cfg.MCPMockServers {
		provider := mcpMockProviderFromConfig(serverConfig)
		if _, err := provider.Initialize(context.Background()); err != nil {
			return nil, nil, err
		}
		if err := mux.RegisterProvider(registry, provider, provider); err != nil {
			return nil, nil, err
		}
	}
	if cfg.Fobrain.Enabled {
		provider := fobrainProviderFromConfig(cfg.Fobrain)
		if err := mux.RegisterProvider(registry, provider, provider); err != nil {
			return nil, nil, err
		}
	}
	return registry, mux, nil
}

// capabilitiesFromConfig 将 YAML capability 元数据转换为运行期契约对象。
func capabilitiesFromConfig(configs []bootstrap.CapabilityConfig) ([]capabilities.Capability, error) {
	out := make([]capabilities.Capability, 0, len(configs))
	for _, cfg := range configs {
		capability := capabilityFromConfig(cfg)
		if capability.ID == "" || capability.ProviderID == "" || capability.ToolName == "" {
			return nil, fmt.Errorf("invalid capability config %q", cfg.ID)
		}
		out = append(out, capability)
	}
	return out, nil
}

// capabilityFromConfig 保持旧测试入口和新 runtime 入口共用同一份字段映射。
func capabilityFromConfig(cfg bootstrap.CapabilityConfig) capabilities.Capability {
	return capabilities.Capability{
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
	}
}

// mcpMockProviderFromConfig 把文件化 mock MCP catalog 转成可运行 provider。
func mcpMockProviderFromConfig(config bootstrap.MCPMockServerConfig) *capabilities.MockMCPProvider {
	tools := make([]capabilities.MCPToolDefinition, 0, len(config.Tools))
	for _, tool := range config.Tools {
		tools = append(tools, capabilities.MCPToolDefinition{
			Name:         tool.Name,
			Title:        tool.Title,
			Description:  tool.Description,
			InputSchema:  mcpJSONSchemaFromConfig(tool.InputSchema),
			OutputSchema: mcpJSONSchemaFromConfig(tool.OutputSchema),
			Annotations: capabilities.MCPToolAnnotations{
				ReadOnlyHint:    tool.Annotations.ReadOnlyHint,
				DestructiveHint: tool.Annotations.DestructiveHint,
				IdempotentHint:  tool.Annotations.IdempotentHint,
			},
			ProjectPolicy: capabilities.MCPProjectPolicy{
				RiskLevel:           capabilities.RiskLevel(tool.ProjectPolicy.RiskLevel),
				SideEffect:          capabilities.SideEffect(tool.ProjectPolicy.SideEffect),
				PolicyRef:           tool.ProjectPolicy.PolicyRef,
				PermissionScope:     capabilities.PermissionScope(tool.ProjectPolicy.PermissionScope),
				ApprovalRequired:    tool.ProjectPolicy.ApprovalRequired,
				IdempotencyRequired: tool.ProjectPolicy.IdempotencyRequired,
			},
		})
	}
	results := make(map[string]capabilities.MCPToolResult, len(config.Results))
	for name, result := range config.Results {
		results[name] = capabilities.MCPToolResult{
			Content:           mcpContentFromConfig(result.Content),
			StructuredContent: mcpStructuredContentFromConfig(result.StructuredContent),
			IsError:           result.IsError,
		}
	}
	return capabilities.NewMockMCPProvider(capabilities.MCPMockConfig{
		ServerID: config.ServerID,
		Tools:    tools,
		Results:  results,
	})
}

// fobrainProviderFromConfig 把文件化 Fobrain 配置转成 provider 边界对象。
// API token 只进入 CredentialResolver，不进入 capability metadata 或 Product Facts。
func fobrainProviderFromConfig(config bootstrap.FobrainConfig) *fobrain.Provider {
	connectorStatus := capabilities.ConnectorStatusUnavailable
	if config.ConnectorStatus.Available {
		connectorStatus = capabilities.ConnectorStatusAvailable
	}
	providerConfig := fobrain.ProviderConfig{
		WorkspaceID:       config.WorkspaceID,
		CredentialBinding: fobrainCredentialBindingFromConfig(config.CredentialBinding),
		ConnectorStatus:   connectorStatus,
		Client:            fobrain.MockClient{},
	}
	if config.Credential.APIToken != "" {
		providerConfig.CredentialResolver = fobrain.StaticCredentialResolver{
			WorkspaceID: config.WorkspaceID,
			APIToken:    config.Credential.APIToken,
		}
	}
	return fobrain.NewProvider(providerConfig)
}

// policyContextsFromConfig 把 provider 配置中的安全凭据摘要传给 policy gate。
func policyContextsFromConfig(config bootstrap.Config) map[string]capabilities.PolicyContext {
	contexts := map[string]capabilities.PolicyContext{}
	if config.Fobrain.Enabled {
		status := capabilities.ConnectorStatusUnavailable
		if config.Fobrain.ConnectorStatus.Available {
			status = capabilities.ConnectorStatusAvailable
		}
		contexts[fobrain.CapabilityCurrentUserContext] = capabilities.PolicyContext{
			WorkspaceID:       config.Fobrain.WorkspaceID,
			CredentialBinding: fobrainCredentialBindingFromConfig(config.Fobrain.CredentialBinding),
			ConnectorStatus:   status,
		}
	}
	return contexts
}

// fobrainCredentialBindingFromConfig 转换安全凭据摘要，不包含真实 token。
func fobrainCredentialBindingFromConfig(config bootstrap.CredentialBinding) capabilities.CredentialBinding {
	return capabilities.CredentialBinding{
		SchemaVersion: config.SchemaVersion,
		WorkspaceID:   config.WorkspaceID,
		System:        config.System,
		Status:        capabilities.CredentialStatus(config.Status),
		DisplayRef:    config.DisplayRef,
		OwnerScope:    capabilities.PermissionScope(config.OwnerScope),
		UpdatedAt:     config.UpdatedAt,
		AuditRef:      config.AuditRef,
	}
}

// mcpJSONSchemaFromConfig 转换 Phase 4.5 支持的 MCP JSON Schema 子集。
func mcpJSONSchemaFromConfig(config bootstrap.MCPJSONSchemaConfig) capabilities.MCPJSONSchema {
	return capabilities.MCPJSONSchema{
		Type:       config.Type,
		Properties: config.Properties,
		Required:   config.Required,
	}
}

// mcpContentFromConfig 转换 mock MCP result content，生产输出仍只使用 StructuredResult。
func mcpContentFromConfig(config []bootstrap.MCPContentConfig) []capabilities.MCPContent {
	out := make([]capabilities.MCPContent, 0, len(config))
	for _, content := range config {
		out = append(out, capabilities.MCPContent{Type: content.Type, Text: content.Text})
	}
	return out
}

// mcpStructuredContentFromConfig 转换 structuredContent，值只作为 candidate 字段输入。
func mcpStructuredContentFromConfig(config map[string]string) map[string]any {
	out := make(map[string]any, len(config))
	for key, value := range config {
		out[key] = value
	}
	return out
}
