package capabilities

import (
	"context"
	"errors"

	"agent-platform-eino/internal/einoapp/product"
)

// ErrInvokerNotRegistered 表示 capability 已被请求但没有对应 provider invoker。
var ErrInvokerNotRegistered = errors.New("capability invoker not registered")

// InvokerMux 按 capability id 路由到注册时绑定的 provider invoker。
// 它避免 execution 根据工具名或 provider 名称写硬编码分支。
type InvokerMux struct {
	byCapabilityID map[string]Invoker
}

// NewInvokerMux 创建空 invoker 路由器。
func NewInvokerMux() *InvokerMux {
	return &InvokerMux{byCapabilityID: map[string]Invoker{}}
}

// RegisterProvider 把 provider 暴露的 capability 注册到 registry，并绑定对应 invoker。
func (mux *InvokerMux) RegisterProvider(registry *Registry, provider Provider, invoker Invoker) error {
	capabilities, err := provider.ListCapabilities()
	if err != nil {
		return err
	}
	for _, capability := range capabilities {
		if err := registry.Register(capability); err != nil {
			return err
		}
		mux.byCapabilityID[capability.ID] = invoker
	}
	return nil
}

// Invoke 根据 capability id 找到对应 provider；找不到时返回安全错误。
func (mux *InvokerMux) Invoke(ctx context.Context, request InvocationRequest) (product.StructuredResultCandidate, error) {
	invoker, ok := mux.byCapabilityID[request.CapabilityID]
	if !ok {
		return product.StructuredResultCandidate{}, ErrInvokerNotRegistered
	}
	return invoker.Invoke(ctx, request)
}
