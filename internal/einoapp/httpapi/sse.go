package httpapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SSEEvent 是 HTTP SSE 编码输入，data 必须已经是产品投影事件。
type SSEEvent struct {
	ID    string
	Event string
	Data  any
}

// SSEHeaders 返回防代理缓冲的标准 SSE 响应头。
func SSEHeaders() http.Header {
	header := make(http.Header)
	header.Set("Content-Type", "text/event-stream; charset=utf-8")
	header.Set("Cache-Control", "no-cache")
	header.Set("X-Accel-Buffering", "no")
	return header
}

// EncodeSSEEvent 将事件编码到任意 writer，便于单元测试。
func EncodeSSEEvent(w io.Writer, event SSEEvent) error {
	data, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	if event.ID != "" {
		if _, err := fmt.Fprintf(w, "id: %s\n", event.ID); err != nil {
			return err
		}
	}
	if event.Event != "" {
		if _, err := fmt.Fprintf(w, "event: %s\n", event.Event); err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(w, "data: %s\n\n", data)
	return err
}

// WriteSSEEvent 写单个 SSE 事件并在支持时 flush。
func WriteSSEEvent(w http.ResponseWriter, event SSEEvent) error {
	data, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	for key, values := range SSEHeaders() {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(http.StatusOK)
	if err := writeEncodedSSEEvent(w, event, data); err != nil {
		return err
	}
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}

// writeEncodedSSEEvent 复用已经 JSON 编码的数据，避免重复编码失败路径不一致。
func writeEncodedSSEEvent(w io.Writer, event SSEEvent, data []byte) error {
	if event.ID != "" {
		if _, err := fmt.Fprintf(w, "id: %s\n", event.ID); err != nil {
			return err
		}
	}
	if event.Event != "" {
		if _, err := fmt.Fprintf(w, "event: %s\n", event.Event); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "data: %s\n\n", data)
	return err
}
