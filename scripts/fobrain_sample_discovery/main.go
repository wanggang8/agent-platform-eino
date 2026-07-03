package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"agent-platform-eino/internal/einoapp/bootstrap"
)

const (
	discoveryReportSchemaVersion = "eino.fobrain_sample_discovery_report.v1"
	discoveryScenario            = "fobrain-sample-discovery"
)

// discoveryReport 是样本发现的脱敏机器报告，只记录是否找到样本和验证结果，不记录真实样本值。
type discoveryReport struct {
	SchemaVersion      string           `json:"schema_version"`
	Scenario           string           `json:"scenario"`
	Status             string           `json:"status"`
	ProviderMode       string           `json:"provider_mode"`
	WorkspaceID        string           `json:"workspace_id"`
	TLSPrivateCertMode bool             `json:"tls_private_cert_mode"`
	SamplePresence     samplePresence   `json:"sample_presence"`
	SourceChecks       []discoveryCheck `json:"source_checks"`
	VerificationChecks []discoveryCheck `json:"verification_checks"`
	SampleFile         sampleFileReport `json:"sample_file"`
	RedactionChecks    []string         `json:"redaction_checks"`
	FailureCategory    string           `json:"failure_category"`
	BlocksClaims       []string         `json:"blocks_claims"`
	ReportCreatedAt    string           `json:"report_created_at"`
}

// samplePresence 只表达样本是否存在，避免报告泄漏人员、部门、IP 或详情 ID。
type samplePresence struct {
	OwnerPresent               bool `json:"owner_present"`
	DepartmentPresent          bool `json:"department_present"`
	IPPresent                  bool `json:"ip_present"`
	AssetDetailPresent         bool `json:"asset_detail_present"`
	VulnerabilityDetailPresent bool `json:"vulnerability_detail_present"`
	BusinessPresent            bool `json:"business_present"`
	ThreatNamePresent          bool `json:"threat_name_present"`
}

// discoveryCheck 记录某个只读接口或验证查询是否成功，不保留 query 参数值。
type discoveryCheck struct {
	CheckID         string `json:"check_id"`
	Status          string `json:"status"`
	PathFamily      string `json:"path_family"`
	ItemCount       int    `json:"item_count"`
	FailureCategory string `json:"failure_category"`
}

// sampleFileReport 只记录本地样本文件是否写出，不在报告中记录文件内容。
type sampleFileReport struct {
	Written bool   `json:"written"`
	Path    string `json:"path,omitempty"`
}

type discoveredSamples struct {
	Owner             string `json:"owner"`
	Department        string `json:"department"`
	IP                string `json:"ip"`
	AssetID           string `json:"asset_id,omitempty"`
	AssetNetworkType  string `json:"asset_network_type,omitempty"`
	VulnerabilityID   string `json:"vulnerability_id,omitempty"`
	BusinessName      string `json:"business_name,omitempty"`
	VulnerabilityName string `json:"vulnerability_name,omitempty"`
}

type sampleFile struct {
	SchemaVersion string            `json:"schema_version"`
	Scenario      string            `json:"scenario"`
	Samples       discoveredSamples `json:"samples"`
	CreatedAt     string            `json:"created_at"`
}

type discoveryClient struct {
	baseURL    *url.URL
	httpClient *http.Client
	authParam  string
	apiToken   string
	timeout    time.Duration
}

type payloadPage struct {
	items []map[string]any
}

func main() {
	configPath := flag.String("config", "configs/eino-workbench.local.yaml", "local config file")
	outputPath := flag.String("output", "test-results/eino-workbench-fobrain-sample-discovery-report.json", "safe report output path")
	samplesOutputPath := flag.String("samples-output", "", "optional local-only samples output path; contains business sample values")
	pageSize := flag.Int("page-size", 20, "discovery page size")
	flag.Parse()

	if err := runDiscovery(*configPath, *outputPath, *samplesOutputPath, *pageSize); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// runDiscovery 执行只读样本发现，并始终优先写出安全报告。
func runDiscovery(configPath string, outputPath string, samplesOutputPath string, pageSize int) error {
	cfg, err := bootstrap.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	report, samples, err := executeDiscovery(context.Background(), cfg, pageSize)
	if strings.TrimSpace(samplesOutputPath) != "" && err == nil && samples.complete() {
		if writeErr := writeSamplesFile(samplesOutputPath, samples); writeErr != nil {
			report.Status = "failed"
			report.FailureCategory = "sample_file_write_failed"
			report.BlocksClaims = defaultDiscoveryBlocks()
			err = writeErr
		} else {
			report.SampleFile = sampleFileReport{Written: true, Path: samplesOutputPath}
		}
	}
	if writeErr := writeDiscoveryReport(outputPath, report); writeErr != nil {
		return writeErr
	}
	if err != nil {
		return err
	}
	return nil
}

func executeDiscovery(ctx context.Context, cfg bootstrap.Config, pageSize int) (discoveryReport, discoveredSamples, error) {
	report := newDiscoveryReport(cfg)
	if !cfg.Fobrain.Enabled || cfg.Fobrain.ConnectorStatus.Mode != "live" {
		report.Status = "failed"
		report.FailureCategory = "live_config_unavailable"
		return report, discoveredSamples{}, errors.New("fobrain live config unavailable")
	}
	client, err := newDiscoveryClient(cfg.Fobrain)
	if err != nil {
		report.Status = "failed"
		report.FailureCategory = "client_unavailable"
		return report, discoveredSamples{}, err
	}
	pageSize = boundedDiscoveryPageSize(pageSize)
	samples, err := discoverSamples(ctx, client, pageSize, &report)
	report.SamplePresence = samplePresence{
		OwnerPresent:               strings.TrimSpace(samples.Owner) != "",
		DepartmentPresent:          strings.TrimSpace(samples.Department) != "",
		IPPresent:                  strings.TrimSpace(samples.IP) != "",
		AssetDetailPresent:         strings.TrimSpace(samples.AssetID) != "",
		VulnerabilityDetailPresent: strings.TrimSpace(samples.VulnerabilityID) != "",
		BusinessPresent:            strings.TrimSpace(samples.BusinessName) != "",
		ThreatNamePresent:          strings.TrimSpace(samples.VulnerabilityName) != "",
	}
	if err != nil {
		report.Status = "blocked"
		report.FailureCategory = stableDiscoveryFailure(err)
		report.BlocksClaims = discoveryBlocksForFailure(report.FailureCategory)
		return report, samples, nil
	}
	report.VerificationChecks = verifySamples(ctx, client, samples, pageSize)
	if failed := firstFailedCheck(report.VerificationChecks); failed != "" {
		report.Status = "blocked"
		report.FailureCategory = failed
		return report, samples, nil
	}
	report.Status = "passed"
	report.FailureCategory = "none"
	report.BlocksClaims = []string{}
	return report, samples, nil
}

func discoverSamples(ctx context.Context, client *discoveryClient, pageSize int, report *discoveryReport) (discoveredSamples, error) {
	assetPage, check := client.fetchPage(ctx, "asset_seed", "asset_list", []string{"/api/asset", "/api/v1/asset"}, url.Values{
		"page":     {"1"},
		"per_page": {fmt.Sprint(pageSize)},
	})
	report.SourceChecks = append(report.SourceChecks, check)
	if check.Status != "passed" {
		return discoveredSamples{}, errors.New(check.FailureCategory)
	}
	threatPage, threatCheck := client.fetchPage(ctx, "threat_seed", "threat_center", []string{"/api/threat_center", "/api/v1/threat_center"}, url.Values{
		"page":       {"1"},
		"per_page":   {fmt.Sprint(pageSize)},
		"data_range": {"4"},
	})
	report.SourceChecks = append(report.SourceChecks, threatCheck)

	staffPage, staffCheck := client.fetchPage(ctx, "staff_seed", "staff_list", []string{"/api/v1/user_access/staff_list"}, url.Values{
		"page":     {"1"},
		"per_page": {fmt.Sprint(pageSize)},
	})
	report.SourceChecks = append(report.SourceChecks, staffCheck)

	departmentPage, departmentCheck := client.fetchPage(ctx, "department_seed", "department_list", []string{"/api/v1/personnel_departments"}, url.Values{
		"page":     {"1"},
		"per_page": {fmt.Sprint(pageSize)},
	})
	report.SourceChecks = append(report.SourceChecks, departmentCheck)

	businessPage, businessCheck := client.fetchPage(ctx, "business_seed", "business_list", []string{"/api/business", "/api/v1/business"}, url.Values{
		"page":     {"1"},
		"per_page": {fmt.Sprint(pageSize)},
	})
	report.SourceChecks = append(report.SourceChecks, businessCheck)

	candidates := collectDiscoveryCandidates(assetPage.items, threatPage.items, staffPage.items, departmentPage.items)
	samples := discoveredSamples{}
	samples.AssetID, samples.AssetNetworkType = firstAssetDetailSample(assetPage.items)
	samples.VulnerabilityID, samples.VulnerabilityName = firstVulnerabilityDetailSample(threatPage.items)
	samples.BusinessName = firstBusinessSample(assetPage.items, threatPage.items, businessPage.items)
	if strings.TrimSpace(samples.VulnerabilityName) == "" {
		samples.VulnerabilityName = firstThreatNameSample(threatPage.items)
	}
	var checks []discoveryCheck
	samples.Owner, checks = findVerifiedCandidate(ctx, client, "owner", candidates.owners, pageSize)
	report.SourceChecks = append(report.SourceChecks, checks...)
	samples.Department, checks = findVerifiedCandidate(ctx, client, "department", candidates.departments, pageSize)
	report.SourceChecks = append(report.SourceChecks, checks...)
	samples.IP, checks = findVerifiedCandidate(ctx, client, "ip", candidates.ips, pageSize)
	report.SourceChecks = append(report.SourceChecks, checks...)
	if samples.complete() {
		return samples, nil
	}
	if samples.batchDComplete() {
		return samples, errors.New("batch_e_sample_candidate_not_found")
	}
	return samples, errors.New("sample_candidate_not_found")
}

type discoveryCandidates struct {
	owners      []string
	departments []string
	ips         []string
}

func firstAssetDetailSample(assets []map[string]any) (string, string) {
	for _, asset := range assets {
		id := firstPresentText(asset, "id", "asset_id", "assetId")
		if id == "" {
			continue
		}
		networkType := firstPresentText(asset, "network_type", "ip_type")
		if networkType == "" {
			networkType = "internal"
		}
		return id, networkType
	}
	return "", ""
}

func firstVulnerabilityDetailSample(threats []map[string]any) (string, string) {
	for _, threat := range threats {
		id := firstPresentText(threat, "id", "vulnerability_id", "vuln_id")
		if id == "" {
			continue
		}
		return id, firstThreatNameFromItem(threat)
	}
	return "", ""
}

func firstBusinessSample(itemGroups ...[]map[string]any) string {
	for _, items := range itemGroups {
		for _, item := range items {
			if text := firstNestedText(item["business"], "name"); text != "" {
				return text
			}
			if text := firstNestedText(item["business_info"], "name"); text != "" {
				return text
			}
			if text := firstNestedText(item["business_system"], "name"); text != "" {
				return text
			}
			if text := firstPresentText(item, "business_name"); text != "" {
				return text
			}
		}
	}
	return ""
}

func firstThreatNameSample(threats []map[string]any) string {
	for _, threat := range threats {
		if text := firstThreatNameFromItem(threat); text != "" {
			return text
		}
	}
	return ""
}

func firstThreatNameFromItem(threat map[string]any) string {
	return firstPresentText(threat, "name", "title", "threat_name", "vulnerability_name", "vul_name")
}

func collectDiscoveryCandidates(assets []map[string]any, threats []map[string]any, staffs []map[string]any, departments []map[string]any) discoveryCandidates {
	candidates := discoveryCandidates{}
	for _, asset := range assets {
		candidates.owners = append(candidates.owners, nestedTexts(asset["oper_info"], "name")...)
		candidates.departments = append(candidates.departments, nestedTexts(asset["business_department"], "name")...)
		candidates.departments = append(candidates.departments, nestedTexts(asset["oper_department"], "name")...)
		candidates.ips = append(candidates.ips, firstText(asset["ip"]))
		candidates.ips = append(candidates.ips, textsFromAny(asset["ips"])...)
	}
	for _, threat := range threats {
		candidates.owners = append(candidates.owners, nestedTexts(threat["person_info"], "name")...)
		candidates.departments = append(candidates.departments, nestedTexts(threat["person_department"], "name")...)
		candidates.departments = append(candidates.departments, nestedTexts(threat["business_department"], "name")...)
		candidates.ips = append(candidates.ips, firstText(threat["ip"]))
		candidates.ips = append(candidates.ips, textsFromAny(threat["ips"])...)
	}
	for _, staff := range staffs {
		candidates.owners = append(candidates.owners, firstText(staff["name"]))
		candidates.departments = append(candidates.departments, textsFromAny(staff["department"])...)
	}
	for _, department := range departments {
		candidates.departments = append(candidates.departments, firstText(department["name"]))
	}
	candidates.owners = uniqueNonEmpty(candidates.owners, 20)
	candidates.departments = uniqueNonEmpty(candidates.departments, 20)
	candidates.ips = uniqueNonEmpty(candidates.ips, 20)
	return candidates
}

func findVerifiedCandidate(ctx context.Context, client *discoveryClient, kind string, values []string, pageSize int) (string, []discoveryCheck) {
	checks := []discoveryCheck{}
	for index, value := range values {
		pair := candidateVerificationQueries(kind, value, pageSize)
		if len(pair) != 2 {
			continue
		}
		passed := true
		for _, spec := range pair {
			page, check := client.fetchPage(ctx, fmt.Sprintf("%s_candidate_%d_%s", kind, index+1, spec.suffix), spec.family, spec.paths, spec.query)
			if check.Status == "passed" && len(page.items) == 0 {
				check.Status = "blocked"
				check.FailureCategory = "empty_result"
			}
			checks = append(checks, check)
			if check.Status != "passed" || check.ItemCount == 0 {
				passed = false
			}
		}
		if passed {
			return value, checks
		}
	}
	return "", checks
}

type candidateQuerySpec struct {
	suffix string
	family string
	paths  []string
	query  url.Values
}

func candidateVerificationQueries(kind string, value string, pageSize int) []candidateQuerySpec {
	switch kind {
	case "owner":
		return []candidateQuerySpec{
			{suffix: "asset", family: "asset_list", paths: []string{"/api/asset", "/api/v1/asset"}, query: queryWithSearch("oper_info.name", value, pageSize)},
			{suffix: "threat", family: "threat_center", paths: []string{"/api/threat_center", "/api/v1/threat_center"}, query: queryWithSearch("person_info.name", value, pageSize, "data_range", "4")},
		}
	case "department":
		return []candidateQuerySpec{
			{suffix: "asset", family: "asset_list", paths: []string{"/api/asset", "/api/v1/asset"}, query: queryWithSearch("business_department.name.keyword", value, pageSize)},
			{suffix: "threat", family: "threat_center", paths: []string{"/api/threat_center", "/api/v1/threat_center"}, query: queryWithSearch("person_department.name.keyword", value, pageSize, "data_range", "4")},
		}
	case "ip":
		return []candidateQuerySpec{
			{suffix: "asset", family: "asset_list", paths: []string{"/api/asset", "/api/v1/asset"}, query: queryWithSearch("ip", value, pageSize, "keyword", value)},
			// 漏洞 IP 验证使用 search_condition，避免旧接口 ip 参数把 data_range 覆盖成未处理范围。
			{suffix: "threat", family: "threat_center", paths: []string{"/api/threat_center", "/api/v1/threat_center"}, query: queryWithSearch("ip", value, pageSize, "data_range", "4")},
		}
	default:
		return nil
	}
}

func verifySamples(ctx context.Context, client *discoveryClient, samples discoveredSamples, pageSize int) []discoveryCheck {
	checks := []discoveryCheck{}
	for _, spec := range []struct {
		id     string
		family string
		paths  []string
		query  url.Values
	}{
		{id: "assets_by_owner", family: "asset_list", paths: []string{"/api/asset", "/api/v1/asset"}, query: queryWithSearch("oper_info.name", samples.Owner, pageSize)},
		{id: "vulnerabilities_by_owner", family: "threat_center", paths: []string{"/api/threat_center", "/api/v1/threat_center"}, query: queryWithSearch("person_info.name", samples.Owner, pageSize, "data_range", "4")},
		{id: "assets_by_department", family: "asset_list", paths: []string{"/api/asset", "/api/v1/asset"}, query: queryWithSearch("business_department.name.keyword", samples.Department, pageSize)},
		{id: "vulnerabilities_by_department", family: "threat_center", paths: []string{"/api/threat_center", "/api/v1/threat_center"}, query: queryWithSearch("person_department.name.keyword", samples.Department, pageSize, "data_range", "4")},
		{id: "assets_by_ip", family: "asset_list", paths: []string{"/api/asset", "/api/v1/asset"}, query: queryWithSearch("ip", samples.IP, pageSize, "keyword", samples.IP)},
		{id: "vulnerabilities_by_ip", family: "threat_center", paths: []string{"/api/threat_center", "/api/v1/threat_center"}, query: queryWithSearch("ip", samples.IP, pageSize, "data_range", "4")},
	} {
		page, check := client.fetchPage(ctx, spec.id, spec.family, spec.paths, spec.query)
		if check.Status == "passed" && len(page.items) == 0 {
			check.Status = "blocked"
			check.FailureCategory = "empty_result"
		}
		checks = append(checks, check)
	}
	return checks
}

func queryWithSearch(field string, value string, pageSize int, extra ...string) url.Values {
	values := url.Values{"page": {"1"}, "per_page": {fmt.Sprint(pageSize)}}
	payload := map[string]any{field: []string{strings.TrimSpace(value)}, "operation_type_string": "=="}
	encoded, _ := json.Marshal(payload)
	values.Add("search_condition", string(encoded))
	for i := 0; i+1 < len(extra); i += 2 {
		values.Set(extra[i], extra[i+1])
	}
	return values
}

func newDiscoveryClient(config bootstrap.FobrainConfig) (*discoveryClient, error) {
	parsed, err := url.Parse(strings.TrimSpace(config.BaseURL))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, errors.New("fobrain base url invalid")
	}
	if strings.TrimSpace(config.Credential.APIToken) == "" {
		return nil, errors.New("fobrain credential missing")
	}
	timeout := config.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	transport := http.DefaultTransport
	if config.TLS.InsecureSkipVerify {
		cloned := http.DefaultTransport.(*http.Transport).Clone()
		cloned.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		transport = cloned
	}
	authParam := strings.TrimSpace(config.Credential.AuthParam)
	if authParam == "" {
		return nil, errors.New("fobrain auth param missing")
	}
	return &discoveryClient{
		baseURL: parsed,
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		authParam: authParam,
		apiToken:  config.Credential.APIToken,
		timeout:   timeout,
	}, nil
}

func (client *discoveryClient) fetchPage(ctx context.Context, checkID string, family string, paths []string, query url.Values) (payloadPage, discoveryCheck) {
	check := discoveryCheck{CheckID: checkID, PathFamily: family, Status: "failed", FailureCategory: "not_run"}
	for _, path := range paths {
		items, ok, err := client.fetchPageAtPath(ctx, path, query)
		if err != nil {
			check.FailureCategory = stableHTTPFailure(err)
			return payloadPage{}, check
		}
		if !ok {
			continue
		}
		check.Status = "passed"
		check.ItemCount = len(items)
		check.FailureCategory = "none"
		return payloadPage{items: items}, check
	}
	check.Status = "blocked"
	check.FailureCategory = "path_unavailable"
	return payloadPage{}, check
}

func (client *discoveryClient) fetchPageAtPath(ctx context.Context, path string, query url.Values) ([]map[string]any, bool, error) {
	requestCtx := ctx
	cancel := func() {}
	if _, ok := ctx.Deadline(); !ok && client.timeout > 0 {
		requestCtx, cancel = context.WithTimeout(ctx, client.timeout)
	}
	defer cancel()
	req, err := http.NewRequestWithContext(requestCtx, http.MethodGet, client.endpointWithQuery(path, query), nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set(client.authParam, client.apiToken)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed {
		return nil, false, nil
	}
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, false, errors.New("auth_failure")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, false, errors.New("http_failure")
	}
	var payload any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, false, err
	}
	return unwrapItems(payload), true, nil
}

func (client *discoveryClient) endpointWithQuery(path string, query url.Values) string {
	endpoint := *client.baseURL
	basePath := strings.TrimRight(endpoint.Path, "/")
	requestPath := strings.TrimLeft(strings.TrimSpace(path), "/")
	if strings.HasSuffix(basePath, "/api/v1") && strings.HasPrefix(requestPath, "api/v1/") {
		requestPath = strings.TrimPrefix(requestPath, "api/v1/")
	} else if strings.HasSuffix(basePath, "/api") && strings.HasPrefix(requestPath, "api/") {
		requestPath = strings.TrimPrefix(requestPath, "api/")
	}
	if basePath == "" {
		endpoint.Path = "/" + requestPath
	} else {
		endpoint.Path = strings.TrimRight(basePath+"/"+requestPath, "/")
	}
	endpoint.RawQuery = query.Encode()
	endpoint.Fragment = ""
	return endpoint.String()
}

func unwrapItems(payload any) []map[string]any {
	item, ok := payload.(map[string]any)
	if !ok {
		return []map[string]any{}
	}
	if data, ok := item["data"].(map[string]any); ok {
		if items, ok := data["items"].([]any); ok {
			return mapsFromAny(items)
		}
	}
	if items, ok := item["items"].([]any); ok {
		return mapsFromAny(items)
	}
	return []map[string]any{}
}

func mapsFromAny(items []any) []map[string]any {
	out := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if mapped, ok := item.(map[string]any); ok {
			out = append(out, mapped)
		}
	}
	return out
}

func firstNestedText(value any, key string) string {
	texts := nestedTexts(value, key)
	if len(texts) == 0 {
		return ""
	}
	return texts[0]
}

func nestedTexts(value any, key string) []string {
	out := []string{}
	switch typed := value.(type) {
	case []any:
		for _, item := range typed {
			if mapped, ok := item.(map[string]any); ok {
				if text := firstText(mapped[key]); text != "" {
					out = append(out, text)
				}
			}
		}
	case []map[string]any:
		for _, item := range typed {
			if text := firstText(item[key]); text != "" {
				out = append(out, text)
			}
		}
	case map[string]any:
		if text := firstText(typed[key]); text != "" {
			out = append(out, text)
		}
	}
	return out
}

func firstText(value any) string {
	texts := textsFromAny(value)
	if len(texts) == 0 {
		return ""
	}
	return texts[0]
}

func firstPresentText(item map[string]any, keys ...string) string {
	for _, key := range keys {
		if text := firstText(item[key]); text != "" {
			return text
		}
	}
	return ""
}

func textsFromAny(value any) []string {
	out := []string{}
	switch typed := value.(type) {
	case string:
		if text := strings.TrimSpace(typed); text != "" {
			out = append(out, text)
		}
	case []any:
		for _, item := range typed {
			if text := firstText(item); text != "" {
				out = append(out, text)
			}
		}
	case []string:
		for _, item := range typed {
			if text := strings.TrimSpace(item); text != "" {
				out = append(out, text)
			}
		}
	}
	return out
}

func uniqueNonEmpty(values []string, limit int) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out
}

func (samples discoveredSamples) complete() bool {
	return samples.batchDComplete() &&
		strings.TrimSpace(samples.AssetID) != "" &&
		strings.TrimSpace(samples.VulnerabilityID) != "" &&
		strings.TrimSpace(samples.BusinessName) != "" &&
		strings.TrimSpace(samples.VulnerabilityName) != ""
}

func (samples discoveredSamples) batchDComplete() bool {
	return strings.TrimSpace(samples.Owner) != "" && strings.TrimSpace(samples.Department) != "" && strings.TrimSpace(samples.IP) != ""
}

func firstFailedCheck(checks []discoveryCheck) string {
	for _, check := range checks {
		if check.Status != "passed" || check.ItemCount == 0 {
			if check.FailureCategory != "" && check.FailureCategory != "none" {
				return check.FailureCategory
			}
			return "verification_failed"
		}
	}
	return ""
}

func newDiscoveryReport(cfg bootstrap.Config) discoveryReport {
	return discoveryReport{
		SchemaVersion:      discoveryReportSchemaVersion,
		Scenario:           discoveryScenario,
		Status:             "failed",
		ProviderMode:       cfg.Fobrain.ConnectorStatus.Mode,
		WorkspaceID:        cfg.Fobrain.WorkspaceID,
		TLSPrivateCertMode: cfg.Fobrain.TLS.InsecureSkipVerify,
		SourceChecks:       []discoveryCheck{},
		VerificationChecks: []discoveryCheck{},
		SampleFile:         sampleFileReport{Written: false},
		RedactionChecks: []string{
			"secret_absent",
			"auth_header_absent",
			"raw_payload_absent",
			"sample_values_absent",
			"binding_reference_absent",
		},
		FailureCategory: "not_run",
		BlocksClaims:    defaultDiscoveryBlocks(),
		ReportCreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func defaultDiscoveryBlocks() []string {
	return []string{"fobrain-batch-d live pass", "fobrain-batch-e live pass", "Fobrain 24 readonly final acceptance"}
}

func discoveryBlocksForFailure(failureCategory string) []string {
	// Batch E 独有样本缺失时，不能错误阻断已经具备样本条件的 Batch D live 验收。
	if failureCategory == "batch_e_sample_candidate_not_found" {
		return []string{"fobrain-batch-e live pass", "Fobrain 24 readonly final acceptance"}
	}
	return defaultDiscoveryBlocks()
}

func stableDiscoveryFailure(err error) string {
	if err == nil {
		return "none"
	}
	switch err.Error() {
	case "sample_candidate_not_found":
		return "sample_candidate_not_found"
	case "batch_e_sample_candidate_not_found":
		return "batch_e_sample_candidate_not_found"
	default:
		return "sample_discovery_failed"
	}
}

func stableHTTPFailure(err error) string {
	if err == nil {
		return "none"
	}
	switch err.Error() {
	case "auth_failure":
		return "connector_auth_failure"
	case "http_failure":
		return "connector_execution_failed"
	default:
		return "connector_transport_unavailable"
	}
}

func boundedDiscoveryPageSize(pageSize int) int {
	if pageSize <= 0 {
		return 20
	}
	if pageSize > 100 {
		return 100
	}
	return pageSize
}

func writeDiscoveryReport(outputPath string, report discoveryReport) error {
	return writeJSON0600(outputPath, report)
}

func writeSamplesFile(outputPath string, samples discoveredSamples) error {
	return writeJSON0600(outputPath, sampleFile{
		SchemaVersion: "eino.fobrain_batch_samples.local.v2",
		Scenario:      discoveryScenario,
		Samples:       samples,
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
	})
}

func writeJSON0600(outputPath string, value any) error {
	if strings.TrimSpace(outputPath) == "" {
		return errors.New("output path is required")
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(outputPath, data, 0o600)
}
