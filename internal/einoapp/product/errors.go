package product

type SafeError struct {
	Code           string
	Message        string
	SafeDetail     string
	Retryable      bool
	CorrelationRef string
}

func (err SafeError) Error() string {
	if err.SafeDetail != "" {
		return err.SafeDetail
	}
	return err.Message
}

func NewSafeError(code string, safeDetail string, retryable bool) SafeError {
	return SafeError{
		Code:       code,
		Message:    safeDetail,
		SafeDetail: safeDetail,
		Retryable:  retryable,
	}
}
