package fobrain

import "context"

// FobrainClient 是 provider 边界内的业务读取接口。
// 实现可以是 mock 或 HTTP client，但 raw provider payload 不能离开该边界。
type FobrainClient interface {
	CurrentUserContext(context.Context, ResolvedCredential) (CurrentUserContextResult, error)
}

// CurrentUserContextResult 是当前用户 PoC 的安全业务结果。
type CurrentUserContextResult struct {
	DisplayName string
	Department  string
	Role        string
}

// MockClient 是 Phase 5 smoke 使用的本地 Fobrain client，不触达真实网络。
type MockClient struct {
	Result CurrentUserContextResult
}

// CurrentUserContext 返回安全 fixture 结果，用于证明 provider/Facts 链路。
func (client MockClient) CurrentUserContext(_ context.Context, _ ResolvedCredential) (CurrentUserContextResult, error) {
	if client.Result.DisplayName == "" {
		return CurrentUserContextResult{
			DisplayName: "Fobrain 本地用户",
			Department:  "安全运营",
			Role:        "只读验证",
		}, nil
	}
	return client.Result, nil
}
