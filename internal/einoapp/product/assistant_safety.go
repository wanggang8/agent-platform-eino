package product

import (
	"errors"
	"strings"

	"agent-platform-eino/internal/einoapp/facts"
)

// ErrInvalidAssistantMessage 表示 assistant 候选文本为空，不能写入 Product Facts。
var ErrInvalidAssistantMessage = errors.New("invalid assistant message")

// AssistantSafetyGate 校验 assistant 文本是否可写入 Turn fact。
type AssistantSafetyGate struct{}

// NewAssistantSafetyGate 创建 assistant 文本安全门。
func NewAssistantSafetyGate() AssistantSafetyGate {
	return AssistantSafetyGate{}
}

// Approve 返回可写入 Product Facts 的 assistant 文本，拒绝明显敏感或 raw provider 内容。
func (gate AssistantSafetyGate) Approve(content string) (string, error) {
	approved := strings.TrimSpace(content)
	if approved == "" {
		return "", ErrInvalidAssistantMessage
	}
	if ContainsUnsafeMaterial(approved) {
		return "", facts.ErrUnsafeFactMaterial
	}
	return approved, nil
}
