package capabilities

import (
	"context"
	"errors"
	"strings"

	"github.com/cloudwego/eino/schema"
)

// ErrInvalidToolSchema 表示 capability 输入 schema 无法转换为 Eino ToolInfo。
var ErrInvalidToolSchema = errors.New("invalid tool schema")

// ToolInfoFromCapability 将注册表中的能力元数据转换为 Eino ToolInfo。
// 这里不按工具名称分支，tool name、description 和参数只能来自 Capability。
func ToolInfoFromCapability(_ context.Context, capability Capability) (*schema.ToolInfo, error) {
	if strings.TrimSpace(capability.ToolName) == "" || strings.TrimSpace(capability.Description) == "" {
		return nil, ErrInvalidToolSchema
	}

	info := &schema.ToolInfo{
		Name: capability.ToolName,
		Desc: capability.Description,
		Extra: map[string]any{
			"capability_id": capability.ID,
			"provider_id":   capability.ProviderID,
			"risk_level":    string(capability.RiskLevel),
		},
	}
	if len(capability.InputSchema.Properties) > 0 {
		params := make(map[string]*schema.ParameterInfo, len(capability.InputSchema.Properties))
		required := requiredSchemaFields(capability.InputSchema)
		for name, dataType := range capability.InputSchema.Properties {
			parameterType, err := toEinoDataType(dataType)
			if err != nil {
				return nil, err
			}
			params[name] = &schema.ParameterInfo{
				Type:     parameterType,
				Required: required[name],
			}
		}
		info.ParamsOneOf = schema.NewParamsOneOfByParams(params)
	}
	return info, nil
}

// requiredSchemaFields 返回 schema required 集合；nil 表示旧 capability 未声明 required，显式空切片表示全部可选。
func requiredSchemaFields(input JSONSchema) map[string]bool {
	required := map[string]bool{}
	if input.Required == nil {
		for name := range input.Properties {
			required[name] = true
		}
		return required
	}
	for _, name := range input.Required {
		if _, ok := input.Properties[name]; ok {
			required[name] = true
		}
	}
	return required
}

// toEinoDataType 将项目轻量 schema 类型映射为 Eino 参数类型。
func toEinoDataType(dataType string) (schema.DataType, error) {
	switch strings.TrimSpace(dataType) {
	case "string":
		return schema.String, nil
	case "number":
		return schema.Number, nil
	case "integer":
		return schema.Integer, nil
	case "boolean":
		return schema.Boolean, nil
	default:
		return "", ErrInvalidToolSchema
	}
}
