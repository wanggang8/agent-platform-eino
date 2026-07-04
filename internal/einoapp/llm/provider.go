package llm

import "context"

// Config 是 LLM provider 的安全配置摘要，真实密钥只通过凭据绑定解析。
type Config struct {
	Provider          string
	BaseURL           string
	Model             string
	ModelLabel        string
	TimeoutMillis     int
	NetworkSafety     NetworkSafety
	CredentialBinding CredentialBinding
}

// NetworkSafety 定义模型请求的出站网络安全策略。
type NetworkSafety struct {
	RequireHTTPS         bool
	AllowLocalHTTP       bool
	BlockPrivateNetworks bool
	AllowRedirects       bool
	AllowedHosts         []string
}

// CredentialBinding 是模型凭据的安全显示引用，不包含 token。
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

// Provider 创建具体 ChatModel，实现必须负责网络策略和错误脱敏。
type Provider interface {
	NewChatModel(ctx context.Context, cfg Config) (ChatModel, error)
}

// ChatModel 是 execution 使用的最小模型接口，隔离具体 SDK。
type ChatModel interface {
	Generate(ctx context.Context, req ChatRequest) (ChatResponse, error)
	Stream(ctx context.Context, req ChatRequest) (ChatStream, error)
}

// ChatRequest 是安全上下文投影后的模型输入。
type ChatRequest struct {
	Messages []Message
}

// Message 是传给模型的单条安全消息。
type Message struct {
	Role    string
	Content string
}

// ChatResponse 是模型返回的安全文本候选。
type ChatResponse struct {
	Content string
	Usage   TokenUsage
}

// TokenUsage 是 provider 返回的安全 token 计数，不包含 prompt、completion 或 provider raw body。
type TokenUsage struct {
	InputTokens  int
	OutputTokens int
	TotalTokens  int
}

// ChatStream 抽象模型流式输出，调用方必须负责关闭。
type ChatStream interface {
	Recv() (ChatResponse, error)
	Close() error
}
