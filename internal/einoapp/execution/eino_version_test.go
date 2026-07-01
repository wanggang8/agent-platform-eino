package execution_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestEinoVersionPinnedInGoMod(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller unavailable")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", ".."))
	content, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "github.com/cloudwego/eino v0.9.12") {
		t.Fatal("github.com/cloudwego/eino must remain pinned to v0.9.12")
	}
}
