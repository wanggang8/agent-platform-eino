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
	Server        ServerConfig        `yaml:"server"`
	Database      DatabaseConfig      `yaml:"database"`
	LLM           LLMConfig           `yaml:"llm"`
	Security      SecurityConfig      `yaml:"security"`
	Observability ObservabilityConfig `yaml:"observability"`
	Budgets       BudgetConfig        `yaml:"budgets"`
	Capabilities  []CapabilityConfig  `yaml:"capabilities"`
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
	Model             string              `yaml:"model"`
	ModelLabel        string              `yaml:"model_label"`
	TimeoutMillis     int                 `yaml:"timeout_ms"`
	NetworkSafety     NetworkSafetyConfig `yaml:"network_safety"`
	CredentialBinding CredentialBinding   `yaml:"credential_binding"`
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
	LogLevel         string `yaml:"log_level"`
	EnableRequestLog bool   `yaml:"enable_request_log"`
}

// BudgetConfig 定义默认执行预算，供后续 run/tool 超时使用。
type BudgetConfig struct {
	DefaultTimeout time.Duration `yaml:"default_timeout"`
	MaxToolTimeout time.Duration `yaml:"max_tool_timeout"`
}

// CapabilityConfig 定义由配置文件注册的能力元数据，不包含 provider raw payload。
type CapabilityConfig struct {
	ID               string        `yaml:"id"`
	ProviderID       string        `yaml:"provider_id"`
	ToolName         string        `yaml:"tool_name"`
	DisplayName      string        `yaml:"display_name"`
	Description      string        `yaml:"description"`
	ResultSchema     string        `yaml:"result_schema"`
	RiskLevel        string        `yaml:"risk_level"`
	ApprovalRequired bool          `yaml:"approval_required"`
	Timeout          time.Duration `yaml:"timeout"`
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
	if strings.TrimSpace(cfg.LLM.Model) == "" {
		return errors.New("llm.model is required")
	}
	if cfg.LLM.TimeoutMillis <= 0 {
		return errors.New("llm.timeout_ms must be positive")
	}
	if strings.TrimSpace(cfg.LLM.CredentialBinding.SchemaVersion) == "" {
		return errors.New("llm.credential_binding.schema_version is required")
	}
	if strings.TrimSpace(cfg.LLM.CredentialBinding.WorkspaceID) == "" {
		return errors.New("llm.credential_binding.workspace_id is required")
	}
	if strings.TrimSpace(cfg.LLM.CredentialBinding.System) == "" {
		return errors.New("llm.credential_binding.system is required")
	}
	if strings.TrimSpace(cfg.LLM.CredentialBinding.Status) == "" {
		return errors.New("llm.credential_binding.status is required")
	}
	if strings.TrimSpace(cfg.LLM.CredentialBinding.DisplayRef) == "" {
		return errors.New("llm.credential_binding.display_ref is required")
	}
	if cfg.Server.ReadTimeout <= 0 {
		return errors.New("server.read_timeout must be positive")
	}
	if cfg.Server.WriteTimeout <= 0 {
		return errors.New("server.write_timeout must be positive")
	}
	if cfg.Budgets.DefaultTimeout <= 0 {
		return errors.New("budgets.default_timeout must be positive")
	}
	if cfg.Budgets.MaxToolTimeout <= 0 {
		return errors.New("budgets.max_tool_timeout must be positive")
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
		if capability.RiskLevel != "read_only" && capability.RiskLevel != "write" {
			return fmt.Errorf("capabilities[%d].risk_level must be read_only or write", index)
		}
	}
	return nil
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
			"provider":   cfg.LLM.Provider,
			"base_url":   redactURL(cfg.LLM.BaseURL),
			"model":      cfg.LLM.Model,
			"timeout_ms": cfg.LLM.TimeoutMillis,
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
		"security": map[string]any{
			"redact_secrets":        cfg.Security.RedactSecrets,
			"allow_private_network": cfg.Security.AllowPrivateNetwork,
		},
		"observability": map[string]any{
			"log_level":          cfg.Observability.LogLevel,
			"enable_request_log": cfg.Observability.EnableRequestLog,
		},
		"budgets": map[string]any{
			"default_timeout":  cfg.Budgets.DefaultTimeout.String(),
			"max_tool_timeout": cfg.Budgets.MaxToolTimeout.String(),
		},
		"capabilities": redactedCapabilitySummaries(cfg.Capabilities),
	}
}

// redactedCapabilitySummaries 只输出 capability 元数据摘要，不输出 schema 或 provider 参数。
func redactedCapabilitySummaries(capabilities []CapabilityConfig) []map[string]any {
	if len(capabilities) == 0 {
		return nil
	}
	summaries := make([]map[string]any, 0, len(capabilities))
	for _, capability := range capabilities {
		summaries = append(summaries, map[string]any{
			"id":                capability.ID,
			"provider_id":       capability.ProviderID,
			"risk_level":        capability.RiskLevel,
			"approval_required": capability.ApprovalRequired,
		})
	}
	return summaries
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
