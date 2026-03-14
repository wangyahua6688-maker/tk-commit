package httpx

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// NewTimeoutClient 创建带超时的 HTTP 客户端。
func NewTimeoutClient(timeout time.Duration) *http.Client {
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	return &http.Client{Timeout: timeout}
}

// GetRange 发起带 Range 头的 GET 请求并读取有限字节。
func GetRange(ctx context.Context, client *http.Client, url string, rangeHeader string, maxRead int64) (statusCode int, contentType string, body []byte, err error) {
	if client == nil {
		client = NewTimeoutClient(3 * time.Second)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, "", nil, err
	}
	if rangeHeader != "" {
		req.Header.Set("Range", rangeHeader)
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, "", nil, err
	}
	defer resp.Body.Close()

	if maxRead <= 0 {
		maxRead = 8192
	}
	body, err = io.ReadAll(io.LimitReader(resp.Body, maxRead))
	if err != nil {
		return resp.StatusCode, resp.Header.Get("Content-Type"), nil, fmt.Errorf("read body failed: %w", err)
	}
	return resp.StatusCode, resp.Header.Get("Content-Type"), body, nil
}
