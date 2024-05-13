package apis

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	netUrl "net/url"
)

func DownloadFileFromUrl(ctx context.Context, url string) ([]byte, error) {
	body, err := HttpRequest(
		ctx,
		"GET",
		url,
		nil,
		nil,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func HttpRequest(
	ctx context.Context,
	method, url string,
	requestBody []byte,
	headers map[string]string,
	query map[string]string,
) ([]byte, error) {
	// 添加请求query参数
	if query != nil && len(query) > 0 {
		queryParams := netUrl.Values{}
		for key, value := range query {
			queryParams.Add(key, value)
		}
		url = url + "?" + queryParams.Encode()
	}
	// 创建 HTTP 请求
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	// 添加 HTTP 请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	req = req.WithContext(ctx)
	// 发送 HTTP 请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	// 读取 HTTP 响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// 检查 HTTP 响应状态码是否为 2xx
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d body:%s", resp.StatusCode, string(body))
	}
	return body, nil
}
