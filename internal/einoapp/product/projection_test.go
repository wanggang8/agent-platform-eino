package product_test

import (
	"context"
	"testing"

	"agent-platform-eino/internal/einoapp/product"
)

func TestEmptyProjectionReturnsWorkbenchViewSchemaDirectly(t *testing.T) {
	projection := product.NewEmptyProjection()

	view, err := projection.WorkbenchView(context.Background(), "ws-demo")
	if err != nil {
		t.Fatal(err)
	}

	if view.SchemaVersion != "eino_workbench_view.v1" {
		t.Fatalf("schema_version = %q", view.SchemaVersion)
	}
	if view.WorkspaceID != "ws-demo" {
		t.Fatalf("workspace_id = %q", view.WorkspaceID)
	}
	if view.Timeline == nil {
		t.Fatal("timeline must be an empty array, not nil")
	}
}

func TestProductErrorConvertsToSafeAPIError(t *testing.T) {
	err := product.NewSafeError("projection_unavailable", "投影暂不可用", false)

	if err.Code != "projection_unavailable" {
		t.Fatalf("code = %q", err.Code)
	}
	if err.SafeDetail == "" {
		t.Fatal("safe detail is empty")
	}
}
