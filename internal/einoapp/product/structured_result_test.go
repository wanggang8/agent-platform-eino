package product_test

import (
	"errors"
	"testing"

	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

func TestStructuredResultSafetyGateApprovesSafeCandidate(t *testing.T) {
	// Safety Gate 是 provider/raw 工具输出进入 Product Facts 前的收敛边界。
	gate := product.NewStructuredResultSafetyGate()

	result, err := gate.Approve(product.StructuredResultCandidate{
		SchemaVersion: product.StructuredResultSchemaVersion,
		ResultRef:     "result:call-1",
		SafeSummary:   "发现 1 个资产",
	})
	if err != nil {
		t.Fatal(err)
	}

	if result.SchemaVersion != product.StructuredResultSchemaVersion || result.ResultRef != "result:call-1" || result.SafeSummary != "发现 1 个资产" {
		t.Fatalf("structured result ref mismatch: %+v", result)
	}
}

func TestStructuredResultSafetyGateRejectsUnsafeMaterial(t *testing.T) {
	gate := product.NewStructuredResultSafetyGate()

	for _, testCase := range []struct {
		name      string
		candidate product.StructuredResultCandidate
	}{
		{
			name: "raw provider marker",
			candidate: product.StructuredResultCandidate{
				SchemaVersion: product.StructuredResultSchemaVersion,
				ResultRef:     "result:call-1",
				SafeSummary:   "raw provider body: {...}",
			},
		},
		{
			name: "credential marker",
			candidate: product.StructuredResultCandidate{
				SchemaVersion: product.StructuredResultSchemaVersion,
				ResultRef:     "result:call-2",
				SafeSummary:   "credential_ref=local-secret",
			},
		},
		{
			name: "raw checkpoint ref",
			candidate: product.StructuredResultCandidate{
				SchemaVersion: product.StructuredResultSchemaVersion,
				ResultRef:     "checkpoint-raw-eino-id",
				SafeSummary:   "查询完成",
			},
		},
		{
			name: "raw interrupt ref",
			candidate: product.StructuredResultCandidate{
				SchemaVersion: product.StructuredResultSchemaVersion,
				ResultRef:     "interrupt-raw-eino-id",
				SafeSummary:   "查询完成",
			},
		},
		{
			name: "password marker",
			candidate: product.StructuredResultCandidate{
				SchemaVersion: product.StructuredResultSchemaVersion,
				ResultRef:     "result:call-3",
				SafeSummary:   "password=local-secret",
			},
		},
		{
			name: "token marker",
			candidate: product.StructuredResultCandidate{
				SchemaVersion: product.StructuredResultSchemaVersion,
				ResultRef:     "result:call-4",
				SafeSummary:   "token=local-secret",
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := gate.Approve(testCase.candidate)
			if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
				t.Fatalf("Approve err = %v, want facts.ErrUnsafeFactMaterial", err)
			}
		})
	}
}

func TestStructuredResultSafetyGateRejectsInvalidCandidate(t *testing.T) {
	gate := product.NewStructuredResultSafetyGate()

	for _, testCase := range []struct {
		name      string
		candidate product.StructuredResultCandidate
	}{
		{name: "missing schema", candidate: product.StructuredResultCandidate{ResultRef: "result:call-1", SafeSummary: "查询完成"}},
		{name: "unsupported schema", candidate: product.StructuredResultCandidate{SchemaVersion: "provider.raw.v1", ResultRef: "result:call-1", SafeSummary: "查询完成"}},
		{name: "missing result ref", candidate: product.StructuredResultCandidate{SchemaVersion: product.StructuredResultSchemaVersion, SafeSummary: "查询完成"}},
		{name: "missing safe summary", candidate: product.StructuredResultCandidate{SchemaVersion: product.StructuredResultSchemaVersion, ResultRef: "result:call-1"}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := gate.Approve(testCase.candidate)
			if !errors.Is(err, product.ErrInvalidStructuredResultCandidate) {
				t.Fatalf("Approve err = %v, want product.ErrInvalidStructuredResultCandidate", err)
			}
		})
	}
}
