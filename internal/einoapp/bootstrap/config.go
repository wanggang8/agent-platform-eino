package bootstrap

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 是后端服务的文件化配置根对象，禁止用环境变量替代这些部署参数。
type Config struct {
	Server         ServerConfig          `yaml:"server"`
	Database       DatabaseConfig        `yaml:"database"`
	LLM            LLMConfig             `yaml:"llm"`
	Fobrain        FobrainConfig         `yaml:"fobrain"`
	Security       SecurityConfig        `yaml:"security"`
	Observability  ObservabilityConfig   `yaml:"observability"`
	Budgets        BudgetConfig          `yaml:"budgets"`
	Capabilities   []CapabilityConfig    `yaml:"capabilities"`
	MCPMockServers []MCPMockServerConfig `yaml:"mcp_mock_servers"`
}

// ServerConfig 定义 HTTP 服务监听地址和超时。
type ServerConfig struct {
	Addr         string        `yaml:"addr"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// DatabaseConfig 定义 Product Facts 等持久化配置。
type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

// LLMConfig 定义模型 provider 配置和凭据绑定的安全展示信息。
type LLMConfig struct {
	Provider          string              `yaml:"provider"`
	BaseURL           string              `yaml:"base_url"`
	APIKey            string              `yaml:"api_key"`
	Model             string              `yaml:"model"`
	Timeout           time.Duration       `yaml:"timeout"`
	ModelLabel        string              `yaml:"-"`
	TimeoutMillis     int                 `yaml:"-"`
	NetworkSafety     NetworkSafetyConfig `yaml:"-"`
	CredentialBinding CredentialBinding   `yaml:"-"`
}

// FobrainConfig 定义 Phase 5 provider PoC 与 Phase 8 live read 的本地文件配置。
// 真实 token 只能在 provider/client 边界内使用，不进入产品事实或日志摘要。
type FobrainConfig struct {
	Enabled           bool                         `yaml:"enabled"`
	ConnectorID       string                       `yaml:"connector_id"`
	WorkspaceID       string                       `yaml:"workspace_id"`
	BaseURL           string                       `yaml:"base_url"`
	Timeout           time.Duration                `yaml:"timeout"`
	Credential        FobrainCredentialConfig      `yaml:"credential"`
	ConnectorStatus   FobrainConnectorStatusConfig `yaml:"connector_status"`
	TLS               FobrainTLSConfig             `yaml:"tls"`
	CredentialBinding CredentialBinding            `yaml:"-"`
}

// FobrainCredentialConfig 保存本机 Fobrain 凭据配置；API token 不得进入 RedactedSummary。
type FobrainCredentialConfig struct {
	Status     string `yaml:"status"`
	DisplayRef string `yaml:"display_ref"`
	OwnerScope string `yaml:"owner_scope"`
	AuthParam  string `yaml:"auth_param"`
	APIToken   string `yaml:"api_token"`
}

// FobrainConnectorStatusConfig 是 Fobrain connector 的安全状态输入，不替代业务读取工具。
type FobrainConnectorStatusConfig struct {
	Mode      string `yaml:"mode"`
	Available bool   `yaml:"available"`
}

// FobrainTLSConfig 定义 Fobrain 专用 TLS 兼容开关，不影响 LLM 或其它 provider。
type FobrainTLSConfig struct {
	InsecureSkipVerify bool `yaml:"insecure_skip_verify"`
}

// NetworkSafetyConfig 定义模型 provider 出站网络策略。
type NetworkSafetyConfig struct {
	RequireHTTPS         bool     `yaml:"require_https"`
	AllowLocalHTTP       bool     `yaml:"allow_local_http"`
	BlockPrivateNetworks bool     `yaml:"block_private_networks"`
	AllowRedirects       bool     `yaml:"allow_redirects"`
	AllowedHosts         []string `yaml:"allowed_hosts"`
}

// CredentialBinding 是可展示的凭据绑定摘要，不包含真实密钥。
type CredentialBinding struct {
	SchemaVersion string `yaml:"schema_version"`
	WorkspaceID   string `yaml:"workspace_id"`
	System        string `yaml:"system"`
	Status        string `yaml:"status"`
	DisplayRef    string `yaml:"display_ref"`
	OwnerScope    string `yaml:"owner_scope"`
	UpdatedAt     string `yaml:"updated_at"`
	AuditRef      string `yaml:"audit_ref"`
}

// SecurityConfig 定义本服务的安全开关。
type SecurityConfig struct {
	RedactSecrets       bool `yaml:"redact_secrets"`
	AllowPrivateNetwork bool `yaml:"allow_private_network"`
}

// ObservabilityConfig 定义日志和诊断开关。
type ObservabilityConfig struct {
	LogLevel         string               `yaml:"log_level"`
	EnableRequestLog bool                 `yaml:"enable_request_log"`
	CostStatistics   CostStatisticsConfig `yaml:"cost_statistics"`
}

// CostStatisticsConfig 定义只用于内部统计的 token 成本单价；默认 0 表示不估算成本。
type CostStatisticsConfig struct {
	InputUnitMicrounits  int64 `yaml:"input_unit_microunits"`
	OutputUnitMicrounits int64 `yaml:"output_unit_microunits"`
}

// BudgetConfig 定义默认执行预算，供后续 run/tool 超时使用。
type BudgetConfig struct {
	DefaultTimeout       time.Duration `yaml:"default_timeout"`
	MaxToolTimeout       time.Duration `yaml:"max_tool_timeout"`
	MaxModelCallsPerRun  int           `yaml:"max_model_calls_per_run"`
	MaxToolCallsPerRun   int           `yaml:"max_tool_calls_per_run"`
	MaxInputTokensPerRun int           `yaml:"max_input_tokens_per_run"`
}

// CapabilityConfig 定义由配置文件注册的能力元数据，不包含 provider raw payload。
type CapabilityConfig struct {
	ID                      string        `yaml:"id"`
	ProviderID              string        `yaml:"provider_id"`
	ToolName                string        `yaml:"tool_name"`
	DisplayName             string        `yaml:"display_name"`
	Description             string        `yaml:"description"`
	ResultSchema            string        `yaml:"result_schema"`
	RiskLevel               string        `yaml:"risk_level"`
	SideEffect              string        `yaml:"side_effect"`
	PolicyRef               string        `yaml:"policy_ref"`
	PermissionScope         string        `yaml:"permission_scope"`
	CredentialBindingPolicy string        `yaml:"credential_binding_policy"`
	ConnectorID             string        `yaml:"connector_id"`
	ApprovalRequired        bool          `yaml:"approval_required"`
	IdempotencyRequired     bool          `yaml:"idempotency_required"`
	Timeout                 time.Duration `yaml:"timeout"`
}

// MCPMockServerConfig 定义 Phase 4.5 可运行 mock MCP server catalog，不用于生产 MCP 连接。
type MCPMockServerConfig struct {
	ServerID string                             `yaml:"server_id"`
	Tools    []MCPMockToolConfig                `yaml:"tools"`
	Results  map[string]MCPMockToolResultConfig `yaml:"results"`
}

// MCPMockToolConfig 定义 mock MCP tools/list 中单个工具的安全元数据。
type MCPMockToolConfig struct {
	Name          string                   `yaml:"name"`
	Title         string                   `yaml:"title"`
	Description   string                   `yaml:"description"`
	InputSchema   MCPJSONSchemaConfig      `yaml:"input_schema"`
	OutputSchema  MCPJSONSchemaConfig      `yaml:"output_schema"`
	Annotations   MCPToolAnnotationsConfig `yaml:"annotations"`
	ProjectPolicy MCPProjectPolicyConfig   `yaml:"project_policy"`
}

// MCPJSONSchemaConfig 是配置文件中的 MCP JSON Schema 子集。
type MCPJSONSchemaConfig struct {
	Type       string            `yaml:"type"`
	Properties map[string]string `yaml:"properties"`
	Required   []string          `yaml:"required"`
}

// MCPToolAnnotationsConfig 是配置文件中的 MCP annotations 子集。
type MCPToolAnnotationsConfig struct {
	ReadOnlyHint    bool `yaml:"read_only_hint"`
	DestructiveHint bool `yaml:"destructive_hint"`
	IdempotentHint  bool `yaml:"idempotent_hint"`
}

// MCPProjectPolicyConfig 是项目可信 policy 覆盖，不能由 MCP server 自行决定。
type MCPProjectPolicyConfig struct {
	RiskLevel           string `yaml:"risk_level"`
	SideEffect          string `yaml:"side_effect"`
	PolicyRef           string `yaml:"policy_ref"`
	PermissionScope     string `yaml:"permission_scope"`
	ApprovalRequired    bool   `yaml:"approval_required"`
	IdempotencyRequired bool   `yaml:"idempotency_required"`
}

// MCPMockToolResultConfig 定义 mock tools/call 的固定安全结果。
type MCPMockToolResultConfig struct {
	Content           []MCPContentConfig `yaml:"content"`
	StructuredContent map[string]string  `yaml:"structured_content"`
	IsError           bool               `yaml:"is_error"`
}

// MCPContentConfig 是 mock MCP result content 的配置形态。
type MCPContentConfig struct {
	Type string `yaml:"type"`
	Text string `yaml:"text"`
}

// RedactedSummary 是可打印的配置摘要，必须保证不泄漏密钥。
type RedactedSummary map[string]any

// LoadConfig 从本地配置文件读取服务配置并执行校验。
func LoadConfig(path string) (Config, error) {
	if strings.TrimSpace(path) == "" {
		return Config{}, errors.New("config path is required")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config file: %w", err)
	}
	cfg.applyDerivedDefaults()
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Validate 校验启动必需项，避免服务以半配置状态运行。
func (cfg Config) Validate() error {
	if strings.TrimSpace(cfg.Server.Addr) == "" {
		return errors.New("server.addr is required")
	}
	if strings.TrimSpace(cfg.Database.Driver) == "" {
		return errors.New("database.driver is required")
	}
	if strings.TrimSpace(cfg.Database.DSN) == "" {
		return errors.New("database.dsn is required")
	}
	if strings.TrimSpace(cfg.LLM.Provider) == "" {
		return errors.New("llm.provider is required")
	}
	if strings.TrimSpace(cfg.LLM.BaseURL) == "" {
		return errors.New("llm.base_url is required")
	}
	if err := validateBaseURL(cfg.LLM.BaseURL); err != nil {
		return err
	}
	if strings.TrimSpace(cfg.LLM.Model) == "" {
		return errors.New("llm.model is required")
	}
	if providerRequiresAPIKey(cfg.LLM.Provider) && strings.TrimSpace(cfg.LLM.APIKey) == "" {
		return errors.New("llm.api_key is required for provider")
	}
	if cfg.LLM.Timeout <= 0 {
		return errors.New("llm.timeout must be positive")
	}
	if err := cfg.Fobrain.validate(); err != nil {
		return err
	}
	if cfg.Server.ReadTimeout <= 0 {
		return errors.New("server.read_timeout must be positive")
	}
	if cfg.Server.WriteTimeout <= 0 {
		return errors.New("server.write_timeout must be positive")
	}
	if cfg.Observability.CostStatistics.InputUnitMicrounits < 0 {
		return errors.New("observability.cost_statistics.input_unit_microunits must not be negative")
	}
	if cfg.Observability.CostStatistics.OutputUnitMicrounits < 0 {
		return errors.New("observability.cost_statistics.output_unit_microunits must not be negative")
	}
	if cfg.Budgets.DefaultTimeout <= 0 {
		return errors.New("budgets.default_timeout must be positive")
	}
	if cfg.Budgets.MaxToolTimeout <= 0 {
		return errors.New("budgets.max_tool_timeout must be positive")
	}
	if cfg.Budgets.MaxModelCallsPerRun <= 0 {
		return errors.New("budgets.max_model_calls_per_run must be positive")
	}
	if cfg.Budgets.MaxToolCallsPerRun <= 0 {
		return errors.New("budgets.max_tool_calls_per_run must be positive")
	}
	if cfg.Budgets.MaxInputTokensPerRun <= 0 {
		return errors.New("budgets.max_input_tokens_per_run must be positive")
	}
	for index, capability := range cfg.Capabilities {
		if strings.TrimSpace(capability.ID) == "" {
			return fmt.Errorf("capabilities[%d].id is required", index)
		}
		if strings.TrimSpace(capability.ProviderID) == "" {
			return fmt.Errorf("capabilities[%d].provider_id is required", index)
		}
		if strings.TrimSpace(capability.ToolName) == "" {
			return fmt.Errorf("capabilities[%d].tool_name is required", index)
		}
		if capability.RiskLevel != "none" &&
			capability.RiskLevel != "low" &&
			capability.RiskLevel != "medium" &&
			capability.RiskLevel != "high" {
			return fmt.Errorf("capabilities[%d].risk_level must be none, low, medium or high", index)
		}
		if strings.TrimSpace(capability.SideEffect) == "" {
			return fmt.Errorf("capabilities[%d].side_effect is required", index)
		}
		if strings.TrimSpace(capability.PolicyRef) == "" {
			return fmt.Errorf("capabilities[%d].policy_ref is required", index)
		}
		if capability.SideEffect != "" &&
			capability.SideEffect != "none" &&
			capability.SideEffect != "read_external" &&
			capability.SideEffect != "write_external" &&
			capability.SideEffect != "local_runtime" {
			return fmt.Errorf("capabilities[%d].side_effect is invalid", index)
		}
		if strings.TrimSpace(capability.PermissionScope) == "" {
			return fmt.Errorf("capabilities[%d].permission_scope is required", index)
		}
		if capability.PermissionScope != "workspace" &&
			capability.PermissionScope != "caller" &&
			capability.PermissionScope != "system" {
			return fmt.Errorf("capabilities[%d].permission_scope is invalid", index)
		}
		if strings.TrimSpace(capability.CredentialBindingPolicy) == "" {
			return fmt.Errorf("capabilities[%d].credential_binding_policy is required", index)
		}
		if capability.CredentialBindingPolicy != "none" &&
			capability.CredentialBindingPolicy != "required" &&
			capability.CredentialBindingPolicy != "optional" {
			return fmt.Errorf("capabilities[%d].credential_binding_policy is invalid", index)
		}
		if capability.SideEffect == "write_external" && (!capability.ApprovalRequired || !capability.IdempotencyRequired) {
			return fmt.Errorf("capabilities[%d].write_external requires approval and idempotency", index)
		}
	}
	for index, server := range cfg.MCPMockServers {
		if strings.TrimSpace(server.ServerID) == "" {
			return fmt.Errorf("mcp_mock_servers[%d].server_id is required", index)
		}
		for toolIndex, tool := range server.Tools {
			if strings.TrimSpace(tool.Name) == "" {
				return fmt.Errorf("mcp_mock_servers[%d].tools[%d].name is required", index, toolIndex)
			}
			if strings.TrimSpace(tool.Description) == "" {
				return fmt.Errorf("mcp_mock_servers[%d].tools[%d].description is required", index, toolIndex)
			}
			if strings.TrimSpace(tool.InputSchema.Type) != "object" {
				return fmt.Errorf("mcp_mock_servers[%d].tools[%d].input_schema.type must be object", index, toolIndex)
			}
			if err := validateMCPProjectPolicy(tool.ProjectPolicy, index, toolIndex); err != nil {
				return err
			}
		}
	}
	return nil
}

// validateMCPProjectPolicy 校验 mock MCP 的项目可信策略覆盖，避免拼写错误静默降低风险等级。
func validateMCPProjectPolicy(policy MCPProjectPolicyConfig, serverIndex int, toolIndex int) error {
	if policy.RiskLevel != "" &&
		policy.RiskLevel != "none" &&
		policy.RiskLevel != "low" &&
		policy.RiskLevel != "medium" &&
		policy.RiskLevel != "high" {
		return fmt.Errorf("mcp_mock_servers[%d].tools[%d].project_policy.risk_level is invalid", serverIndex, toolIndex)
	}
	if policy.SideEffect != "" &&
		policy.SideEffect != "none" &&
		policy.SideEffect != "read_external" &&
		policy.SideEffect != "write_external" &&
		policy.SideEffect != "local_runtime" {
		return fmt.Errorf("mcp_mock_servers[%d].tools[%d].project_policy.side_effect is invalid", serverIndex, toolIndex)
	}
	if policy.PermissionScope != "" &&
		policy.PermissionScope != "workspace" &&
		policy.PermissionScope != "caller" &&
		policy.PermissionScope != "system" {
		return fmt.Errorf("mcp_mock_servers[%d].tools[%d].project_policy.permission_scope is invalid", serverIndex, toolIndex)
	}
	if policy.SideEffect == "write_external" && (!policy.ApprovalRequired || !policy.IdempotencyRequired) {
		return fmt.Errorf("mcp_mock_servers[%d].tools[%d].project_policy.write_external requires approval and idempotency", serverIndex, toolIndex)
	}
	return nil
}

// applyDerivedDefaults 从最小本地配置派生运行时安全字段，避免要求用户维护内部状态。
func (cfg *Config) applyDerivedDefaults() {
	cfg.LLM.applyDerivedDefaults()
	cfg.Fobrain.applyDerivedDefaults()
	cfg.Budgets.applyDerivedDefaults()
}

// applyDerivedDefaults 为新增预算计数字段提供安全默认值，避免旧本地配置升级时启动失败。
func (cfg *BudgetConfig) applyDerivedDefaults() {
	if cfg.MaxModelCallsPerRun == 0 {
		cfg.MaxModelCallsPerRun = 8
	}
	if cfg.MaxToolCallsPerRun == 0 {
		cfg.MaxToolCallsPerRun = 24
	}
	if cfg.MaxInputTokensPerRun == 0 {
		cfg.MaxInputTokensPerRun = 32000
	}
}

// applyDerivedDefaults 派生 model label、timeout_ms、credential binding 和网络安全默认值。
func (cfg *LLMConfig) applyDerivedDefaults() {
	cfg.Provider = strings.TrimSpace(cfg.Provider)
	cfg.BaseURL = strings.TrimSpace(cfg.BaseURL)
	cfg.Model = strings.TrimSpace(cfg.Model)
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	cfg.ModelLabel = cfg.Model
	cfg.TimeoutMillis = int(cfg.Timeout / time.Millisecond)
	cfg.CredentialBinding = derivedCredentialBinding(cfg.Provider, cfg.APIKey)
	cfg.NetworkSafety = derivedNetworkSafety(cfg.BaseURL)
}

// applyDerivedDefaults 归一化 Fobrain 配置，并派生安全 credential binding 摘要。
func (cfg *FobrainConfig) applyDerivedDefaults() {
	cfg.ConnectorID = strings.TrimSpace(cfg.ConnectorID)
	cfg.WorkspaceID = strings.TrimSpace(cfg.WorkspaceID)
	cfg.BaseURL = strings.TrimSpace(cfg.BaseURL)
	cfg.Credential.Status = strings.TrimSpace(cfg.Credential.Status)
	cfg.Credential.DisplayRef = strings.TrimSpace(cfg.Credential.DisplayRef)
	cfg.Credential.OwnerScope = strings.TrimSpace(cfg.Credential.OwnerScope)
	cfg.Credential.AuthParam = strings.TrimSpace(cfg.Credential.AuthParam)
	cfg.Credential.APIToken = strings.TrimSpace(cfg.Credential.APIToken)
	cfg.ConnectorStatus.Mode = strings.TrimSpace(cfg.ConnectorStatus.Mode)
	if !cfg.Enabled {
		return
	}
	if cfg.ConnectorID == "" {
		cfg.ConnectorID = "fobrain"
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 10 * time.Second
	}
	cfg.CredentialBinding = CredentialBinding{
		SchemaVersion: "eino.provider_credential_binding.v1",
		WorkspaceID:   safeConfigIdentifierOrDefault(cfg.WorkspaceID, "workspace"),
		System:        "fobrain",
		Status:        cfg.Credential.Status,
		DisplayRef:    cfg.Credential.DisplayRef,
		OwnerScope:    cfg.Credential.OwnerScope,
		AuditRef:      "audit:fobrain:credential",
	}
}

// validate 校验 Fobrain connector 配置，避免 provider 实现期硬编码 endpoint 或 token 规则。
func (cfg FobrainConfig) validate() error {
	if !cfg.Enabled {
		return nil
	}
	if cfg.ConnectorID != "fobrain" {
		return errors.New("fobrain.connector_id must be fobrain")
	}
	if cfg.WorkspaceID == "" || !validSummaryIdentifier(cfg.WorkspaceID) {
		return errors.New("fobrain.workspace_id is required")
	}
	if cfg.BaseURL == "" {
		return errors.New("fobrain.base_url is required")
	}
	if err := validateFobrainBaseURL(cfg.BaseURL); err != nil {
		return err
	}
	if cfg.Timeout <= 0 {
		return errors.New("fobrain.timeout must be positive")
	}
	if !validFobrainCredentialStatus(cfg.Credential.Status) {
		return errors.New("fobrain.credential.status is invalid")
	}
	if cfg.Credential.DisplayRef == "" || unsafeConfigSummaryText(cfg.Credential.DisplayRef) {
		return errors.New("fobrain.credential.display_ref is invalid")
	}
	if cfg.Credential.OwnerScope != "workspace" {
		return errors.New("fobrain.credential.owner_scope must be workspace")
	}
	if !validAuthParamName(cfg.Credential.AuthParam) {
		return errors.New("fobrain.credential.auth_param is invalid")
	}
	if cfg.ConnectorStatus.Mode != "mock" && cfg.ConnectorStatus.Mode != "live" {
		return errors.New("fobrain.connector_status.mode must be mock or live")
	}
	if cfg.ConnectorStatus.Mode == "live" &&
		cfg.Credential.AuthParam != "authorization" {
		return errors.New("fobrain.credential.auth_param must be authorization in live mode")
	}
	if cfg.ConnectorStatus.Mode == "live" &&
		cfg.Credential.Status != "configured" &&
		cfg.Credential.Status != "bound" {
		return errors.New("fobrain.credential.status must be configured or bound in live mode")
	}
	if cfg.ConnectorStatus.Mode == "live" && cfg.Credential.APIToken == "" {
		return errors.New("fobrain.credential.api_token is required for live mode")
	}
	if cfg.ConnectorStatus.Mode == "live" {
		parsed, _ := url.Parse(cfg.BaseURL)
		if parsed == nil || (parsed.Scheme != "https" && !(parsed.Scheme == "http" && isLocalHost(parsed.Hostname()))) {
			return errors.New("fobrain.base_url must use https or local http in live mode")
		}
	}
	return nil
}

// validAuthParamName 只允许安全 HTTP header 名，值本身仍只能留在 provider 边界。
func validAuthParamName(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" || unsafeConfigSummaryText(value) || strings.ContainsAny(value, " \t\r\n:") {
		return false
	}
	for _, char := range value {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' || char == '-' {
			continue
		}
		return false
	}
	return true
}

// validateFobrainBaseURL 只要求绝对 URL；敏感 userinfo/query 会在摘要中移除。
func validateFobrainBaseURL(value string) error {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return errors.New("fobrain.base_url must be an absolute URL")
	}
	return nil
}

// validFobrainCredentialStatus 固定 provider credential 状态枚举。
func validFobrainCredentialStatus(status string) bool {
	switch status {
	case "missing", "unbound", "configured", "bound":
		return true
	default:
		return false
	}
}

// validateBaseURL 校验本地 YAML 中的 LLM 地址，避免派生出空网络策略。
func validateBaseURL(value string) error {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return errors.New("llm.base_url must be an absolute URL")
	}
	return nil
}

// providerRequiresAPIKey 标记真实模型 provider 的最小凭据要求。
func providerRequiresAPIKey(provider string) bool {
	switch provider {
	case "openai_compatible", "openai", "azure_openai", "custom":
		return true
	default:
		return false
	}
}

// derivedCredentialBinding 只生成可展示安全状态，不暴露 api_key。
func derivedCredentialBinding(provider string, apiKey string) CredentialBinding {
	status := "missing"
	displayRef := "missing"
	if strings.TrimSpace(apiKey) != "" || provider == "mock" {
		status = "bound"
		displayRef = "bound:llm:" + provider
	}
	return CredentialBinding{
		SchemaVersion: "eino.provider_credential_binding.v1",
		WorkspaceID:   "local",
		System:        "llm",
		Status:        status,
		DisplayRef:    displayRef,
		OwnerScope:    "system",
	}
}

// derivedNetworkSafety 从 base_url 推导默认网络策略；高级覆盖后续在真实 provider 阶段再引入。
func derivedNetworkSafety(baseURL string) NetworkSafetyConfig {
	parsed, err := url.Parse(baseURL)
	host := ""
	if err == nil {
		host = parsed.Hostname()
	}
	localHTTP := parsed != nil && parsed.Scheme == "http" && isLocalHost(host)
	return NetworkSafetyConfig{
		RequireHTTPS:         !localHTTP,
		AllowLocalHTTP:       localHTTP,
		BlockPrivateNetworks: true,
		AllowRedirects:       false,
		AllowedHosts:         nonEmptyHosts(host),
	}
}

// isLocalHost 判断本地开发允许的 loopback host。
func isLocalHost(host string) bool {
	switch host {
	case "127.0.0.1", "localhost", "::1":
		return true
	default:
		return false
	}
}

// nonEmptyHosts 返回网络策略可用的 host 列表。
func nonEmptyHosts(host string) []string {
	if host == "" {
		return nil
	}
	return []string{host}
}

// RedactedSummary 返回可写入日志或诊断的脱敏配置摘要。
func (cfg Config) RedactedSummary() RedactedSummary {
	return RedactedSummary{
		"server": map[string]any{
			"addr":          cfg.Server.Addr,
			"read_timeout":  cfg.Server.ReadTimeout.String(),
			"write_timeout": cfg.Server.WriteTimeout.String(),
		},
		"database": map[string]any{
			"driver": cfg.Database.Driver,
			"dsn":    redactValue(cfg.Database.DSN),
		},
		"llm": map[string]any{
			"provider": cfg.LLM.Provider,
			"base_url": redactURL(cfg.LLM.BaseURL),
			"model":    cfg.LLM.Model,
			"timeout":  cfg.LLM.Timeout.String(),
			"credential_binding": map[string]any{
				"schema_version": cfg.LLM.CredentialBinding.SchemaVersion,
				"workspace_id":   cfg.LLM.CredentialBinding.WorkspaceID,
				"system":         cfg.LLM.CredentialBinding.System,
				"status":         cfg.LLM.CredentialBinding.Status,
				"display_ref":    cfg.LLM.CredentialBinding.DisplayRef,
				"owner_scope":    cfg.LLM.CredentialBinding.OwnerScope,
				"updated_at":     cfg.LLM.CredentialBinding.UpdatedAt,
				"audit_ref":      cfg.LLM.CredentialBinding.AuditRef,
			},
			"network_safety": map[string]any{
				"require_https":          cfg.LLM.NetworkSafety.RequireHTTPS,
				"allow_local_http":       cfg.LLM.NetworkSafety.AllowLocalHTTP,
				"block_private_networks": cfg.LLM.NetworkSafety.BlockPrivateNetworks,
				"allow_redirects":        cfg.LLM.NetworkSafety.AllowRedirects,
				"allowed_hosts":          cfg.LLM.NetworkSafety.AllowedHosts,
			},
		},
		"fobrain": redactedFobrainSummary(cfg.Fobrain),
		"security": map[string]any{
			"redact_secrets":        cfg.Security.RedactSecrets,
			"allow_private_network": cfg.Security.AllowPrivateNetwork,
		},
		"observability": map[string]any{
			"log_level":          cfg.Observability.LogLevel,
			"enable_request_log": cfg.Observability.EnableRequestLog,
			"cost_statistics": map[string]any{
				"input_unit_microunits":  cfg.Observability.CostStatistics.InputUnitMicrounits,
				"output_unit_microunits": cfg.Observability.CostStatistics.OutputUnitMicrounits,
			},
		},
		"budgets": map[string]any{
			"default_timeout":         cfg.Budgets.DefaultTimeout.String(),
			"max_tool_timeout":        cfg.Budgets.MaxToolTimeout.String(),
			"max_model_calls_per_run": cfg.Budgets.MaxModelCallsPerRun,
			"max_tool_calls_per_run":  cfg.Budgets.MaxToolCallsPerRun,
			"max_input_units_per_run": cfg.Budgets.MaxInputTokensPerRun,
		},
		"capabilities": redactedCapabilitySummaries(cfg.Capabilities),
	}
}

// redactedFobrainSummary 输出 Fobrain provider 的安全配置摘要，不能包含 token 或 raw payload。
func redactedFobrainSummary(config FobrainConfig) map[string]any {
	if !config.Enabled {
		return map[string]any{"enabled": false}
	}
	return map[string]any{
		"enabled":                  true,
		"connector_id":             safeSummaryIdentifier(config.ConnectorID),
		"workspace_id":             safeSummaryIdentifier(config.WorkspaceID),
		"base_url":                 redactURL(config.BaseURL),
		"timeout":                  config.Timeout.String(),
		"credential_status":        config.CredentialBinding.Status,
		"auth_param":               safeSummaryIdentifier(config.Credential.AuthParam),
		"display_ref":              safeCredentialDisplay(config.CredentialBinding.DisplayRef),
		"credential_scope":         config.CredentialBinding.OwnerScope,
		"connector_mode":           config.ConnectorStatus.Mode,
		"connector_available":      config.ConnectorStatus.Available,
		"tls_insecure_skip_verify": config.TLS.InsecureSkipVerify,
	}
}

// safeCredentialDisplay 只允许安全 display_ref 进入摘要，非法内容折叠。
func safeCredentialDisplay(value string) string {
	if unsafeConfigSummaryText(value) {
		return "redacted"
	}
	if value == "configured" || value == "missing" || value == "unbound" {
		return value
	}
	parts := strings.Split(value, ":")
	if len(parts) == 3 && parts[0] == "bound" && validSummaryIdentifier(parts[1]) && validSummaryIdentifier(parts[2]) {
		return value
	}
	return "redacted"
}

// safeConfigIdentifierOrDefault 给派生安全摘要提供稳定 fallback。
func safeConfigIdentifierOrDefault(value string, fallback string) string {
	if validSummaryIdentifier(value) && !unsafeConfigSummaryText(value) {
		return value
	}
	return fallback
}

// redactedCapabilitySummaries 只输出 capability 元数据摘要，不输出 schema 或 provider 参数。
func redactedCapabilitySummaries(capabilities []CapabilityConfig) []map[string]any {
	if len(capabilities) == 0 {
		return nil
	}
	summaries := make([]map[string]any, 0, len(capabilities))
	for _, capability := range capabilities {
		summaries = append(summaries, map[string]any{
			"id":                        safeSummaryIdentifier(capability.ID),
			"provider_id":               safeSummaryIdentifier(capability.ProviderID),
			"risk_level":                capability.RiskLevel,
			"side_effect":               capability.SideEffect,
			"policy_ref":                safeSummaryPolicyRef(capability.PolicyRef),
			"permission_scope":          capability.PermissionScope,
			"credential_binding_policy": capability.CredentialBindingPolicy,
			"connector_id":              safeSummaryOptionalIdentifier(capability.ConnectorID),
			"approval_required":         capability.ApprovalRequired,
			"idempotency_required":      capability.IdempotencyRequired,
		})
	}
	return summaries
}

// safeSummaryPolicyRef 只保留 policy:<safe-part...> 形式，避免把 URL、DSN 或 query secret 打进摘要。
func safeSummaryPolicyRef(value string) string {
	value = strings.TrimSpace(value)
	parts := strings.Split(value, ":")
	if len(parts) < 2 || parts[0] != "policy" || unsafeConfigSummaryText(value) {
		return "policy:redacted"
	}
	for _, part := range parts[1:] {
		if !validSummaryIdentifier(part) {
			return "policy:redacted"
		}
	}
	return value
}

// safeSummaryOptionalIdentifier 清洗可选标识；空值保留为空，非法值折叠为 redacted。
func safeSummaryOptionalIdentifier(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	return safeSummaryIdentifier(value)
}

// safeSummaryIdentifier 清洗配置摘要中的标识字段，避免误配的 URL/DSN 被打印。
func safeSummaryIdentifier(value string) string {
	value = strings.TrimSpace(value)
	if unsafeConfigSummaryText(value) || !validSummaryIdentifier(value) {
		return "redacted"
	}
	return value
}

// validSummaryIdentifier 限制配置摘要标识为简单 ASCII 标识；允许点用于 capability id。
func validSummaryIdentifier(value string) bool {
	if value == "" {
		return false
	}
	for _, char := range value {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' || char == '-' || char == '.' {
			continue
		}
		return false
	}
	return true
}

// unsafeConfigSummaryText 是配置摘要的最后防线，覆盖密钥、原始 payload 和 raw config 变体。
func unsafeConfigSummaryText(value string) bool {
	lower := strings.ToLower(value)
	for _, marker := range []string{"authorization:", "bearer ", "api_key", "apikey", "token", "password", "credential", "secret", "raw provider", "raw body", "raw error", "raw payload", "raw_payload", "provider_payload", "raw config", "raw_config", "raw-config", "checkpoint-raw", "interrupt-raw"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

// String 将脱敏摘要编码为 YAML，便于人工审查。
func (summary RedactedSummary) String() string {
	encoded, err := yaml.Marshal(map[string]any(summary))
	if err != nil {
		return "<redacted-summary>"
	}
	return string(encoded)
}

// redactValue 对普通敏感字符串做完全遮蔽。
func redactValue(value string) string {
	if value == "" {
		return ""
	}
	return "***"
}

// redactURL 保留 URL 的安全 origin/path，移除 userinfo、query 和 fragment。
func redactURL(value string) string {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return redactValue(value)
	}
	parsed.User = nil
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String()
}
