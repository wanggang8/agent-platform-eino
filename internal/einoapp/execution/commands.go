package execution

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"agent-platform-eino/internal/einoapp/facts"
)

var ErrRunNotFound = errors.New("run not found")

type AcceptedRun struct {
	RunID  string
	Status string
}

type MessageCommand struct {
	WorkspaceID     string
	Message         string
	ClientRequestID string
	RunID           string
}

type ActionCommand struct {
	WorkspaceID     string
	ActionID        string
	ClientRequestID string
	CapabilityHint  string
	InputText       string
	Attachments     []Attachment
	Context         RequestContext
}

type ResumeCommand struct {
	WorkspaceID     string
	RunID           string
	ResumeRef       string
	ClientRequestID string
	Decision        string
	SelectedRefs    []string
	FreeText        string
	Comment         string
}

type Attachment struct {
	AttachmentRef string
	MediaType     string
	SafeName      string
}

type RequestContext struct {
	Timezone      string
	Locale        string
	SafeUserLabel string
}

type Commands interface {
	StartMessage(context.Context, MessageCommand) (AcceptedRun, error)
	StartAction(context.Context, ActionCommand) (AcceptedRun, error)
	Resume(context.Context, ResumeCommand) (AcceptedRun, error)
}

type StaticCommands struct {
	repository facts.Repository
	newRunID   func() (string, error)
}

func NewStaticCommands() StaticCommands {
	return StaticCommands{newRunID: randomRunID}
}

func NewFactCommands(repository facts.Repository) StaticCommands {
	return StaticCommands{
		repository: repository,
		newRunID:   randomRunID,
	}
}

func (commands StaticCommands) StartMessage(ctx context.Context, command MessageCommand) (AcceptedRun, error) {
	return commands.acceptRun(ctx, command.WorkspaceID, command.RunID)
}

func (commands StaticCommands) StartAction(ctx context.Context, command ActionCommand) (AcceptedRun, error) {
	return commands.acceptRun(ctx, command.WorkspaceID, "")
}

func (commands StaticCommands) Resume(ctx context.Context, command ResumeCommand) (AcceptedRun, error) {
	if commands.repository != nil {
		if _, err := commands.repository.GetRun(ctx, command.RunID); err != nil {
			if errors.Is(err, facts.ErrNotFound) {
				return AcceptedRun{}, ErrRunNotFound
			}
			return AcceptedRun{}, err
		}
	}
	return AcceptedRun{RunID: command.RunID, Status: "accepted"}, nil
}

func (commands StaticCommands) acceptRun(ctx context.Context, workspaceID string, requestedRunID string) (AcceptedRun, error) {
	runID := strings.TrimSpace(requestedRunID)
	if runID == "" {
		generated, err := commands.newRunID()
		if err != nil {
			return AcceptedRun{}, err
		}
		runID = generated
	}
	if commands.repository != nil {
		now := time.Now().UTC()
		if err := commands.repository.CreateRun(ctx, facts.Run{
			RunID:       runID,
			WorkspaceID: workspaceID,
			Status:      facts.RunStatusCreated,
			CreatedAt:   now,
			UpdatedAt:   now,
		}); err != nil {
			return AcceptedRun{}, err
		}
	}
	return AcceptedRun{RunID: runID, Status: "accepted"}, nil
}

func randomRunID() (string, error) {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return "run_" + hex.EncodeToString(bytes[:]), nil
}
