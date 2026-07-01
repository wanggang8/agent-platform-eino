package llm

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"slices"
)

// validateNetworkPolicy 校验模型 provider 出站目标是否符合配置派生的安全策略。
func validateNetworkPolicy(baseURL string, policy NetworkSafety) (string, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", redactedProviderError("config", "provider_config_invalid", "模型地址配置无效", false)
	}
	if parsed.User != nil {
		return "", redactedProviderError("config", "provider_config_invalid", "模型地址不能包含 userinfo", false)
	}
	host := parsed.Hostname()
	if policy.RequireHTTPS && parsed.Scheme != "https" {
		return "", redactedProviderError("network", "provider_insecure_scheme_blocked", "模型网络策略阻止了非 HTTPS 地址", false)
	}
	if parsed.Scheme == "http" && !(policy.AllowLocalHTTP && isLoopbackHost(host)) {
		return "", redactedProviderError("network", "provider_insecure_scheme_blocked", "模型网络策略阻止了 HTTP 地址", false)
	}
	if len(policy.AllowedHosts) > 0 && !slices.Contains(policy.AllowedHosts, host) {
		return "", redactedProviderError("network", "provider_host_not_allowed", "模型网络策略阻止了未授权主机", false)
	}
	if policy.BlockPrivateNetworks && isPrivateHost(host) && !(policy.AllowLocalHTTP && isLoopbackHost(host)) {
		return "", redactedProviderError("network", "provider_network_blocked", "模型网络策略阻止了私有网络地址", false)
	}
	return host, nil
}

// policyHTTPClient 创建受网络策略约束的 HTTP client，覆盖 redirect 和 dial 边界。
func policyHTTPClient(base *http.Client, policy NetworkSafety) *http.Client {
	client := &http.Client{}
	if base != nil {
		*client = *base
	}
	client.CheckRedirect = func(req *http.Request, _ []*http.Request) error {
		if !policy.AllowRedirects {
			return redactedProviderError("network", "provider_redirect_blocked", "模型网络策略阻止了重定向", false)
		}
		if _, err := validateNetworkPolicy(req.URL.String(), policy); err != nil {
			return err
		}
		return nil
	}
	client.Transport = policyTransport(policy)
	return client
}

// policyTransport 在 DNS 解析和 dial 前后都执行私网阻断，避免只校验字符串 host。
func policyTransport(policy NetworkSafety) http.RoundTripper {
	dialer := &net.Dialer{}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.Proxy = nil
	transport.DialContext = func(ctx context.Context, network string, address string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return nil, redactedProviderError("network", "provider_config_invalid", "模型网络地址无效", false)
		}
		ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
		if err != nil {
			return nil, redactedProviderError("network", "provider_unavailable", "模型网络解析失败", true)
		}
		for _, ip := range ips {
			if policyBlocksIP(ip.IP, policy) {
				return nil, redactedProviderError("network", "provider_network_blocked", "模型网络策略阻止了私有网络地址", false)
			}
		}
		if len(ips) == 0 {
			return nil, redactedProviderError("network", "provider_unavailable", "模型网络解析失败", true)
		}
		return dialer.DialContext(ctx, network, net.JoinHostPort(ips[0].IP.String(), port))
	}
	return transport
}

// isLoopbackHost 判断本地开发可显式允许的 loopback 地址。
func isLoopbackHost(host string) bool {
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

// isPrivateHost 判断明显私有或链路本地 IP；域名解析检查留给真实 dialer 阶段扩展。
func isPrivateHost(host string) bool {
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	return ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()
}

// policyBlocksIP 判断 DNS/dial 解析出的目标 IP 是否被当前策略阻断。
func policyBlocksIP(ip net.IP, policy NetworkSafety) bool {
	if ip == nil || !policy.BlockPrivateNetworks {
		return false
	}
	if policy.AllowLocalHTTP && ip.IsLoopback() {
		return false
	}
	return ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified()
}

// redactedProviderError 创建 provider 边界安全错误，避免调用方拼接原始错误。
func redactedProviderError(category string, reasonCode string, safeSummary string, retryable bool) RedactedError {
	return RedactProviderError(nil, RedactionInput{
		Category:    category,
		ReasonCode:  reasonCode,
		SafeSummary: safeSummary,
		Retryable:   retryable,
	})
}
