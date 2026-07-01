package product_test

import (
	"errors"
	"testing"

	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

func TestAssistantSafetyGateApprovesSafeText(t *testing.T) {
	// Assistant 文本进入 Turn fact 前也必须经过安全门，避免绕过 StructuredResult 事实边界。
	gate := product.NewAssistantSafetyGate()

	content, err := gate.Approve("已完成查询，发现 1 个资产。")
	if err != nil {
		t.Fatal(err)
	}

	if content != "已完成查询，发现 1 个资产。" {
		t.Fatalf("content = %q", content)
	}
}

func TestAssistantSafetyGateRejectsUnsafeText(t *testing.T) {
	gate := product.NewAssistantSafetyGate()

	for _, content := range []string{
		"Authorization: Bearer local-secret",
		"raw provider body: {...}",
		"checkpoint-raw-eino-id",
		"interrupt-raw-eino-id",
		"api_key=secret",
		"password=local-secret",
		"token=local-secret",
	} {
		t.Run(content, func(t *testing.T) {
			_, err := gate.Approve(content)
			if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
				t.Fatalf("Approve err = %v, want facts.ErrUnsafeFactMaterial", err)
			}
		})
	}
}

func TestAssistantSafetyGateRejectsEmptyText(t *testing.T) {
	gate := product.NewAssistantSafetyGate()

	_, err := gate.Approve("  ")
	if !errors.Is(err, product.ErrInvalidAssistantMessage) {
		t.Fatalf("Approve err = %v, want product.ErrInvalidAssistantMessage", err)
	}
}
