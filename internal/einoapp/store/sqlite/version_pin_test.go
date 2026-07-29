package sqlite_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func TestModerncSQLiteVersionPinnedToWALSafetyFloor(t *testing.T) {
	// WAL writer 只能使用 ADR 已验证的 SQLite 3.51.3 安全下限。
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime caller unavailable")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "..", ".."))
	content, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "modernc.org/sqlite v1.46.2") {
		t.Fatal("modernc.org/sqlite must be pinned to v1.46.2")
	}
}

func TestModerncSQLiteRuntimeMeetsWALSafetyFloor(t *testing.T) {
	// 依赖版本之外再验证实际链接的 SQLite，避免模块标签与运行库发生偏差。
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var version string
	if err := db.QueryRow("SELECT sqlite_version()").Scan(&version); err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		t.Fatalf("malformed sqlite_version %q", version)
	}
	numbers := make([]int, 3)
	for i, part := range parts {
		numbers[i], err = strconv.Atoi(part)
		if err != nil {
			t.Fatalf("malformed sqlite_version %q", version)
		}
	}
	if numbers[0] < 3 || (numbers[0] == 3 && numbers[1] < 51) || (numbers[0] == 3 && numbers[1] == 51 && numbers[2] < 3) {
		t.Fatalf("sqlite_version %q is below WAL safety floor 3.51.3", version)
	}
}
