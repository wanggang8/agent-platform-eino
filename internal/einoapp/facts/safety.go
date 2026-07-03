package facts

import (
	"regexp"
	"strings"
)

// StructuredResultSchemaVersion 是 Product Facts 接受的工具结果 schema 版本。
const StructuredResultSchemaVersion = "tool.structured_result.v1"

// CheckpointRefPrefix 是 Product Facts 可保存的 checkpoint 安全引用前缀。
const CheckpointRefPrefix = "checkpoint_ref:"

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

// SafeCheckpointRef 判断 checkpoint_ref 是否为产品安全引用，而不是 Eino 内部 checkpoint id。
func SafeCheckpointRef(value string) bool {
	ref := strings.TrimSpace(value)
	return ref == value &&
		strings.HasPrefix(ref, CheckpointRefPrefix) &&
		len(ref) > len(CheckpointRefPrefix) &&
		!ContainsUnsafeMaterial(ref)
}

// UnsafePendingCandidates 判断 clarification 候选是否包含原始 id、凭据、手机号、邮箱等不安全材料。
func UnsafePendingCandidates(candidates []PendingCandidate) bool {
	for _, candidate := range candidates {
		if strings.TrimSpace(candidate.CandidateRef) == "" ||
			strings.TrimSpace(candidate.Label) == "" ||
			!validPendingCandidateEntityType(candidate.EntityType) {
			return true
		}
		if unsafeCandidateText(candidate.CandidateRef) ||
			unsafeCandidateText(candidate.Label) ||
			unsafeCandidateText(candidate.Description) ||
			unsafeCandidateText(candidate.EntityType) {
			return true
		}
		for _, field := range candidate.SafeFields {
			if unsafeCandidateText(field.Label) || unsafeCandidateText(field.Value) {
				return true
			}
		}
	}
	return false
}

// validPendingCandidateEntityType 保持 facts 候选与 pending interaction schema 的 entity_type 枚举一致。
func validPendingCandidateEntityType(entityType string) bool {
	switch strings.TrimSpace(entityType) {
	case "person", "asset", "vulnerability", "department", "business":
		return true
	default:
		return false
	}
}

// unsafeCandidateText 是候选字段专用的更严格规则，弥补通用 facts 标记无法识别邮箱/手机号的问题。
func unsafeCandidateText(value string) bool {
	text := strings.TrimSpace(value)
	lower := strings.ToLower(text)
	if ContainsUnsafeMaterial(lower) {
		return true
	}
	if emailPattern.MatchString(text) || longDigitPattern.MatchString(text) {
		return true
	}
	if looksLikePhoneNumber(text) {
		return true
	}
	for _, marker := range unsafeCandidateMarkers {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return false
}

// looksLikePhoneNumber 识别带空格、短横线等格式化分隔符的常见电话号码。
func looksLikePhoneNumber(value string) bool {
	digits := digitOnly(value)
	if len(digits) == 11 && strings.HasPrefix(digits, "1") {
		return true
	}
	if len(digits) >= 10 && (strings.Contains(value, "-") || strings.Contains(value, " ") || strings.Contains(value, "(")) {
		return true
	}
	return false
}

func digitOnly(value string) string {
	var builder strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' {
			builder.WriteRune(r)
		}
	}
	return builder.String()
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

var unsafeCandidateMarkers = []string{
	"raw",
	"internal",
	"employee",
	"fb_user",
	"user_",
	"u_",
	"staff",
	"phone",
	"mobile",
	"email",
	"authorization",
	"_id",
	"id:",
	"id=",
	"内部",
	"手机号",
	"手机",
	"电话",
	"员工号",
	"工号",
	"邮箱",
}

var (
	emailPattern     = regexp.MustCompile(`(?i)[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}`)
	longDigitPattern = regexp.MustCompile(`\d{8,}`)
)
