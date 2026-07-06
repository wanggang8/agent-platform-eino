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
		fobrain.CapabilityMyAssets,
		fobrain.CapabilityMyDepartmentAssets,
		fobrain.CapabilityMyVulnerabilities,
		fobrain.CapabilityMyDepartmentVulnerabilities,
		fobrain.CapabilityMyBusinessSystems,
		fobrain.CapabilityMyImportantBusinessSystems,
		fobrain.CapabilityBusinessList,
		fobrain.CapabilityExternalHighRiskAssets,
		fobrain.CapabilityVulnerabilityStatusSummary,
		fobrain.CapabilityPendingTickets,
		fobrain.CapabilityIPStats,
		fobrain.CapabilityVulStats,
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

	// Mock provider 只证明 catalog/mapper 覆盖；live、视觉和 Action API 同源验收仍由后续门禁声明。
	for toolID := range matrixByID {
		if !catalogIDs[toolID] {
			t.Fatalf("readonly tool %s missing from mock provider catalog", toolID)
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
