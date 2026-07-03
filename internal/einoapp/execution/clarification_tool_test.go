package execution_test

import (
	"errors"
	"testing"

	"agent-platform-eino/internal/einoapp/execution"
	"agent-platform-eino/internal/einoapp/facts"
)

func TestNewClarificationPendingEventBuildsSafeRunnerEvent(t *testing.T) {
	event, err := execution.NewClarificationPendingEvent(execution.ClarificationRequest{
		RunID:         "run-clarify",
		PendingID:     "pending-clarify-1",
		ResumeRef:     "resume-safe-clarify-1",
		CheckpointRef: "checkpoint_ref:clarify-1",
		Question:      "请选择要查询的人员",
		InputMode:     facts.PendingInputModeSingleChoice,
		Candidates: []facts.PendingCandidate{
			{
				CandidateRef: "candidate:fobrain:person:1",
				Label:        "张三",
				Description:  "安全部 / 安全运营",
				EntityType:   "person",
				SafeFields: []facts.PendingCandidateField{
					{Label: "部门", Value: "安全部"},
					{Label: "角色", Value: "安全运营"},
				},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if event.Kind != execution.RunnerEventPending || event.PendingKind != string(facts.PendingKindClarification) || event.PendingStatus != string(facts.PendingStatusWaiting) {
		t.Fatalf("runner event mismatch: %+v", event)
	}
	if event.Question != "请选择要查询的人员" || event.InputMode != string(facts.PendingInputModeSingleChoice) || len(event.Candidates) != 1 {
		t.Fatalf("clarification payload mismatch: %+v", event)
	}
}

func TestNewClarificationPendingEventRejectsUnsafeCandidate(t *testing.T) {
	for _, testCase := range []struct {
		name      string
		candidate facts.PendingCandidate
	}{
		{name: "email", candidate: facts.PendingCandidate{CandidateRef: "candidate:fobrain:person:1", Label: "张三", EntityType: "person", SafeFields: []facts.PendingCandidateField{{Label: "邮箱", Value: "zhangsan@example.com"}}}},
		{name: "formatted phone", candidate: facts.PendingCandidate{CandidateRef: "candidate:fobrain:person:1", Label: "张三", EntityType: "person", SafeFields: []facts.PendingCandidateField{{Label: "手机号", Value: "138-0013-8000"}}}},
		{name: "raw provider id", candidate: facts.PendingCandidate{CandidateRef: "candidate:fobrain:person:1", Label: "fb_user_736281", EntityType: "person"}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := execution.NewClarificationPendingEvent(execution.ClarificationRequest{
				RunID:         "run-clarify",
				PendingID:     "pending-clarify-1",
				ResumeRef:     "resume-safe-clarify-1",
				CheckpointRef: "checkpoint_ref:clarify-1",
				Question:      "请选择要查询的人员",
				InputMode:     facts.PendingInputModeSingleChoice,
				Candidates:    []facts.PendingCandidate{testCase.candidate},
			})
			if !errors.Is(err, facts.ErrUnsafeFactMaterial) {
				t.Fatalf("unsafe candidate err = %v", err)
			}
		})
	}
}

func TestNewClarificationPendingEventRejectsInvalidInputMode(t *testing.T) {
	_, err := execution.NewClarificationPendingEvent(execution.ClarificationRequest{
		RunID:         "run-clarify",
		PendingID:     "pending-clarify-1",
		ResumeRef:     "resume-safe-clarify-1",
		CheckpointRef: "checkpoint_ref:clarify-1",
		Question:      "请选择要查询的人员",
		InputMode:     facts.PendingInputMode("dropdown"),
		Candidates:    []facts.PendingCandidate{{CandidateRef: "candidate:fobrain:person:1", Label: "张三", EntityType: "person"}},
	})
	if !errors.Is(err, execution.ErrInvalidClarificationRequest) {
		t.Fatalf("invalid input mode err = %v", err)
	}
}
