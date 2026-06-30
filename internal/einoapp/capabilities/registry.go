package capabilities

import (
	"errors"
	"slices"
)

type Registry struct {
	byID map[string]Capability
}

func NewRegistry() *Registry {
	return &Registry{byID: map[string]Capability{}}
}

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
