package httpapi

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"agent-platform-eino/internal/einoapp/product"
)

const requestIDHeader = "X-Request-Id"

type APIError = product.SafeError

type errorEnvelope struct {
	SchemaVersion string      `json:"schema_version"`
	RequestID     string      `json:"request_id"`
	Error         errorDetail `json:"error"`
}

type errorDetail struct {
	Code           string `json:"code"`
	Message        string `json:"message"`
	SafeDetail     string `json:"safe_detail"`
	Retryable      bool   `json:"retryable"`
	CorrelationRef string `json:"correlation_ref,omitempty"`
}

func WriteJSON(w http.ResponseWriter, r *http.Request, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set(requestIDHeader, requestID(r))
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, apiErr APIError) {
	if apiErr.Code == "" {
		apiErr.Code = "internal_error"
	}
	if apiErr.Message == "" {
		apiErr.Message = "请求处理失败"
	}
	if apiErr.SafeDetail == "" {
		apiErr.SafeDetail = apiErr.Message
	}
	reqID := requestID(r)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set(requestIDHeader, reqID)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorEnvelope{
		SchemaVersion: "eino_error_envelope.v1",
		RequestID:     reqID,
		Error: errorDetail{
			Code:           apiErr.Code,
			Message:        apiErr.Message,
			SafeDetail:     apiErr.SafeDetail,
			Retryable:      apiErr.Retryable,
			CorrelationRef: apiErr.CorrelationRef,
		},
	})
}

func requestID(r *http.Request) string {
	if r == nil {
		return newRequestID()
	}
	if existing := r.Header.Get(requestIDHeader); existing != "" {
		return existing
	}
	return newRequestID()
}

func newRequestID() string {
	var raw [12]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "req_unavailable"
	}
	return "req_" + hex.EncodeToString(raw[:])
}
