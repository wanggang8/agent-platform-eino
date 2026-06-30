package llm

import "context"

type Config struct {
	Provider          string
	BaseURL           string
	Model             string
	ModelLabel        string
	TimeoutMillis     int
	NetworkSafety     NetworkSafety
	CredentialBinding CredentialBinding
}

type NetworkSafety struct {
	RequireHTTPS         bool
	AllowLocalHTTP       bool
	BlockPrivateNetworks bool
	AllowRedirects       bool
	AllowedHosts         []string
}

type CredentialBinding struct {
	SchemaVersion string
	WorkspaceID   string
	System        string
	Status        string
	DisplayRef    string
	OwnerScope    string
	UpdatedAt     string
	AuditRef      string
}

type Provider interface {
	NewChatModel(ctx context.Context, cfg Config) (ChatModel, error)
}

type ChatModel interface {
	Generate(ctx context.Context, req ChatRequest) (ChatResponse, error)
	Stream(ctx context.Context, req ChatRequest) (ChatStream, error)
}

type ChatRequest struct {
	Messages []Message
}

type Message struct {
	Role    string
	Content string
}

type ChatResponse struct {
	Content string
}

type ChatStream interface {
	Recv() (ChatResponse, error)
	Close() error
}
