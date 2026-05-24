package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// DoubaoClient 豆包 LLM 客户端
type DoubaoClient struct {
	apiURL        string
	apiKey        string
	defaultModel  string
	maxTokens     int
	temperature   float64
	timeout       time.Duration
	maxRetries    int
	retryInterval time.Duration
	httpClient    *http.Client
}

// NewDoubaoClient 创建豆包客户端
func NewDoubaoClient(apiURL, apiKey, model string, maxTokens int, temperature float64, timeout time.Duration, maxRetries int, retryInterval time.Duration) *DoubaoClient {
	return &DoubaoClient{
		apiURL:        apiURL,
		apiKey:        apiKey,
		defaultModel:  model,
		maxTokens:     maxTokens,
		temperature:   temperature,
		timeout:       timeout,
		maxRetries:    maxRetries,
		retryInterval: retryInterval,
		httpClient:    &http.Client{Timeout: timeout},
	}
}

// ChatCompletion 非流式调用
func (c *DoubaoClient) ChatCompletion(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	req.Stream = false
	if req.Model == "" {
		req.Model = c.defaultModel
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = c.maxTokens
	}
	if req.Temperature == 0 {
		req.Temperature = c.temperature
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			slog.Warn("retrying LLM call", "attempt", attempt, "max_retries", c.maxRetries)
			time.Sleep(c.retryInterval)
		}

		resp, err := c.doRequest(ctx, req)
		if err == nil {
			return resp, nil
		}
		lastErr = err
	}

	return nil, fmt.Errorf("LLM call failed after %d retries: %w", c.maxRetries, lastErr)
}

// ChatCompletionStream 流式调用
func (c *DoubaoClient) ChatCompletionStream(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error) {
	req.Stream = true
	if req.Model == "" {
		req.Model = c.defaultModel
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = c.maxTokens
	}
	if req.Temperature == 0 {
		req.Temperature = c.temperature
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.apiURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		httpResp.Body.Close()
		return nil, fmt.Errorf("HTTP %d: %s", httpResp.StatusCode, string(body))
	}

	chunkCh := make(chan StreamChunk, 10)
	go c.readStream(httpResp.Body, chunkCh)

	return chunkCh, nil
}

// doRequest 执行非流式请求
func (c *DoubaoClient) doRequest(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.apiURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	slog.Debug("LLM request", "model", req.Model, "messages", len(req.Messages), "tools", len(req.Tools))

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", httpResp.StatusCode, string(body))
	}

	var resp ChatResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	slog.Debug("LLM response", "id", resp.ID, "finish_reason", resp.Choices[0].FinishReason, "tokens", resp.Usage.TotalTokens)

	return &resp, nil
}

// readStream 读取 SSE 流
func (c *DoubaoClient) readStream(body io.ReadCloser, chunkCh chan<- StreamChunk) {
	defer close(chunkCh)
	defer body.Close()

	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}

		var chunk StreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			slog.Warn("failed to parse stream chunk", "error", err, "data", data)
			continue
		}

		chunkCh <- chunk
	}

	if err := scanner.Err(); err != nil {
		slog.Error("stream read error", "error", err)
		chunkCh <- StreamChunk{Error: &StreamError{Code: -1, Message: err.Error()}}
	}
}
