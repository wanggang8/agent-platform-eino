package observability

// TokenCostRates 是配置化 token 单价，单位为 micro-units/token；默认 0 表示只统计 token 不估算成本。
type TokenCostRates struct {
	InputMicrounitsPerToken  int64
	OutputMicrounitsPerToken int64
}

// Estimate 根据安全 token 计数估算成本，只返回统计值，不触发预算限制或 lifecycle。
func (rates TokenCostRates) Estimate(inputTokens int, outputTokens int) int64 {
	inputCost := nonNegativeInt64(rates.InputMicrounitsPerToken) * int64(nonNegativeInt(inputTokens))
	outputCost := nonNegativeInt64(rates.OutputMicrounitsPerToken) * int64(nonNegativeInt(outputTokens))
	return inputCost + outputCost
}

func nonNegativeInt(value int) int {
	if value < 0 {
		return 0
	}
	return value
}

func nonNegativeInt64(value int64) int64 {
	if value < 0 {
		return 0
	}
	return value
}
