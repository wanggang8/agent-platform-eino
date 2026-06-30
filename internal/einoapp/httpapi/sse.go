package httpapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SSEEvent struct {
	ID    string
	Event string
	Data  any
}

func SSEHeaders() http.Header {
	header := make(http.Header)
	header.Set("Content-Type", "text/event-stream; charset=utf-8")
	header.Set("Cache-Control", "no-cache")
	header.Set("X-Accel-Buffering", "no")
	return header
}

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
