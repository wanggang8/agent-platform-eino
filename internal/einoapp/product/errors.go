package product

// SafeError 是产品/API 可展示错误，禁止携带 raw provider payload 或密钥。
type SafeError struct {
	Code           string
	Message        string
	SafeDetail     string
	Retryable      bool
	CorrelationRef string
}

// Error 返回安全错误详情。
func (err SafeError) Error() string {
	if err.SafeDetail != "" {
		return err.SafeDetail
	}
	return err.Message
}

// NewSafeError 创建统一安全错误。
func NewSafeError(code string, safeDetail string, retryable bool) SafeError {
	return SafeError{
		Code:       code,
		Message:    safeDetail,
		SafeDetail: safeDetail,
		Retryable:  retryable,
	}
}
