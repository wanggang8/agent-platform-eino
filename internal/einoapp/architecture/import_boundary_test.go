package architecture_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const modulePrefix = "agent-platform-eino/internal/einoapp/"

// TestImportBoundaryDoesNotUseLegacyProjectPackages 防止新项目重新依赖旧项目内部包。
func TestImportBoundaryDoesNotUseLegacyProjectPackages(t *testing.T) {
	forbidden := []string{
		"ai-agent",
		"internal/interfaces/presentation",
		"internal/registry/capability",
		"internal/registry/providerapi",
		"internal/business/fobrain",
		"internal/infrastructure/eino",
		"internal/infrastructure/llm",
	}

	root := filepath.Join("..", "..", "..")
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() && shouldSkipProjectDir(entry.Name()) {
			return filepath.SkipDir
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, spec := range file.Imports {
			checkImport(t, path, spec, forbidden)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// shouldSkipProjectDir 跳过不会参与 Go import 边界检查的目录。
func shouldSkipProjectDir(name string) bool {
	switch name {
	case ".git", "node_modules", "test-results":
		return true
	default:
		return false
	}
}

// TestPhaseOnePackageSkeletonExists 确认 Phase 1 固定的分层目录仍然存在。
func TestPhaseOnePackageSkeletonExists(t *testing.T) {
	requiredDirs := []string{
		"architecture",
		"bootstrap",
		"httpapi",
		"facts",
		"execution",
		"product",
		"capabilities",
		filepath.Join("store", "sqlite"),
		"llm",
		"observability",
		filepath.Join("providers", "fobrain"),
	}

	root := filepath.Join("..")
	for _, dir := range requiredDirs {
		path := filepath.Join(root, dir)
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("required Phase 1 package directory %s missing: %v", path, err)
		}
		if !info.IsDir() {
			t.Fatalf("required Phase 1 package path %s is not a directory", path)
		}
	}
}

// TestImportBoundaryPreservesLayering 防止低层包反向依赖 HTTP/provider 等边界。
func TestImportBoundaryPreservesLayering(t *testing.T) {
	root := filepath.Join("..")
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		layer := packageLayer(root, path)
		for _, spec := range file.Imports {
			checkLayerImport(t, path, layer, spec)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestProductLayersDoNotImportEino 防止产品层直接消费 Eino 内部事件。
func TestProductLayersDoNotImportEino(t *testing.T) {
	root := filepath.Join("..")
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		layer := packageLayer(root, path)
		if !mustStayEinoFree(layer) {
			return nil
		}

		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		for _, spec := range file.Imports {
			importPath, err := strconv.Unquote(spec.Path.Value)
			if err != nil {
				t.Fatalf("%s has invalid import %s", path, spec.Path.Value)
			}
			if strings.HasPrefix(importPath, "github.com/cloudwego/eino") {
				t.Fatalf("%s layer %q must not import Eino package %q directly", path, layer, importPath)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// checkImport 检查单个 import 是否命中旧项目禁用标记。
func checkImport(t *testing.T, filePath string, spec *ast.ImportSpec, forbidden []string) {
	t.Helper()

	importPath, err := strconv.Unquote(spec.Path.Value)
	if err != nil {
		t.Fatalf("%s has invalid import %s", filePath, spec.Path.Value)
	}
	for _, marker := range forbidden {
		if strings.Contains(importPath, marker) {
			t.Fatalf("%s imports forbidden legacy boundary %q via %q", filePath, marker, importPath)
		}
	}
}

// packageLayer 从文件路径推导 internal/einoapp 下的逻辑层。
func packageLayer(root, filePath string) string {
	rel, err := filepath.Rel(root, filepath.Dir(filePath))
	if err != nil {
		return ""
	}
	parts := strings.Split(rel, string(filepath.Separator))
	if len(parts) == 0 {
		return ""
	}
	if parts[0] == "store" && len(parts) > 1 {
		return "store/" + parts[1]
	}
	if parts[0] == "providers" && len(parts) > 1 {
		return "providers/" + parts[1]
	}
	return parts[0]
}

// checkLayerImport 校验当前层是否导入了禁止的内部层。
func checkLayerImport(t *testing.T, filePath, layer string, spec *ast.ImportSpec) {
	t.Helper()

	importPath, err := strconv.Unquote(spec.Path.Value)
	if err != nil {
		t.Fatalf("%s has invalid import %s", filePath, spec.Path.Value)
	}
	if !strings.HasPrefix(importPath, modulePrefix) {
		return
	}
	importedLayer := strings.TrimPrefix(importPath, modulePrefix)
	importedLayer = strings.Split(importedLayer, "/")[0]
	if strings.HasPrefix(strings.TrimPrefix(importPath, modulePrefix), "store/sqlite") {
		importedLayer = "store/sqlite"
	}
	if strings.HasPrefix(strings.TrimPrefix(importPath, modulePrefix), "providers/fobrain") {
		importedLayer = "providers/fobrain"
	}

	forbidden := forbiddenImportsForLayer(layer)
	for _, blocked := range forbidden {
		if importedLayer == blocked {
			t.Fatalf("%s layer %q must not import %q via %q", filePath, layer, blocked, importPath)
		}
	}
}

// mustStayEinoFree 标记不能直接依赖 Eino 的产品出口层。
func mustStayEinoFree(layer string) bool {
	switch layer {
	case "httpapi", "facts", "product", "providers/fobrain":
		return true
	default:
		return false
	}
}

// forbiddenImportsForLayer 定义各层禁止依赖的内部层。
func forbiddenImportsForLayer(layer string) []string {
	switch layer {
	case "httpapi":
		return []string{"providers/fobrain", "llm", "capabilities"}
	case "execution":
		return []string{"httpapi", "providers/fobrain"}
	case "facts":
		return []string{"httpapi", "execution", "product", "capabilities", "providers/fobrain", "llm", "observability"}
	case "product":
		return []string{"httpapi", "execution", "providers/fobrain", "llm"}
	case "capabilities":
		return []string{"httpapi", "providers/fobrain"}
	case "llm":
		return []string{"httpapi", "execution", "providers/fobrain"}
	case "observability":
		return []string{"httpapi", "providers/fobrain"}
	case "providers/fobrain":
		return []string{"httpapi", "execution"}
	case "store/sqlite":
		return []string{"httpapi", "execution", "product", "capabilities", "providers/fobrain", "llm", "observability"}
	default:
		return nil
	}
}
