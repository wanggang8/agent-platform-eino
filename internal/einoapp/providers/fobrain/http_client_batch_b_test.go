package fobrain_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"agent-platform-eino/internal/einoapp/capabilities"
	"agent-platform-eino/internal/einoapp/providers/fobrain"
)

func TestHTTPClientMyScopeQueryBuildsCurrentUserAndDepartmentFilters(t *testing.T) {
	// Batch B 的“我的/本部门”范围必须从当前用户接口派生，不能要求用户显式传 owner 或 department。
	tests := []struct {
		name          string
		capabilityID  string
		wantPath      string
		wantCondition string
		wantValue     any
		wantOperation string
		wantDataRange string
		wantItemRef   string
	}{
		{
			name:          "my assets",
			capabilityID:  fobrain.CapabilityMyAssets,
			wantPath:      "/api/asset",
			wantCondition: "oper_info.name",
			wantValue:     "王五",
			wantOperation: "==",
			wantItemRef:   "asset:fobrain:asset-1",
		},
		{
			name:          "department assets",
			capabilityID:  fobrain.CapabilityMyDepartmentAssets,
			wantPath:      "/api/asset",
			wantCondition: "business_department.name.keyword",
			wantValue:     "安全部",
			wantOperation: "==",
			wantItemRef:   "asset:fobrain:asset-1",
		},
		{
			name:          "my vulnerabilities",
			capabilityID:  fobrain.CapabilityMyVulnerabilities,
			wantPath:      "/api/threat_center",
			wantCondition: "person_info.name",
			wantValue:     "王五",
			wantOperation: "==",
			wantDataRange: "4",
			wantItemRef:   "vuln:fobrain:vul-1",
		},
		{
			name:          "department vulnerabilities",
			capabilityID:  fobrain.CapabilityMyDepartmentVulnerabilities,
			wantPath:      "/api/threat_center",
			wantCondition: "person_department.name.keyword",
			wantValue:     "安全部",
			wantOperation: "==",
			wantDataRange: "4",
			wantItemRef:   "vuln:fobrain:vul-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			var gotDataRange string
			gotConditions := map[string]batchBHTTPCondition{}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.URL.Path {
				case "/api/v1/user":
					_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{
						"display_name":    "王五",
						"department_name": "安全部",
					}})
				case "/api/asset", "/api/threat_center":
					gotPath = r.URL.Path
					gotDataRange = r.URL.Query().Get("data_range")
					gotConditions = parseBatchBHTTPConditions(t, r.URL.Query()["search_condition"])
					_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"items": batchBHTTPItems(tt.capabilityID)}})
				default:
					t.Fatalf("unexpected path = %s", r.URL.Path)
				}
			}))
			defer server.Close()

			client, err := fobrain.NewHTTPClient(fobrain.HTTPClientConfig{
				BaseURL:    server.URL + "/api",
				Timeout:    time.Second,
				HTTPClient: server.Client(),
			})
			if err != nil {
				t.Fatal(err)
			}
			result, err := client.MyScopeQuery(context.Background(), fobrain.ResolvedCredential{
				WorkspaceID: "ws_fobrain",
				AuthParam:   "authorization",
				APIToken:    "workspace-token",
			}, tt.capabilityID, fobrain.MyScopeQuery{Keyword: "生产", Page: 1, PageSize: 20})
			if err != nil {
				t.Fatal(err)
			}
			if gotPath != tt.wantPath {
				t.Fatalf("path = %q, want %q", gotPath, tt.wantPath)
			}
			if !batchBHTTPConditionMatches(gotConditions[tt.wantCondition], tt.wantValue, tt.wantOperation) {
				t.Fatalf("conditions = %+v, want %s", gotConditions, tt.wantCondition)
			}
			if gotDataRange != tt.wantDataRange {
				t.Fatalf("data_range = %q, want %q", gotDataRange, tt.wantDataRange)
			}
			if len(result.Items) != 1 || result.Items[0].EntityRef != tt.wantItemRef {
				t.Fatalf("result mismatch: %+v", result)
			}
		})
	}
}

func TestHTTPClientMyScopeQueryBusinessSystemsUsesOwnerScopeAndImportantFilter(t *testing.T) {
	// 业务系统 owner scope 必须来自当前用户条件，用户 keyword 只能作为业务名称收窄过滤。
	var gotPath string
	var gotQuery string
	gotConditions := map[string]batchBHTTPCondition{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/user":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{
				"display_name":    "王五",
				"department_name": "安全部",
			}})
		case "/api/business":
			gotPath = r.URL.Path
			gotQuery = r.URL.RawQuery
			gotConditions = parseBatchBHTTPConditions(t, r.URL.Query()["search_condition"])
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"items": []map[string]any{{
				"id":              "biz-1",
				"business_name":   "核心支付",
				"department_name": "安全部",
				"person_info":     []map[string]any{{"name": "王五"}},
				"risk_num":        2,
			}}}})
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client, err := fobrain.NewHTTPClient(fobrain.HTTPClientConfig{
		BaseURL:    server.URL + "/api",
		Timeout:    time.Second,
		HTTPClient: server.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.MyScopeQuery(context.Background(), fobrain.ResolvedCredential{
		WorkspaceID: "ws_fobrain",
		AuthParam:   "authorization",
		APIToken:    "workspace-token",
	}, fobrain.CapabilityMyImportantBusinessSystems, fobrain.MyScopeQuery{Keyword: "不应放宽", Page: 1, PageSize: 20})
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/api/business" {
		t.Fatalf("path = %q, want /api/business", gotPath)
	}
	for _, want := range []string{"page=1", "per_page=20"} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query %q missing %q", gotQuery, want)
		}
	}
	if strings.Contains(gotQuery, "keyword=") {
		t.Fatalf("business owner scope must use search_condition, got query %q", gotQuery)
	}
	if !batchBHTTPConditionMatches(gotConditions["person_base.name"], "王五", "==") {
		t.Fatalf("conditions = %+v, want owner person_base.name", gotConditions)
	}
	if !batchBHTTPConditionMatches(gotConditions["assets_attribute.important_types"], []any{float64(1), float64(2)}, "in") {
		t.Fatalf("conditions = %+v, want important_types filter", gotConditions)
	}
	if !batchBHTTPConditionMatches(gotConditions["business_name"], "不应放宽", "==") {
		t.Fatalf("conditions = %+v, want business_name filter for user keyword", gotConditions)
	}
	if len(result.Items) != 1 || result.Items[0].EntityRef != "business_system:fobrain:biz-1" || result.Items[0].OwnerName != "王五" {
		t.Fatalf("business result mismatch: %+v", result)
	}
}

func TestHTTPClientMyScopeQueryFallsBackToLegacyPath(t *testing.T) {
	// Batch B 与 Batch C/D 保持一致：当前 /api 路径不存在时才回退 /api/v1。
	var gotPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPaths = append(gotPaths, r.URL.Path)
		switch r.URL.Path {
		case "/api/v1/user":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"display_name": "王五"}})
		case "/api/asset":
			http.NotFound(w, r)
		case "/api/v1/asset":
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"items": []map[string]any{{
				"id":       "asset-1",
				"hostname": "prod-web-01",
			}}}})
		default:
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client, err := fobrain.NewHTTPClient(fobrain.HTTPClientConfig{
		BaseURL:    server.URL + "/api",
		Timeout:    time.Second,
		HTTPClient: server.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.MyScopeQuery(context.Background(), fobrain.ResolvedCredential{
		WorkspaceID: "ws_fobrain",
		AuthParam:   "authorization",
		APIToken:    "workspace-token",
	}, fobrain.CapabilityMyAssets, fobrain.MyScopeQuery{})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(gotPaths, ",") != "/api/v1/user,/api/asset,/api/v1/asset" {
		t.Fatalf("paths = %+v", gotPaths)
	}
	if len(result.Items) != 1 || result.Items[0].EntityRef != "asset:fobrain:asset-1" {
		t.Fatalf("fallback result mismatch: %+v", result)
	}
}

func TestHTTPClientMyScopeQueryRejectsUnsafeBusinessError(t *testing.T) {
	// provider 业务错误只输出安全原因，不能把 raw body、token 或 URL 泄漏到上层。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/user" {
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"display_name": "王五"}})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 500,
			"msg":  "authorization workspace-token raw /api/asset failed",
		})
	}))
	defer server.Close()

	client, err := fobrain.NewHTTPClient(fobrain.HTTPClientConfig{
		BaseURL:    server.URL + "/api",
		Timeout:    time.Second,
		HTTPClient: server.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.MyScopeQuery(context.Background(), fobrain.ResolvedCredential{
		WorkspaceID: "ws_fobrain",
		AuthParam:   "authorization",
		APIToken:    "workspace-token",
	}, fobrain.CapabilityMyAssets, fobrain.MyScopeQuery{})
	if !fobrain.HasReason(err, capabilities.PolicyReasonConnectorExecutionFailed) {
		t.Fatalf("err = %v, want safe execution failure", err)
	}
	lower := strings.ToLower(err.Error())
	for _, forbidden := range []string{"workspace-token", "authorization", "/api/asset", "raw"} {
		if strings.Contains(lower, forbidden) {
			t.Fatalf("error leaked %q: %v", forbidden, err)
		}
	}
}

func TestHTTPClientMyScopeQueryUsesStableReasonWhenCurrentUserScopeMissing(t *testing.T) {
	// 当前用户上下文缺字段时必须返回稳定原因码，smoke/report 不能依赖中文错误文案做分类。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/user" {
			_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"display_name": "王五"}})
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	client, err := fobrain.NewHTTPClient(fobrain.HTTPClientConfig{
		BaseURL:    server.URL + "/api",
		Timeout:    time.Second,
		HTTPClient: server.Client(),
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.MyScopeQuery(context.Background(), fobrain.ResolvedCredential{
		WorkspaceID: "ws_fobrain",
		AuthParam:   "authorization",
		APIToken:    "workspace-token",
	}, fobrain.CapabilityMyDepartmentAssets, fobrain.MyScopeQuery{})
	if !fobrain.HasReason(err, fobrain.PolicyReasonMissingCurrentUserScope) {
		t.Fatalf("err = %v, want missing current user scope reason", err)
	}
}

func batchBHTTPItems(capabilityID string) []map[string]any {
	switch capabilityID {
	case fobrain.CapabilityMyAssets, fobrain.CapabilityMyDepartmentAssets:
		return []map[string]any{{"id": "asset-1", "hostname": "prod-web-01", "ip": "10.10.11.12"}}
	case fobrain.CapabilityMyVulnerabilities, fobrain.CapabilityMyDepartmentVulnerabilities:
		return []map[string]any{{"id": "vul-1", "name": "SQL 注入", "status": 10, "level": 3}}
	default:
		return nil
	}
}

type batchBHTTPCondition struct {
	Value     any
	Operation string
}

func parseBatchBHTTPConditions(t *testing.T, encodedConditions []string) map[string]batchBHTTPCondition {
	t.Helper()
	out := map[string]batchBHTTPCondition{}
	for _, encoded := range encodedConditions {
		var condition map[string]any
		if err := json.Unmarshal([]byte(encoded), &condition); err != nil {
			t.Fatalf("search_condition is not JSON: %v", err)
		}
		operation, _ := condition["operation_type_string"].(string)
		for key, value := range condition {
			if key != "operation_type_string" {
				out[key] = batchBHTTPCondition{Value: value, Operation: operation}
			}
		}
	}
	return out
}

func batchBHTTPConditionMatches(condition batchBHTTPCondition, wantValue any, wantOperation string) bool {
	if condition.Operation != wantOperation {
		return false
	}
	values, ok := condition.Value.([]any)
	if !ok || len(values) == 0 {
		return false
	}
	switch want := wantValue.(type) {
	case string:
		value, _ := values[0].(string)
		return value == want
	case []any:
		if len(values) != len(want) {
			return false
		}
		for index := range want {
			if values[index] != want[index] {
				return false
			}
		}
		return true
	default:
		return values[0] == want
	}
}
