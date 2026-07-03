package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"agent-platform-eino/internal/einoapp/bootstrap"
)

func TestRunDiscoveryWritesSafeReportAndLocalSamples(t *testing.T) {
	// discovery report 不能落真实样本值；本地 samples 文件用于人工执行 Batch D smoke。
	const secret = "sample-discovery-secret"
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authorization") != secret {
			t.Fatalf("authorization header mismatch")
		}
		switch r.URL.Path {
		case "/api/asset":
			_, _ = w.Write([]byte(`{"code":0,"data":{"items":[{"id":"asset-1","ip":"10.10.11.12","network_type":"internal","hostname":"prod-web-01","oper_info":[{"name":"张三"}],"business_department":[{"name":"安全部"}],"business":[{"name":"核心业务"}]}],"page":1,"per_page":20,"total":1}}`))
		case "/api/threat_center":
			// 旧 Fobrain 的 ip query 参数会强制 data_range=1；脚本必须用 search_condition 验证完整漏洞范围。
			if r.URL.Query().Get("ip") != "" {
				_, _ = w.Write([]byte(`{"code":0,"data":{"items":[],"page":1,"per_page":20,"total":0}}`))
				return
			}
			if conditions := strings.Join(r.URL.Query()["search_condition"], "\n"); conditions != "" && !strings.Contains(conditions, `"ip"`) && !strings.Contains(conditions, "person_") {
				_, _ = w.Write([]byte(`{"code":0,"data":{"items":[],"page":1,"per_page":20,"total":0}}`))
				return
			}
			_, _ = w.Write([]byte(`{"code":0,"data":{"items":[{"id":"vuln-1","ip":"10.10.11.12","name":"高危组件漏洞","person_info":[{"name":"张三"}],"person_department":[{"name":"安全部"}],"business":[{"name":"核心业务"}]}],"page":1,"per_page":20,"total":1}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "eino-workbench.yaml")
	reportPath := filepath.Join(tempDir, "report.json")
	samplesPath := filepath.Join(tempDir, "samples.local.json")
	writeDiscoveryTestConfig(t, configPath, server.URL, secret)

	if err := runDiscovery(configPath, reportPath, samplesPath, 20); err != nil {
		t.Fatal(err)
	}
	report := readDiscoveryReport(t, reportPath)
	if report.Status != "passed" || report.FailureCategory != "none" || len(report.BlocksClaims) != 0 {
		t.Fatalf("report status mismatch: %+v", report)
	}
	if !report.SamplePresence.OwnerPresent || !report.SamplePresence.DepartmentPresent || !report.SamplePresence.IPPresent {
		t.Fatalf("sample presence mismatch: %+v", report.SamplePresence)
	}
	if !report.SamplePresence.AssetDetailPresent ||
		!report.SamplePresence.VulnerabilityDetailPresent ||
		!report.SamplePresence.BusinessPresent ||
		!report.SamplePresence.ThreatNamePresent {
		t.Fatalf("Batch E sample presence mismatch: %+v", report.SamplePresence)
	}
	if len(report.VerificationChecks) != 6 {
		t.Fatalf("verification count = %d", len(report.VerificationChecks))
	}
	for _, check := range report.VerificationChecks {
		if check.Status != "passed" || check.ItemCount == 0 {
			t.Fatalf("verification check mismatch: %+v", check)
		}
	}
	if !report.SampleFile.Written {
		t.Fatalf("sample file should be marked written: %+v", report.SampleFile)
	}
	assertDiscoveryReportNoLeak(t, reportPath, secret, "张三", "安全部", "10.10.11.12", "asset-1", "vuln-1", "核心业务", "高危组件漏洞")
	assertSamplesFileContains(t, samplesPath, "张三", "安全部", "10.10.11.12", "asset-1", "vuln-1", "核心业务", "高危组件漏洞")
}

func TestDiscoveryClientRequiresExplicitAuthParam(t *testing.T) {
	_, err := newDiscoveryClient(bootstrapFobrainConfigForDiscoveryTest("https://fobrain.example.test/api", "sample-discovery-secret", ""))
	if err == nil || !strings.Contains(err.Error(), "auth param") {
		t.Fatalf("newDiscoveryClient err = %v, want missing auth param", err)
	}
}

func TestRunDiscoveryBlocksWhenNoStableCandidate(t *testing.T) {
	// 资产列表为空时只能生成 blocking report，不能写出本地样本文件。
	const secret = "sample-discovery-secret"
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/asset" {
			_, _ = w.Write([]byte(`{"code":0,"data":{"items":[],"page":1,"per_page":20,"total":0}}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "eino-workbench.yaml")
	reportPath := filepath.Join(tempDir, "report.json")
	samplesPath := filepath.Join(tempDir, "samples.local.json")
	writeDiscoveryTestConfig(t, configPath, server.URL, secret)

	if err := runDiscovery(configPath, reportPath, samplesPath, 20); err != nil {
		t.Fatal(err)
	}
	report := readDiscoveryReport(t, reportPath)
	if report.Status != "blocked" || report.FailureCategory != "sample_candidate_not_found" {
		t.Fatalf("report status mismatch: %+v", report)
	}
	if len(report.BlocksClaims) == 0 || report.SampleFile.Written {
		t.Fatalf("blocked report mismatch: %+v", report)
	}
	if _, err := os.Stat(samplesPath); !os.IsNotExist(err) {
		t.Fatalf("samples file should not exist, stat err=%v", err)
	}
	assertDiscoveryReportNoLeak(t, reportPath, secret)
}

func TestRunDiscoveryBlocksWhenBatchESamplesMissing(t *testing.T) {
	// Batch E 样本不完整时不能写出 local samples，避免后续 live smoke 使用半成品样本。
	const secret = "sample-discovery-secret"
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authorization") != secret {
			t.Fatalf("authorization header mismatch")
		}
		switch r.URL.Path {
		case "/api/asset":
			_, _ = w.Write([]byte(`{"code":0,"data":{"items":[{"ip":"10.10.11.12","oper_info":[{"name":"张三"}],"business_department":[{"name":"安全部"}]}],"page":1,"per_page":20,"total":1}}`))
		case "/api/threat_center":
			_, _ = w.Write([]byte(`{"code":0,"data":{"items":[{"ip":"10.10.11.12","person_info":[{"name":"张三"}],"person_department":[{"name":"安全部"}]}],"page":1,"per_page":20,"total":1}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "eino-workbench.yaml")
	reportPath := filepath.Join(tempDir, "report.json")
	samplesPath := filepath.Join(tempDir, "samples.local.json")
	writeDiscoveryTestConfig(t, configPath, server.URL, secret)

	if err := runDiscovery(configPath, reportPath, samplesPath, 20); err != nil {
		t.Fatal(err)
	}
	report := readDiscoveryReport(t, reportPath)
	if report.Status != "blocked" || report.FailureCategory != "batch_e_sample_candidate_not_found" {
		t.Fatalf("report status mismatch: %+v", report)
	}
	assertDiscoveryBlocksClaims(t, report.BlocksClaims, "fobrain-batch-e live pass", "Fobrain 24 readonly final acceptance")
	assertDiscoveryBlocksClaimsAbsent(t, report.BlocksClaims, "fobrain-batch-d live pass")
	if report.SampleFile.Written {
		t.Fatalf("sample file should not be written for incomplete Batch E samples: %+v", report.SampleFile)
	}
	if _, err := os.Stat(samplesPath); !os.IsNotExist(err) {
		t.Fatalf("samples file should not exist, stat err=%v", err)
	}
}

func TestRunDiscoveryFindsBusinessSampleFromBusinessList(t *testing.T) {
	// 资产和漏洞页未携带业务字段时，discovery 应按旧项目接口证据从业务系统列表补充业务样本。
	const secret = "sample-discovery-secret"
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("authorization") != secret {
			t.Fatalf("authorization header mismatch")
		}
		switch r.URL.Path {
		case "/api/asset":
			_, _ = w.Write([]byte(`{"code":0,"data":{"items":[{"id":"asset-1","ip":"10.10.11.12","network_type":"internal","oper_info":[{"name":"张三"}],"business_department":[{"name":"安全部"}]}],"page":1,"per_page":20,"total":1}}`))
		case "/api/threat_center":
			_, _ = w.Write([]byte(`{"code":0,"data":{"items":[{"id":"vuln-1","ip":"10.10.11.12","name":"高危组件漏洞","person_info":[{"name":"张三"}],"person_department":[{"name":"安全部"}]}],"page":1,"per_page":20,"total":1}}`))
		case "/api/v1/business":
			_, _ = w.Write([]byte(`{"code":0,"data":{"items":[{"id":"biz-1","business_name":"核心业务"}],"page":1,"per_page":20,"total":1}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "eino-workbench.yaml")
	reportPath := filepath.Join(tempDir, "report.json")
	samplesPath := filepath.Join(tempDir, "samples.local.json")
	writeDiscoveryTestConfig(t, configPath, server.URL, secret)

	if err := runDiscovery(configPath, reportPath, samplesPath, 20); err != nil {
		t.Fatal(err)
	}
	report := readDiscoveryReport(t, reportPath)
	if report.Status != "passed" || report.FailureCategory != "none" || !report.SamplePresence.BusinessPresent {
		t.Fatalf("business list sample should pass discovery: %+v", report)
	}
	assertDiscoveryReportNoLeak(t, reportPath, "核心业务", "biz-1")
	assertSamplesFileContains(t, samplesPath, "核心业务")
}

func assertDiscoveryBlocksClaims(t *testing.T, claims []string, expected ...string) {
	t.Helper()
	for _, want := range expected {
		found := false
		for _, claim := range claims {
			if claim == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("blocks_claims missing %q: %v", want, claims)
		}
	}
}

func assertDiscoveryBlocksClaimsAbsent(t *testing.T, claims []string, unexpected ...string) {
	t.Helper()
	for _, reject := range unexpected {
		for _, claim := range claims {
			if claim == reject {
				t.Fatalf("blocks_claims should not contain %q: %v", reject, claims)
			}
		}
	}
}

func readDiscoveryReport(t *testing.T, outputPath string) discoveryReport {
	t.Helper()
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	var report discoveryReport
	if err := json.Unmarshal(data, &report); err != nil {
		t.Fatal(err)
	}
	return report
}

func assertDiscoveryReportNoLeak(t *testing.T, outputPath string, forbidden ...string) {
	t.Helper()
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	encoded := strings.ToLower(string(data))
	for _, marker := range append(forbidden, "api_token", "credential_ref", "raw provider", "raw body", "safe_summary") {
		if strings.TrimSpace(marker) == "" {
			continue
		}
		if strings.Contains(encoded, strings.ToLower(marker)) {
			t.Fatalf("report leaked forbidden marker %q: %s", marker, encoded)
		}
	}
}

func assertSamplesFileContains(t *testing.T, outputPath string, markers ...string) {
	t.Helper()
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	encoded := string(data)
	for _, marker := range markers {
		if !strings.Contains(encoded, marker) {
			t.Fatalf("samples file missing %q: %s", marker, encoded)
		}
	}
}

func writeDiscoveryTestConfig(t *testing.T, path string, baseURL string, secret string) {
	t.Helper()
	body := `
server:
  addr: "127.0.0.1:0"
  read_timeout: "10s"
  write_timeout: "30s"
database:
  driver: "sqlite"
  dsn: "` + filepath.Join(t.TempDir(), "eino-workbench.db") + `"
llm:
  provider: "mock"
  base_url: "http://127.0.0.1/mock-llm"
  model: "mock-chat"
  timeout: "30s"
security:
  redact_secrets: true
  allow_private_network: false
observability:
  log_level: "debug"
  enable_request_log: true
budgets:
  default_timeout: "60s"
  max_tool_timeout: "120s"
fobrain:
  enabled: true
  connector_id: "fobrain"
  workspace_id: "ws_fobrain"
  base_url: "` + baseURL + `/api"
  timeout: "10s"
  credential:
    status: "bound"
    display_ref: "bound:fobrain:local"
    owner_scope: "workspace"
    auth_param: "authorization"
    api_token: "` + secret + `"
  connector_status:
    mode: "live"
    available: true
  tls:
    insecure_skip_verify: true
`
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func bootstrapFobrainConfigForDiscoveryTest(baseURL string, secret string, authParam string) bootstrap.FobrainConfig {
	return bootstrap.FobrainConfig{
		Enabled:     true,
		ConnectorID: "fobrain",
		WorkspaceID: "ws_fobrain",
		BaseURL:     baseURL,
		Credential: bootstrap.FobrainCredentialConfig{
			Status:     "bound",
			DisplayRef: "bound:fobrain:local",
			OwnerScope: "workspace",
			AuthParam:  authParam,
			APIToken:   secret,
		},
		ConnectorStatus: bootstrap.FobrainConnectorStatusConfig{Mode: "live", Available: true},
	}
}
