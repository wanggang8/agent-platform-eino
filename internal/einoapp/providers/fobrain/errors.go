package fobrain

import (
	"errors"
	"fmt"

	"agent-platform-eino/internal/einoapp/capabilities"
)

// PolicyReasonMissingCurrentUserScope 表示当前用户上下文不足，无法构造“我的范围”查询。
const PolicyReasonMissingCurrentUserScope capabilities.PolicyReason = "missing_current_user_scope"

// SafeError 是 Fobrain provider 边界可向上返回的脱敏错误。
type SafeError struct {
	ReasonCode  capabilities.PolicyReason
	SafeSummary string
}

// NewSafeError 创建稳定原因码错误，不包含 raw provider payload 或凭据。
func NewSafeError(reason capabilities.PolicyReason, summary string) SafeError {
	if summary == "" {
		summary = "Fobrain provider 暂不可用"
	}
	return SafeError{ReasonCode: reason, SafeSummary: summary}
}

// Error 返回可展示的安全摘要。
func (err SafeError) Error() string {
	return err.SafeSummary
}

// StableReasonCode 返回 policy/report 可记录的稳定原因码。
func (err SafeError) StableReasonCode() string {
	return string(err.ReasonCode)
}

// WrapSafeError 用于测试和边界包装，确保 errors.As 仍能识别稳定原因码。
func WrapSafeError(err SafeError) error {
	return fmt.Errorf("fobrain provider: %w", err)
}

// HasReason 判断错误链是否包含指定安全原因码。
func HasReason(err error, reason capabilities.PolicyReason) bool {
	var safeErr SafeError
	if errors.As(err, &safeErr) {
		return safeErr.ReasonCode == reason
	}
	var stable interface{ StableReasonCode() string }
	if errors.As(err, &stable) {
		return stable.StableReasonCode() == string(reason)
	}
	return false
}

// foldProviderError 将未知 provider 错误折叠成安全错误，避免上层拼接原始错误。
func foldProviderError(err error) error {
	if err == nil {
		return nil
	}
	var safeErr SafeError
	if errors.As(err, &safeErr) {
		return safeErr
	}
	var stable interface{ StableReasonCode() string }
	if errors.As(err, &stable) {
		return NewSafeError(capabilities.PolicyReason(stable.StableReasonCode()), "Fobrain provider 调用失败")
	}
	return NewSafeError(capabilities.PolicyReasonConnectorExecutionFailed, "Fobrain provider 调用失败")
}
