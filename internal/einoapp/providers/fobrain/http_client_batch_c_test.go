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

func TestHTTPClientDirectReadBusinessListUsesCurrentAPIPath(t *testing.T) {
	// Batch C business_list 走业务系统列表接口，响应只归一为业务系统安全行摘要。
	var gotPath string
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"data": map[string]any{
				"items": []map[string]any{
					{
						"id":              "biz-1",
						"business_name":   "核心支付",
						"department_name": "交易平台部",
						"person_info":     []map[string]any{{"name": "张三"}},
						"risk_num":        3,
					},
				},
				"page": 1, "per_page": 20, "total": 1,
			},
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
	result, err := client.DirectRead(context.Background(), fobrain.ResolvedCredential{
		WorkspaceID: "ws_fobrain",
		AuthParam:   "authorization",
		APIToken:    "workspace-token",
	}, fobrain.CapabilityBusinessList, fobrain.DirectReadQuery{BusinessName: "核心支付", Page: 1, PageSize: 20})
	if err != nil {
		t.Fatal(err)
	}
	if gotPath != "/api/business" {
		t.Fatalf("path = %q, want /api/business", gotPath)
	}
	for _, want := range []string{"page=1", "per_page=20", "keyword="} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query %q missing %q", gotQuery, want)
		}
	}
	if len(result.Items) != 1 {
		t.Fatalf("items = %+v", result.Items)
	}
	item := result.Items[0]
	if item.EntityRef != "business_system:fobrain:biz-1" ||
		item.DisplayName != "核心支付" ||
		item.OwnerName != "张三" ||
		item.Affected != 3 ||
		!strings.Contains(item.Summary, "交易平台部") {
		t.Fatalf("business item mismatch: %+v", item)
	}
}

func TestHTTPClientDirectReadExternalAssetsFallsBackToLegacyPath(t *testing.T) {
	// 外部高风险资产优先当前私有路径，404 时退回 /api/v1 兼容路径。
	var gotPaths []string
	var gotConditions []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPaths = append(gotPaths, r.URL.Path)
		if r.URL.Path == "/api/external_ip_asset" {
			http.NotFound(w, r)
			return
		}
		if r.URL.Path != "/api/v1/external_ip_asset" {
			t.Fatalf("unexpected path = %s", r.URL.Path)
		}
		for _, encoded := range r.URL.Query()["search_condition"] {
			var condition map[string]any
			if err := json.Unmarshal([]byte(encoded), &condition); err != nil {
				t.Fatalf("search_condition is not JSON: %v", err)
			}
			for key := range condition {
				if key != "operation_type_string" {
					gotConditions = append(gotConditions, key)
				}
			}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"data": map[string]any{
				"items": []map[string]any{
					{"id": "asset-1", "ip": "203.0.113.10", "hostname": "edge-gw", "poc_num": 5},
				},
			},
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
	result, err := client.DirectRead(context.Background(), fobrain.ResolvedCredential{
		WorkspaceID: "ws_fobrain",
		AuthParam:   "authorization",
		APIToken:    "workspace-token",
	}, fobrain.CapabilityExternalHighRiskAssets, fobrain.DirectReadQuery{Severity: "high"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(gotPaths, ",") != "/api/external_ip_asset,/api/v1/external_ip_asset" {
		t.Fatalf("paths = %+v", gotPaths)
	}
	if !containsString(gotConditions, "risk") {
		t.Fatalf("conditions = %+v, want risk", gotConditions)
	}
	if len(result.Items) != 1 || result.Items[0].EntityRef != "asset:fobrain:asset-1" || result.Items[0].DisplayName != "edge-gw" || result.Items[0].Affected != 5 {
		t.Fatalf("external asset result mismatch: %+v", result)
	}
}

func TestHTTPClientDirectReadStatsAndTicketsUseBatchCPaths(t *testing.T) {
	tests := []struct {
		name         string
		capabilityID string
		query        fobrain.DirectReadQuery
		wantPath     string
		assert       func(*testing.T, fobrain.DirectReadResult, string)
	}{
		{
			name:         "vulnerability status summary",
			capabilityID: fobrain.CapabilityVulnerabilityStatusSummary,
			query:        fobrain.DirectReadQuery{Severity: "high"},
			wantPath:     "/api/threat_center/count",
			assert: func(t *testing.T, result fobrain.DirectReadResult, body string) {
				t.Helper()
				if !strings.Contains(body, `"count_name":"vulnerability_status"`) ||
					!strings.Contains(body, `"aggregation_field":"status"`) {
					t.Fatalf("count body mismatch: %s", body)
				}
				if len(result.Metrics) != 2 ||
					result.Metrics[0].Label != "待修复" ||
					result.Metrics[0].Count != 7 ||
					result.Metrics[1].Label != "复测通过" ||
					result.Metrics[1].Count != 3 {
					t.Fatalf("status metrics mismatch: %+v", result.Metrics)
				}
			},
		},
		{
			name:         "ip stats",
			capabilityID: fobrain.CapabilityIPStats,
			query:        fobrain.DirectReadQuery{Keyword: "10.10"},
			wantPath:     "/api/threat_center/relevance/ip_stats",
			assert: func(t *testing.T, result fobrain.DirectReadResult, _ string) {
				t.Helper()
				if len(result.Metrics) != 1 || result.Metrics[0].Label != "10.10.11.12" || result.Metrics[0].Count != 4 {
					t.Fatalf("ip metrics mismatch: %+v", result.Metrics)
				}
			},
		},
		{
			name:         "vul stats",
			capabilityID: fobrain.CapabilityVulStats,
			query:        fobrain.DirectReadQuery{Keyword: "SQL"},
			wantPath:     "/api/threat_center/relevance/vul_stats",
			assert: func(t *testing.T, result fobrain.DirectReadResult, _ string) {
				t.Helper()
				if len(result.Metrics) != 1 || result.Metrics[0].Label != "SQL 注入" || result.Metrics[0].Count != 9 {
					t.Fatalf("vul metrics mismatch: %+v", result.Metrics)
				}
			},
		},
		{
			name:         "pending tickets",
			capabilityID: fobrain.CapabilityPendingTickets,
			query:        fobrain.DirectReadQuery{Status: "open", Page: 1, PageSize: 20},
			wantPath:     "/api/ticket/pending",
			assert: func(t *testing.T, result fobrain.DirectReadResult, _ string) {
				t.Helper()
				if len(result.Items) != 1 || result.Items[0].EntityRef != "ticket:fobrain:ticket-1" || result.Items[0].DisplayName != "高危组件漏洞" {
					t.Fatalf("ticket result mismatch: %+v", result.Items)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			var gotBody string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				if r.Body != nil {
					var raw mapOrSlice
					if err := json.NewDecoder(r.Body).Decode(&raw.value); err == nil && raw.value != nil {
						encoded, _ := json.Marshal(raw.value)
						gotBody = string(encoded)
					}
				}
				switch tt.capabilityID {
				case fobrain.CapabilityVulnerabilityStatusSummary:
					_ = json.NewEncoder(w).Encode(map[string]any{
						"code": 0,
						"data": []map[string]any{
							{
								"count_name": "vulnerability_status",
								"count":      10,
								"aggregation": []map[string]any{
									{"key": "待修复", "origin_key": "10", "count": 7},
									{"key": "复测通过", "origin_key": "30", "count": 3},
								},
							},
						},
					})
				case fobrain.CapabilityIPStats:
					_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"items": []map[string]any{{"ip": "10.10.11.12", "vul_count": 4}}}})
				case fobrain.CapabilityVulStats:
					_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"items": []map[string]any{{"vul_name": "SQL 注入", "ip_count": 9}}}})
				case fobrain.CapabilityPendingTickets:
					_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "data": map[string]any{"list": []map[string]any{{"ticket_id": "ticket-1", "vul_name": "高危组件漏洞", "status": "pending"}}}})
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
			result, err := client.DirectRead(context.Background(), fobrain.ResolvedCredential{
				WorkspaceID: "ws_fobrain",
				AuthParam:   "authorization",
				APIToken:    "workspace-token",
			}, tt.capabilityID, tt.query)
			if err != nil {
				t.Fatal(err)
			}
			if gotPath != tt.wantPath {
				t.Fatalf("path = %q, want %q", gotPath, tt.wantPath)
			}
			tt.assert(t, result, gotBody)
		})
	}
}

func TestHTTPClientBatchCDirectReadBusinessErrorDoesNotLeakRawMaterial(t *testing.T) {
	// Batch C 业务错误必须折叠成安全错误，不能泄漏 token、header 或 raw 响应。
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code":    50001,
			"message": "workspace-token authorization raw payload",
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
	_, err = client.DirectRead(context.Background(), fobrain.ResolvedCredential{
		WorkspaceID: "ws_fobrain",
		AuthParam:   "authorization",
		APIToken:    "workspace-token",
	}, fobrain.CapabilityIPStats, fobrain.DirectReadQuery{})
	if !fobrain.HasReason(err, capabilities.PolicyReasonConnectorExecutionFailed) {
		t.Fatalf("err = %v, want connector execution failed", err)
	}
	for _, forbidden := range []string{"workspace-token", "authorization", "raw payload"} {
		if strings.Contains(strings.ToLower(err.Error()), forbidden) {
			t.Fatalf("safe error leaked %q: %s", forbidden, err.Error())
		}
	}
}

type mapOrSlice struct {
	value any
}
