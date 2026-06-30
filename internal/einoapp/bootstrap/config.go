package bootstrap

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server        ServerConfig        `yaml:"server"`
	Database      DatabaseConfig      `yaml:"database"`
	LLM           LLMConfig           `yaml:"llm"`
	Security      SecurityConfig      `yaml:"security"`
	Observability ObservabilityConfig `yaml:"observability"`
	Budgets       BudgetConfig        `yaml:"budgets"`
}

type ServerConfig struct {
	Addr         string        `yaml:"addr"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type DatabaseConfig struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type LLMConfig struct {
	Provider          string              `yaml:"provider"`
	BaseURL           string              `yaml:"base_url"`
	Model             string              `yaml:"model"`
	ModelLabel        string              `yaml:"model_label"`
	TimeoutMillis     int                 `yaml:"timeout_ms"`
	NetworkSafety     NetworkSafetyConfig `yaml:"network_safety"`
	CredentialBinding CredentialBinding   `yaml:"credential_binding"`
}

type NetworkSafetyConfig struct {
	RequireHTTPS         bool     `yaml:"require_https"`
	AllowLocalHTTP       bool     `yaml:"allow_local_http"`
	BlockPrivateNetworks bool     `yaml:"block_private_networks"`
	AllowRedirects       bool     `yaml:"allow_redirects"`
	AllowedHosts         []string `yaml:"allowed_hosts"`
}

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

type SecurityConfig struct {
	RedactSecrets       bool `yaml:"redact_secrets"`
	AllowPrivateNetwork bool `yaml:"allow_private_network"`
}

type ObservabilityConfig struct {
	LogLevel         string `yaml:"log_level"`
	EnableRequestLog bool   `yaml:"enable_request_log"`
}

type BudgetConfig struct {
	DefaultTimeout time.Duration `yaml:"default_timeout"`
	MaxToolTimeout time.Duration `yaml:"max_tool_timeout"`
}

type RedactedSummary map[string]any

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
	return nil
}

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
			"base_url":   cfg.LLM.BaseURL,
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
	}
}

func (summary RedactedSummary) String() string {
	encoded, err := yaml.Marshal(map[string]any(summary))
	if err != nil {
		return "<redacted-summary>"
	}
	return string(encoded)
}

func redactValue(value string) string {
	if value == "" {
		return ""
	}
	return "***"
}
