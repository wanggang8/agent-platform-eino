package capabilities

import "time"

type RiskLevel string

const (
	RiskReadOnly RiskLevel = "read_only"
	RiskWrite    RiskLevel = "write"
)

type Capability struct {
	ID               string
	ProviderID       string
	ToolName         string
	DisplayName      string
	Description      string
	InputSchema      JSONSchema
	ResultSchema     string
	RiskLevel        RiskLevel
	ApprovalRequired bool
	Timeout          time.Duration
}

type JSONSchema struct {
	SchemaVersion string
	Properties    map[string]string
}

type Provider interface {
	ID() string
	ListCapabilities() ([]Capability, error)
}
