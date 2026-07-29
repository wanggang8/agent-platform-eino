package bootstrap

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// M1Config 只包含 walking skeleton 实际装配项，不保留真实 provider 或写域开关。
type M1Config struct {
	Server     M1ServerConfig     `yaml:"server"`
	Database   M1DatabaseConfig   `yaml:"database"`
	Fixture    M1FixtureConfig    `yaml:"fixture"`
	Freshness  M1FreshnessConfig  `yaml:"freshness"`
	Capability M1CapabilityConfig `yaml:"capability"`
}

type M1ServerConfig struct {
	Addr         string        `yaml:"addr"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type M1DatabaseConfig struct {
	Driver      string        `yaml:"driver"`
	DSN         string        `yaml:"dsn"`
	BusyTimeout time.Duration `yaml:"busy_timeout"`
	PoolSize    int           `yaml:"pool_size"`
}

type M1FixtureConfig struct {
	Path string `yaml:"path"`
	Case string `yaml:"case"`
}

type M1FreshnessConfig struct {
	PolicyVersion string        `yaml:"policy_version"`
	TTL           time.Duration `yaml:"ttl"`
}

type M1CapabilityConfig struct {
	Timeout time.Duration `yaml:"timeout"`
}

// LoadM1Config 严格拒绝未知字段，防止旧 LLM/FOBrain/Action 配置被静默接受。
func LoadM1Config(path string) (M1Config, error) {
	if strings.TrimSpace(path) == "" {
		return M1Config{}, errors.New("config path is required")
	}
	file, err := os.Open(path)
	if err != nil {
		return M1Config{}, fmt.Errorf("open M1 config: %w", err)
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	var config M1Config
	if err := decoder.Decode(&config); err != nil {
		return M1Config{}, fmt.Errorf("parse M1 config: %w", err)
	}
	if err := config.Validate(); err != nil {
		return M1Config{}, err
	}
	return config, nil
}

func (config M1Config) Validate() error {
	if strings.TrimSpace(config.Server.Addr) == "" || config.Server.ReadTimeout <= 0 || config.Server.WriteTimeout <= 0 {
		return errors.New("invalid M1 server config")
	}
	if config.Database.Driver != "sqlite" || strings.TrimSpace(config.Database.DSN) == "" || config.Database.BusyTimeout <= 0 || config.Database.PoolSize < 2 {
		return errors.New("invalid M1 database config")
	}
	if strings.TrimSpace(config.Fixture.Path) == "" || (config.Fixture.Case != "resolved" && config.Fixture.Case != "empty" && config.Fixture.Case != "failed") {
		return errors.New("invalid M1 fixture config")
	}
	if strings.TrimSpace(config.Freshness.PolicyVersion) == "" || config.Freshness.TTL <= 0 {
		return errors.New("invalid M1 freshness config")
	}
	if config.Capability.Timeout <= 0 {
		return errors.New("invalid M1 capability config")
	}
	return nil
}
