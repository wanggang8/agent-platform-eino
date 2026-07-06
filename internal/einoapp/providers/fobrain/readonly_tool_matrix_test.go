package fobrain_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"agent-platform-eino/internal/einoapp/providers/fobrain"
)

type readonlyToolMatrixFixture struct {
	SchemaVersion string                    `json:"schema_version"`
	Tools         []readonlyToolMatrixEntry `json:"tools"`
}

type readonlyToolMatrixEntry struct {
	ToolID         string `json:"tool_id"`
	BatchGate      string `json:"batch_gate"`
	Fixture        string `json:"fixture"`
	ReadonlyMatrix bool   `json:"readonly_matrix"`
}

func TestReadonlyToolMatrixDocumentsProviderGaps(t *testing.T) {
	matrix := loadReadonlyToolMatrix(t)
	if matrix.SchemaVersion != "fobrain.tool_matrix.v1" {
		t.Fatalf("schema_version = %q", matrix.SchemaVersion)
	}
	if len(matrix.Tools) != 24 {
		t.Fatalf("readonly tool count = %d, want 24", len(matrix.Tools))
	}

	batchCounts := map[string]int{}
	matrixByID := map[string]readonlyToolMatrixEntry{}
	for _, tool := range matrix.Tools {
		if _, exists := matrixByID[tool.ToolID]; exists {
			t.Fatalf("duplicate readonly tool %s", tool.ToolID)
		}
		if !tool.ReadonlyMatrix {
			t.Fatalf("tool %s is not marked readonly_matrix", tool.ToolID)
		}
		if tool.Fixture == "" {
			t.Fatalf("tool %s missing fixture", tool.ToolID)
		}
		if !fileExists(t, tool.Fixture) {
			t.Fatalf("tool %s fixture missing: %s", tool.ToolID, tool.Fixture)
		}
		matrixByID[tool.ToolID] = tool
		batchCounts[tool.BatchGate]++
	}
	assertBatchCounts(t, batchCounts, map[string]int{"A": 2, "B": 6, "C": 6, "D": 6, "E": 4})

	provider := fobrain.NewProvider(fobrain.ProviderConfig{Client: fobrain.MockClient{}})
	catalog, err := provider.ListCapabilities()
	if err != nil {
		t.Fatal(err)
	}
	catalogIDs := map[string]bool{}
	for _, capability := range catalog {
		catalogIDs[capability.ID] = true
		if capability.ID != fobrain.CapabilityConnectorSecurity {
			if _, ok := matrixByID[capability.ID]; !ok {
				t.Fatalf("advertised readonly capability %s missing from matrix", capability.ID)
			}
		}
	}

	expectedAdvertised := []string{
		fobrain.CapabilityConnectorSecurity,
		fobrain.CapabilityCurrentUserContext,
		fobrain.CapabilityMyPermissions,
		fobrain.CapabilityListAssetsByOwner,
		fobrain.CapabilityListVulnerabilitiesByOwner,
		fobrain.CapabilityListAssetsByDepartment,
		fobrain.CapabilityListVulnerabilitiesByDepartment,
		fobrain.CapabilityListAssetsByIP,
		fobrain.CapabilityListVulnerabilitiesByIP,
		fobrain.CapabilityGetAssetDetail,
		fobrain.CapabilityGetVulnerabilityDetail,
		fobrain.CapabilityBusinessRiskSummary,
		fobrain.CapabilityThreatRelevanceList,
	}
	assertStringSet(t, catalogIDs, expectedAdvertised)

	// Batch B/C 已登记在 24 只读矩阵中，但当前 provider 尚未实现；该断言防止后续误声明 24/24 可用。
	expectedGaps := []string{
		"tool.fobrain.business_list",
		"tool.fobrain.external_high_risk_assets",
		"tool.fobrain.ip_stats",
		"tool.fobrain.my_assets",
		"tool.fobrain.my_business_systems",
		"tool.fobrain.my_department_assets",
		"tool.fobrain.my_department_vulnerabilities",
		"tool.fobrain.my_important_business_systems",
		"tool.fobrain.my_vulnerabilities",
		"tool.fobrain.pending_tickets",
		"tool.fobrain.vul_stats",
		"tool.fobrain.vulnerability_status_summary",
	}
	for _, toolID := range expectedGaps {
		entry, ok := matrixByID[toolID]
		if !ok {
			t.Fatalf("expected gap %s missing from matrix", toolID)
		}
		if entry.BatchGate != "B" && entry.BatchGate != "C" {
			t.Fatalf("gap %s batch = %s, want B or C", toolID, entry.BatchGate)
		}
		if catalogIDs[toolID] {
			t.Fatalf("gap %s must not be advertised before Batch B/C provider implementation", toolID)
		}
	}
}

func loadReadonlyToolMatrix(t *testing.T) readonlyToolMatrixFixture {
	t.Helper()
	data, err := os.ReadFile(repoPath(t, "docs/fixtures/fobrain/tool-matrix-24.json"))
	if err != nil {
		t.Fatal(err)
	}
	var matrix readonlyToolMatrixFixture
	if err := json.Unmarshal(data, &matrix); err != nil {
		t.Fatal(err)
	}
	return matrix
}

func fileExists(t *testing.T, path string) bool {
	t.Helper()
	_, err := os.Stat(repoPath(t, path))
	return err == nil
}

func repoPath(t *testing.T, parts ...string) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller unavailable")
	}
	args := append([]string{filepath.Dir(file), "../../../.."}, parts...)
	return filepath.Clean(filepath.Join(args...))
}

func assertBatchCounts(t *testing.T, got map[string]int, want map[string]int) {
	t.Helper()
	for batch, count := range want {
		if got[batch] != count {
			t.Fatalf("batch %s count = %d, want %d; all counts = %+v", batch, got[batch], count, got)
		}
	}
}

func assertStringSet(t *testing.T, got map[string]bool, want []string) {
	t.Helper()
	wantSet := map[string]bool{}
	for _, value := range want {
		wantSet[value] = true
		if !got[value] {
			t.Fatalf("catalog missing %s", value)
		}
	}
	if len(got) != len(wantSet) {
		t.Fatalf("catalog ids = %v, want %v", sortedKeys(got), sortedKeys(wantSet))
	}
}

func sortedKeys(values map[string]bool) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
