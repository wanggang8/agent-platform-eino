package llm

import "fmt"

// RedactionInput 是 provider 错误进入安全错误前的人工摘要输入。
type RedactionInput struct {
	Category       string
	ReasonCode     string
	SafeSummary    string
	Retryable      bool
	CorrelationRef string
}

// RedactedError 是可进入 Product Facts/API 的 provider 安全错误。
type RedactedError struct {
	SchemaVersion  string           `json:"schema_version"`
	Category       string           `json:"category"`
	ReasonCode     string           `json:"reason_code"`
	SafeSummary    string           `json:"safe_summary"`
	Retryable      bool             `json:"retryable"`
	CorrelationRef string           `json:"correlation_ref"`
	Redaction      RedactionSummary `json:"redaction"`
}

// RedactionSummary 记录已执行的脱敏类别，便于 audit/replay 展示。
type RedactionSummary struct {
	RemovedSecret     bool `json:"removed_secret"`
	RemovedAuthHeader bool `json:"removed_auth_header"`
	RemovedRawBody    bool `json:"removed_raw_body"`
	RemovedRawPrompt  bool `json:"removed_raw_prompt"`
}

// RedactProviderError 丢弃原始错误细节，仅保留安全摘要和脱敏标记。
func RedactProviderError(_ error, input RedactionInput) RedactedError {
	return RedactedError{
		SchemaVersion:  "eino.provider_redacted_error.v1",
		Category:       input.Category,
		ReasonCode:     input.ReasonCode,
		SafeSummary:    input.SafeSummary,
		Retryable:      input.Retryable,
		CorrelationRef: input.CorrelationRef,
		Redaction: RedactionSummary{
			RemovedSecret:     true,
			RemovedAuthHeader: true,
			RemovedRawBody:    true,
		},
	}
}

// Error 返回安全摘要，不包含原始 provider body 或密钥。
func (err RedactedError) Error() string {
	return fmt.Sprintf("%s: %s", err.ReasonCode, err.SafeSummary)
}
