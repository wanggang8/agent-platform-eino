package capabilities

import (
	"context"
	"errors"
	"fmt"

	"agent-platform-eino/internal/einoapp/facts"
	"agent-platform-eino/internal/einoapp/product"
)

// ErrCapabilityInvocationUnavailable 表示 provider 无法执行指定能力。
var ErrCapabilityInvocationUnavailable = errors.New("capability invocation unavailable")

// InvocationRequest 是 tool adapter 传给 provider 的安全调用请求。
type InvocationRequest struct {
	CapabilityID  string
	Arguments     map[string]any
	PolicyContext PolicyContext
}

// Invoker 是 execution 调用业务 provider 的最小接口。
type Invoker interface {
	Invoke(context.Context, InvocationRequest) (product.StructuredResultCandidate, error)
}

// MockProvider 是 Phase 4 本地只读能力提供方，用于证明 tool loop 和 Product Facts 链路。
type MockProvider struct {
	id           string
	capabilities map[string]Capability
}

var _ Provider = (*MockProvider)(nil)
var _ Invoker = (*MockProvider)(nil)

// NewMockProvider 创建配置驱动的 mock provider；capability 元数据仍由注册表决定。
func NewMockProvider(id string, capabilities []Capability) *MockProvider {
	byID := make(map[string]Capability, len(capabilities))
	for _, capability := range capabilities {
		byID[capability.ID] = capability
	}
	return &MockProvider{id: id, capabilities: byID}
}

// ID 返回 provider id。
func (provider *MockProvider) ID() string {
	return provider.id
}

// ListCapabilities 返回 mock provider 注册的能力元数据。
func (provider *MockProvider) ListCapabilities() ([]Capability, error) {
	out := make([]Capability, 0, len(provider.capabilities))
	for _, capability := range provider.capabilities {
		out = append(out, capability)
	}
	return out, nil
}

// Invoke 返回安全 StructuredResult candidate，不包含请求原文或 provider raw payload。
func (provider *MockProvider) Invoke(_ context.Context, request InvocationRequest) (product.StructuredResultCandidate, error) {
	capability, ok := provider.capabilities[request.CapabilityID]
	if !ok {
		return product.StructuredResultCandidate{}, ErrCapabilityInvocationUnavailable
	}
	if capability.RiskLevel != RiskReadOnly {
		return product.StructuredResultCandidate{}, ErrCapabilityInvocationUnavailable
	}
	return product.StructuredResultCandidate{
		SchemaVersion: facts.StructuredResultSchemaVersion,
		ResultRef:     fmt.Sprintf("result:%s", capability.ID),
		SafeSummary:   fmt.Sprintf("%s 完成", capability.DisplayName),
	}, nil
}
