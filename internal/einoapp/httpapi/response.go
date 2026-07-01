package httpapi

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"agent-platform-eino/internal/einoapp/product"
)

const requestIDHeader = "X-Request-Id"

// APIError 是 HTTP 层对产品安全错误的别名，避免引入第二套错误模型。
type APIError = product.SafeError

// errorEnvelope 是唯一错误响应包裹，成功响应保持业务 schema 直出。
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

// WriteJSON 写成功业务响应，不额外包一层 data。
func WriteJSON(w http.ResponseWriter, r *http.Request, status int, body any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set(requestIDHeader, requestID(r))
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// WriteError 写统一安全错误 envelope，确保对外只暴露 safe detail。
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

// requestID 复用入站 request id；没有时生成 opaque id。
func requestID(r *http.Request) string {
	if r == nil {
		return newRequestID()
	}
	if existing := r.Header.Get(requestIDHeader); existing != "" {
		return existing
	}
	return newRequestID()
}

// newRequestID 生成不包含业务信息的请求关联 id。
func newRequestID() string {
	var raw [12]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "req_unavailable"
	}
	return "req_" + hex.EncodeToString(raw[:])
}
