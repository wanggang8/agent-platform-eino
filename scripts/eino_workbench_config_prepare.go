package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type redactedLLMSummary struct {
	Provider             string                `json:"provider"`
	Model                string                `json:"model"`
	TimeoutMillis        int                   `json:"timeout_ms"`
	CredentialDisplayRef string                `json:"credential_display_ref"`
	NetworkPolicy        redactedNetworkPolicy `json:"network_policy"`
}

type redactedNetworkPolicy struct {
	RequireHTTPS      bool `json:"require_https"`
	AllowLocalHTTP    bool `json:"allow_local_http"`
	AllowedHostsCount int  `json:"allowed_hosts_count"`
}

func main() {
	source := flag.String("source", "", "source yaml config")
	target := flag.String("target", "", "target yaml config")
	addr := flag.String("addr", "", "server address override")
	dsn := flag.String("dsn", "", "sqlite dsn override")
	summary := flag.String("summary", "", "redacted llm summary json")
	checkLLMAPIKey := flag.Bool("check-llm-api-key", false, "check whether llm.api_key is configured")
	flag.Parse()

	if *checkLLMAPIKey {
		if err := checkLLMAPIKeyConfigured(*source); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}
	if err := prepareConfig(*source, *target, *addr, *dsn, *summary); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// prepareConfig 结构化改写 smoke 运行配置，只覆盖服务地址和临时数据库。
func prepareConfig(source string, target string, addr string, dsn string, summaryPath string) error {
	if source == "" || target == "" || addr == "" || dsn == "" || summaryPath == "" {
		return fmt.Errorf("source, target, addr, dsn and summary are required")
	}
	data, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("read source config: %w", err)
	}
	var root map[string]any
	if err := yaml.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("parse source config: %w", err)
	}
	server := ensureMap(root, "server")
	server["addr"] = addr
	database := ensureMap(root, "database")
	database["dsn"] = dsn

	rendered, err := yaml.Marshal(root)
	if err != nil {
		return fmt.Errorf("render target config: %w", err)
	}
	if err := os.WriteFile(target, rendered, 0o600); err != nil {
		return fmt.Errorf("write target config: %w", err)
	}
	summaryBytes, err := json.MarshalIndent(buildLLMSummary(root), "", "  ")
	if err != nil {
		return fmt.Errorf("render llm summary: %w", err)
	}
	if err := os.WriteFile(summaryPath, summaryBytes, 0o600); err != nil {
		return fmt.Errorf("write llm summary: %w", err)
	}
	return nil
}

// checkLLMAPIKeyConfigured 结构化读取 llm.api_key，避免其他 block 的 api_key 误判凭据已配置。
func checkLLMAPIKeyConfigured(source string) error {
	if source == "" {
		return fmt.Errorf("source is required")
	}
	data, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("read source config: %w", err)
	}
	var root map[string]any
	if err := yaml.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("parse source config: %w", err)
	}
	llm, _ := root["llm"].(map[string]any)
	if strings.TrimSpace(stringValue(llm["api_key"], "")) == "" {
		return fmt.Errorf("llm.api_key is missing")
	}
	return nil
}

// ensureMap 返回指定 YAML block；缺失时创建，避免 shell 正则误判跨 block 字段。
func ensureMap(root map[string]any, key string) map[string]any {
	value, ok := root[key].(map[string]any)
	if ok {
		return value
	}
	value = map[string]any{}
	root[key] = value
	return value
}

// buildLLMSummary 从 llm block 生成脱敏报告摘要，不包含 api_key。
func buildLLMSummary(root map[string]any) redactedLLMSummary {
	llm, _ := root["llm"].(map[string]any)
	provider := stringValue(llm["provider"], "openai_compatible")
	model := stringValue(llm["model"], "configured")
	baseURL := stringValue(llm["base_url"], "")
	timeoutMillis := durationMillis(stringValue(llm["timeout"], "30s"))
	parsed, _ := url.Parse(baseURL)
	host := ""
	if parsed != nil {
		host = parsed.Hostname()
	}
	localHTTP := parsed != nil && parsed.Scheme == "http" && isLocalHost(host)
	allowedHostsCount := 0
	if host != "" {
		allowedHostsCount = 1
	}
	return redactedLLMSummary{
		Provider:             provider,
		Model:                model,
		TimeoutMillis:        timeoutMillis,
		CredentialDisplayRef: "bound:llm:" + provider,
		NetworkPolicy: redactedNetworkPolicy{
			RequireHTTPS:      !localHTTP,
			AllowLocalHTTP:    localHTTP,
			AllowedHostsCount: allowedHostsCount,
		},
	}
}

// stringValue 读取 YAML 标量的字符串表示，避免报告生成依赖字段顺序。
func stringValue(value any, fallback string) string {
	switch typed := value.(type) {
	case string:
		if strings.TrimSpace(typed) != "" {
			return strings.TrimSpace(typed)
		}
	case fmt.Stringer:
		rendered := strings.TrimSpace(typed.String())
		if rendered != "" {
			return rendered
		}
	}
	return fallback
}

// durationMillis 将本地配置的 Go duration 转成报告需要的毫秒摘要。
func durationMillis(value string) int {
	duration, err := time.ParseDuration(value)
	if err != nil || duration <= 0 {
		return int((30 * time.Second) / time.Millisecond)
	}
	return int(duration / time.Millisecond)
}

// isLocalHost 判断本地 fake provider smoke 可以使用的 loopback host。
func isLocalHost(host string) bool {
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}
