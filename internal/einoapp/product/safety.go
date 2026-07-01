package product

import "agent-platform-eino/internal/einoapp/facts"

// ContainsUnsafeMaterial 判断文本是否包含明显不能进入产品事实的敏感或原始材料。
// 该检查是产品层 Safety Gate 的轻量规则，后续可由 policy/红线配置扩展。
func ContainsUnsafeMaterial(value string) bool {
	return facts.ContainsUnsafeMaterial(value)
}
