package capabilities

import (
	"errors"
	"slices"
)

// Registry 保存运行期可用能力，是工具选择的唯一生产入口。
type Registry struct {
	byID map[string]Capability
}

// NewRegistry 创建空能力注册表。
func NewRegistry() *Registry {
	return &Registry{byID: map[string]Capability{}}
}

// Register 注册单个能力，并校验最小元数据完整性。
func (registry *Registry) Register(capability Capability) error {
	if capability.ID == "" {
		return errors.New("capability id is required")
	}
	if capability.ProviderID == "" {
		return errors.New("capability provider id is required")
	}
	if capability.ToolName == "" {
		return errors.New("capability tool name is required")
	}
	if _, exists := registry.byID[capability.ID]; exists {
		return errors.New("capability id already registered")
	}
	registry.byID[capability.ID] = capability
	return nil
}

// List 返回按 capability id 排序的能力列表，保证测试和投影稳定。
func (registry *Registry) List() []Capability {
	out := make([]Capability, 0, len(registry.byID))
	for _, capability := range registry.byID {
		out = append(out, capability)
	}
	slices.SortFunc(out, func(a, b Capability) int {
		if a.ID < b.ID {
			return -1
		}
		if a.ID > b.ID {
			return 1
		}
		return 0
	})
	return out
}

// Get 按 capability id 查找能力，供 execution selection 使用。
func (registry *Registry) Get(id string) (Capability, bool) {
	capability, ok := registry.byID[id]
	return capability, ok
}
