package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

// messageRequest 对应 Workbench 消息入口，只接受契约字段。
type messageRequest struct {
	SchemaVersion   string `json:"schema_version"`
	Message         string `json:"message"`
	ClientRequestID string `json:"client_request_id"`
	RunID           string `json:"run_id,omitempty"`
	WorkspaceID     string `json:"workspace_id,omitempty"`
}

// actionRequest 对应外部 Action API，请求体不会直接进入 provider。
type actionRequest struct {
	SchemaVersion   string        `json:"schema_version"`
	ActionID        string        `json:"action_id"`
	ClientRequestID string        `json:"client_request_id"`
	CapabilityHint  string        `json:"capability_hint,omitempty"`
	Input           actionInput   `json:"input"`
	Context         actionContext `json:"context,omitempty"`
}

type actionInput struct {
	Text        string             `json:"text"`
	Attachments []actionAttachment `json:"attachments,omitempty"`
}

type actionAttachment struct {
	AttachmentRef string `json:"attachment_ref"`
	MediaType     string `json:"media_type"`
	SafeName      string `json:"safe_name,omitempty"`
}

type actionContext struct {
	Timezone      string `json:"timezone,omitempty"`
	Locale        string `json:"locale,omitempty"`
	SafeUserLabel string `json:"safe_user_label,omitempty"`
}

// resumeRequest 对应 approval/clarification 恢复入口，resume_ref 必须是一次性安全引用。
type resumeRequest struct {
	SchemaVersion         string   `json:"schema_version"`
	ResumeRef             string   `json:"resume_ref"`
	ClientRequestID       string   `json:"client_request_id"`
	Decision              string   `json:"decision,omitempty"`
	SelectedCandidateRefs []string `json:"selected_candidate_refs,omitempty"`
	FreeText              string   `json:"free_text,omitempty"`
	Comment               string   `json:"comment,omitempty"`
}

// lifecycleRequest 对应 run 生命周期控制入口，不接受 provider 或工具原始材料。
type lifecycleRequest struct {
	SchemaVersion   string `json:"schema_version"`
	Action          string `json:"action"`
	ClientRequestID string `json:"client_request_id"`
}

// messageResponse 是消息入口的接收结果，后续事件仍通过 Product Facts/SSE 投影读取。
type messageResponse struct {
	SchemaVersion string `json:"schema_version"`
	RunID         string `json:"run_id"`
	Status        string `json:"status"`
	StreamURL     string `json:"stream_url"`
	SnapshotURL   string `json:"snapshot_url,omitempty"`
}

// handleHealth 返回服务健康状态，不触碰业务依赖。
func (api api) handleHealth(w http.ResponseWriter, r *http.Request) {
	WriteJSON(w, r, http.StatusOK, map[string]string{
		"schema_version": "eino_workbench_health.v1",
		"status":         "ok",
	})
}

// handleCurrentView 从 product projection 读取当前 Workbench 视图。
func (api api) handleCurrentView(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	view, err := api.projection.WorkbenchView(r.Context(), workspaceID)
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("projection_failed", "工作台视图暂不可用", true))
		return
	}
	WriteJSON(w, r, http.StatusOK, view)
}

// handleMessage 校验 Workbench 消息契约并交给 execution command。
func (api api) handleMessage(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	var req messageRequest
	if err := decodeJSONRequest(r, &req); err != nil {
		WriteError(w, r, http.StatusBadRequest, product.NewSafeError("invalid_request", "消息请求无效", false))
		return
	}
	if req.SchemaVersion != "eino_workbench_message_request.v1" || strings.TrimSpace(req.Message) == "" || strings.TrimSpace(req.ClientRequestID) == "" {
		WriteError(w, r, http.StatusBadRequest, product.NewSafeError("invalid_request", "消息请求无效", false))
		return
	}
	accepted, err := api.commands.StartMessage(r.Context(), execution.MessageCommand{
		WorkspaceID:     workspaceID,
		Message:         req.Message,
		ClientRequestID: req.ClientRequestID,
		RunID:           req.RunID,
	})
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("execution_failed", "消息暂无法执行", true))
		return
	}
	WriteJSON(w, r, http.StatusOK, messageResponse{
		SchemaVersion: "eino_workbench_message_response.v1",
		RunID:         accepted.RunID,
		Status:        accepted.Status,
		StreamURL:     fmt.Sprintf("/api/workspaces/%s/runs/%s/stream", workspaceID, accepted.RunID),
		SnapshotURL:   fmt.Sprintf("/api/workspaces/%s/runs/%s", workspaceID, accepted.RunID),
	})
}

// handleAgentAction 校验 Action API 请求，并复用 Workbench 同源 facts/projection。
func (api api) handleAgentAction(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	var req actionRequest
	if err := decodeJSONRequest(r, &req); err != nil {
		WriteError(w, r, http.StatusBadRequest, product.NewSafeError("invalid_request", "操作请求无效", false))
		return
	}
	if req.SchemaVersion != "eino_action_request.v1" || strings.TrimSpace(req.ActionID) == "" || strings.TrimSpace(req.ClientRequestID) == "" || strings.TrimSpace(req.Input.Text) == "" {
		WriteError(w, r, http.StatusBadRequest, product.NewSafeError("invalid_request", "操作请求无效", false))
		return
	}
	accepted, err := api.commands.StartAction(r.Context(), execution.ActionCommand{
		WorkspaceID:     workspaceID,
		ActionID:        req.ActionID,
		ClientRequestID: req.ClientRequestID,
		CapabilityHint:  req.CapabilityHint,
		InputText:       req.Input.Text,
		Attachments:     toExecutionAttachments(req.Input.Attachments),
		Context: execution.RequestContext{
			Timezone:      req.Context.Timezone,
			Locale:        req.Context.Locale,
			SafeUserLabel: req.Context.SafeUserLabel,
		},
	})
	if err != nil {
		if errors.Is(err, execution.ErrCapabilityNotRegistered) {
			WriteError(w, r, http.StatusBadRequest, product.NewSafeError("capability_not_registered", "请求的能力未注册", false))
			return
		}
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("execution_failed", "操作暂无法执行", true))
		return
	}
	result, err := api.projection.ActionResult(r.Context(), workspaceID, req.ActionID, accepted.RunID)
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("projection_failed", "操作结果暂不可用", true))
		return
	}
	WriteJSON(w, r, http.StatusOK, result)
}

// handleResume 校验 resume 请求，完整恢复语义由 execution/facts 层控制。
func (api api) handleResume(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	runID := r.PathValue("run_id")
	var req resumeRequest
	if err := decodeJSONRequest(r, &req); err != nil {
		WriteError(w, r, http.StatusBadRequest, product.NewSafeError("invalid_request", "恢复请求无效", false))
		return
	}
	if !validResumeRequest(req) {
		WriteError(w, r, http.StatusBadRequest, product.NewSafeError("invalid_request", "恢复请求无效", false))
		return
	}
	accepted, err := api.commands.Resume(r.Context(), execution.ResumeCommand{
		WorkspaceID:     workspaceID,
		RunID:           runID,
		ResumeRef:       req.ResumeRef,
		ClientRequestID: req.ClientRequestID,
		Decision:        req.Decision,
		SelectedRefs:    req.SelectedCandidateRefs,
		FreeText:        req.FreeText,
		Comment:         req.Comment,
	})
	if err != nil {
		if errors.Is(err, execution.ErrRunNotFound) {
			WriteError(w, r, http.StatusNotFound, product.NewSafeError("run_not_found", "运行不存在或不可恢复", false))
			return
		}
		if errors.Is(err, execution.ErrCheckpointMissing) {
			WriteError(w, r, http.StatusConflict, product.NewSafeError("checkpoint_missing", "恢复点不存在或已失效", false))
			return
		}
		if errors.Is(err, execution.ErrResumeNotAllowed) {
			WriteError(w, r, http.StatusConflict, product.NewSafeError("resume_not_allowed", "当前运行状态不允许恢复", false))
			return
		}
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("execution_failed", "恢复暂无法执行", true))
		return
	}
	result, err := api.projection.ResumeResult(r.Context(), workspaceID, accepted.RunID)
	if err != nil {
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("projection_failed", "恢复结果暂不可用", true))
		return
	}
	WriteJSON(w, r, http.StatusOK, result)
}

// handleRunLifecycle 校验生命周期控制请求，状态迁移由 execution/facts 层完成。
func (api api) handleRunLifecycle(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	runID := r.PathValue("run_id")
	var req lifecycleRequest
	if err := decodeJSONRequest(r, &req); err != nil {
		WriteError(w, r, http.StatusBadRequest, product.NewSafeError("invalid_request", "生命周期请求无效", false))
		return
	}
	if !validLifecycleRequest(req) {
		WriteError(w, r, http.StatusBadRequest, product.NewSafeError("invalid_request", "生命周期请求无效", false))
		return
	}
	action := execution.LifecycleAction(req.Action)
	accepted, err := api.commands.RunLifecycle(r.Context(), execution.LifecycleCommand{
		WorkspaceID:     workspaceID,
		RunID:           runID,
		Action:          action,
		ClientRequestID: req.ClientRequestID,
	})
	if err != nil {
		if errors.Is(err, execution.ErrRunNotFound) {
			WriteError(w, r, http.StatusNotFound, product.NewSafeError("run_not_found", "运行不存在或不可操作", false))
			return
		}
		if errors.Is(err, execution.ErrLifecycleNotAllowed) {
			WriteError(w, r, http.StatusConflict, product.NewSafeError("lifecycle_not_allowed", "当前运行状态不允许该操作", false))
			return
		}
		if errors.Is(err, facts.ErrIdempotencyConflict) {
			WriteError(w, r, http.StatusConflict, product.NewSafeError("lifecycle_conflict", "生命周期操作已过期或已处理", false))
			return
		}
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("execution_failed", "生命周期操作暂无法执行", true))
		return
	}
	result, err := api.projection.ActionResult(r.Context(), workspaceID, "run_lifecycle:"+string(action), accepted.RunID)
	if err != nil {
		if errors.Is(err, facts.ErrNotFound) {
			WriteError(w, r, http.StatusNotFound, product.NewSafeError("run_not_found", "运行不存在", false))
			return
		}
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("projection_failed", "生命周期结果暂不可用", true))
		return
	}
	WriteJSON(w, r, http.StatusOK, result)
}

// handleRunSnapshot 返回指定 run 的产品快照视图。
func (api api) handleRunSnapshot(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	runID := r.PathValue("run_id")
	view, err := api.projection.RunSnapshot(r.Context(), workspaceID, runID)
	if err != nil {
		if errors.Is(err, facts.ErrNotFound) {
			WriteError(w, r, http.StatusNotFound, product.NewSafeError("run_not_found", "运行不存在", false))
			return
		}
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("projection_failed", "运行快照暂不可用", true))
		return
	}
	WriteJSON(w, r, http.StatusOK, view)
}

// handleReplay 返回基于 Product Facts 的回放视图。
func (api api) handleReplay(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	runID := r.PathValue("run_id")
	replay, err := api.projection.ReplayView(r.Context(), workspaceID, runID)
	if err != nil {
		if errors.Is(err, facts.ErrNotFound) {
			WriteError(w, r, http.StatusNotFound, product.NewSafeError("run_not_found", "运行不存在", false))
			return
		}
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("projection_failed", "回放暂不可用", true))
		return
	}
	WriteJSON(w, r, http.StatusOK, replay)
}

// handleRunStream 从产品投影读取 SSE 事件，HTTP 层不接触 Eino event 或 provider payload。
func (api api) handleRunStream(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.PathValue("workspace_id")
	runID := r.PathValue("run_id")
	events, err := api.projection.StreamEvents(r.Context(), workspaceID, runID)
	if err != nil {
		if errors.Is(err, facts.ErrNotFound) {
			WriteError(w, r, http.StatusNotFound, product.NewSafeError("run_not_found", "运行不存在", false))
			return
		}
		WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("projection_failed", "事件暂不可用", true))
		return
	}

	w.Header().Set(requestIDHeader, requestID(r))
	for _, event := range events {
		if err := WriteSSEEvent(w, SSEEvent{
			ID:    event.EventID,
			Event: event.Type,
			Data:  event,
		}); err != nil {
			WriteError(w, r, http.StatusInternalServerError, product.NewSafeError("sse_encode_failed", "事件编码失败", true))
			return
		}
	}
}

// decodeJSONRequest 使用 DisallowUnknownFields 防止静默接受未知契约字段。
func decodeJSONRequest(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		return err
	}
	var extra json.RawMessage
	if err := decoder.Decode(&extra); err != io.EOF {
		return fmt.Errorf("request body must contain a single JSON object")
	}
	return nil
}

func toExecutionAttachments(attachments []actionAttachment) []execution.Attachment {
	if len(attachments) == 0 {
		return nil
	}
	result := make([]execution.Attachment, 0, len(attachments))
	for _, attachment := range attachments {
		result = append(result, execution.Attachment{
			AttachmentRef: attachment.AttachmentRef,
			MediaType:     attachment.MediaType,
			SafeName:      attachment.SafeName,
		})
	}
	return result
}

func validResumeRequest(req resumeRequest) bool {
	if req.SchemaVersion != "eino_workbench_resume_request.v1" || strings.TrimSpace(req.ResumeRef) == "" || strings.TrimSpace(req.ClientRequestID) == "" {
		return false
	}
	hasDecision := req.Decision == "approve" || req.Decision == "reject"
	hasCandidates := len(req.SelectedCandidateRefs) > 0
	hasFreeText := strings.TrimSpace(req.FreeText) != ""
	count := 0
	for _, ok := range []bool{hasDecision, hasCandidates, hasFreeText} {
		if ok {
			count++
		}
	}
	return count == 1
}

func validLifecycleRequest(req lifecycleRequest) bool {
	if req.SchemaVersion != "eino_run_lifecycle_request.v1" || strings.TrimSpace(req.ClientRequestID) == "" {
		return false
	}
	switch execution.LifecycleAction(req.Action) {
	case execution.LifecycleActionCancel, execution.LifecycleActionStop, execution.LifecycleActionProviderTimeout, execution.LifecycleActionPendingTimeout, execution.LifecycleActionRetry:
		return true
	default:
		return false
	}
}
