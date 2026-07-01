package facts

import "strings"

// StructuredResultSchemaVersion 是 Product Facts 接受的工具结果 schema 版本。
const StructuredResultSchemaVersion = "tool.structured_result.v1"

// ContainsUnsafeMaterial 判断文本是否包含明显不能进入 Product Facts 的敏感或原始材料。
// 这是 facts/store 的最后防线；更复杂的安全投影仍由 product Safety Gate 承担。
func ContainsUnsafeMaterial(value string) bool {
	lower := strings.ToLower(value)
	for _, marker := range unsafeMaterialMarkers {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

// UnsafeStructuredResultRef 判断 StructuredResult 引用是否不满足 facts 安全边界。
func UnsafeStructuredResultRef(result StructuredResultRef) bool {
	if strings.TrimSpace(result.SchemaVersion) != StructuredResultSchemaVersion {
		return true
	}
	return ContainsUnsafeMaterial(result.ResultRef) || ContainsUnsafeMaterial(result.SafeSummary)
}

// unsafeMaterialMarkers 是 facts 最后一层拒绝的基础风险标记，避免 raw payload 或密钥入库。
var unsafeMaterialMarkers = []string{
	"authorization:",
	"bearer ",
	"api_key",
	"apikey",
	"token",
	"password",
	"credential",
	"secret",
	"raw provider",
	"raw body",
	"raw error",
	"raw_payload",
	"provider_payload",
	"checkpoint-raw",
	"interrupt-raw",
}
