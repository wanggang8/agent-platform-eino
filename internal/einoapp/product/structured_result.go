package product

import (
	"errors"
	"strings"

	"agent-platform-eino/internal/einoapp/facts"
)

// StructuredResultSchemaVersion 是工具结果进入 Product Facts 的唯一 schema 版本。
const StructuredResultSchemaVersion = facts.StructuredResultSchemaVersion

// ErrInvalidStructuredResultCandidate 表示工具候选结果缺少 Product Facts 所需的安全字段。
var ErrInvalidStructuredResultCandidate = errors.New("invalid structured result candidate")

// StructuredResultCandidate 是 provider/raw 工具输出完成脱敏后的候选材料。
// 只有通过 Safety Gate 的候选结果才能转成 facts.StructuredResultRef。
type StructuredResultCandidate struct {
	SchemaVersion string
	ResultRef     string
	SafeSummary   string
}

// StructuredResultSafetyGate 校验工具候选结果是否可作为唯一事实材料进入 Product Facts。
type StructuredResultSafetyGate struct{}

// NewStructuredResultSafetyGate 创建 StructuredResult 安全门。
func NewStructuredResultSafetyGate() StructuredResultSafetyGate {
	return StructuredResultSafetyGate{}
}

// Approve 将安全候选结果收敛为 facts.StructuredResultRef，拒绝 raw payload、密钥和 raw checkpoint。
func (gate StructuredResultSafetyGate) Approve(candidate StructuredResultCandidate) (facts.StructuredResultRef, error) {
	schemaVersion := strings.TrimSpace(candidate.SchemaVersion)
	resultRef := strings.TrimSpace(candidate.ResultRef)
	safeSummary := strings.TrimSpace(candidate.SafeSummary)

	if schemaVersion != StructuredResultSchemaVersion || resultRef == "" || safeSummary == "" {
		return facts.StructuredResultRef{}, ErrInvalidStructuredResultCandidate
	}
	if ContainsUnsafeMaterial(resultRef) || ContainsUnsafeMaterial(safeSummary) {
		return facts.StructuredResultRef{}, facts.ErrUnsafeFactMaterial
	}

	return facts.StructuredResultRef{
		SchemaVersion: schemaVersion,
		ResultRef:     resultRef,
		SafeSummary:   safeSummary,
	}, nil
}
