package services

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

const defaultHTTPTimeout = 30 * time.Second

// resilientHTTPClient 在检测到环境代理不可用时自动回退到直连。
type resilientHTTPClient struct {
	primary     *http.Client
	direct      *http.Client
	forceDirect bool
}

func newResilientHTTPClient(forceDirect bool) (*resilientHTTPClient, error) {
	primary, err := newHTTPClient(forceDirect)
	if err != nil {
		return nil, err
	}

	direct, err := newHTTPClient(true)
	if err != nil {
		return nil, err
	}

	return &resilientHTTPClient{
		primary:     primary,
		direct:      direct,
		forceDirect: forceDirect,
	}, nil
}

func newHTTPClient(disableProxy bool) (*http.Client, error) {
	baseTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, fmt.Errorf("无法复制默认HTTP传输配置")
	}

	transport := baseTransport.Clone()
	if disableProxy {
		transport.Proxy = nil
	}

	return &http.Client{
		Transport: transport,
		Timeout:   defaultHTTPTimeout,
	}, nil
}

func (c *resilientHTTPClient) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.primary.Do(req)
	if err == nil || c.forceDirect || !isProxyConnectionError(err) {
		return resp, err
	}

	retryReq, cloneErr := cloneRequest(req)
	if cloneErr != nil {
		return nil, fmt.Errorf("%w；检测到代理连接异常，但无法自动切换为直连重试: %v", err, cloneErr)
	}

	resp, directErr := c.direct.Do(retryReq)
	if directErr == nil {
		return resp, nil
	}

	return nil, fmt.Errorf("%w；检测到代理连接异常，已自动切换为直连重试但仍失败: %v", err, directErr)
}

func cloneRequest(req *http.Request) (*http.Request, error) {
	cloned := req.Clone(req.Context())
	if req.Body == nil {
		return cloned, nil
	}

	if req.GetBody == nil {
		return nil, fmt.Errorf("请求体不可重复读取")
	}

	body, err := req.GetBody()
	if err != nil {
		return nil, fmt.Errorf("重新创建请求体失败: %w", err)
	}

	cloned.Body = body
	return cloned, nil
}

func isProxyConnectionError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "proxyconnect tcp") ||
		strings.Contains(msg, "proxy connection") ||
		strings.Contains(msg, "http proxy error")
}
