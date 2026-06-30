package llm

import "fmt"

type RedactionInput struct {
	Category       string
	ReasonCode     string
	SafeSummary    string
	Retryable      bool
	CorrelationRef string
}

type RedactedError struct {
	SchemaVersion  string           `json:"schema_version"`
	Category       string           `json:"category"`
	ReasonCode     string           `json:"reason_code"`
	SafeSummary    string           `json:"safe_summary"`
	Retryable      bool             `json:"retryable"`
	CorrelationRef string           `json:"correlation_ref"`
	Redaction      RedactionSummary `json:"redaction"`
}

type RedactionSummary struct {
	RemovedSecret     bool `json:"removed_secret"`
	RemovedAuthHeader bool `json:"removed_auth_header"`
	RemovedRawBody    bool `json:"removed_raw_body"`
	RemovedRawPrompt  bool `json:"removed_raw_prompt"`
}

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

func (err RedactedError) Error() string {
	return fmt.Sprintf("%s: %s", err.ReasonCode, err.SafeSummary)
}
