package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
	"travel-agent/internal/provider"
)

// WeatherTool 天气查询工具
type WeatherTool struct {
	provider      provider.WeatherProvider
	timeout       time.Duration
	retryCount    int
	retryInterval time.Duration
}

// NewWeatherTool 创建天气查询工具
func NewWeatherTool(p provider.WeatherProvider, timeout time.Duration, retryCount int, retryInterval time.Duration) *WeatherTool {
	return &WeatherTool{
		provider:      p,
		timeout:       timeout,
		retryCount:    retryCount,
		retryInterval: retryInterval,
	}
}

// Name 工具名称
func (t *WeatherTool) Name() string {
	return "weather_query"
}

// Description 工具描述
func (t *WeatherTool) Description() string {
	return "查询目的城市在指定日期范围内的天气情况，返回每日天气和汇总概况"
}

// Parameters 参数 Schema
func (t *WeatherTool) Parameters() json.RawMessage {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"city": map[string]interface{}{
				"type":        "string",
				"description": "城市名称，如\"三亚\"",
			},
			"start_date": map[string]interface{}{
				"type":        "string",
				"description": "开始日期，格式 YYYY-MM-DD",
			},
			"end_date": map[string]interface{}{
				"type":        "string",
				"description": "结束日期，格式 YYYY-MM-DD",
			},
		},
		"required": []string{"city", "start_date", "end_date"},
	}
	data, _ := json.Marshal(schema)
	return data
}

// Execute 执行工具调用
func (t *WeatherTool) Execute(ctx context.Context, args json.RawMessage) (*ToolResult, error) {
	var req provider.WeatherQueryRequest
	if err := json.Unmarshal(args, &req); err != nil {
		return nil, fmt.Errorf("parse args: %w", err)
	}

	slog.Info("weather_query called", "city", req.City, "start_date", req.StartDate, "end_date", req.EndDate)

	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	var result *provider.WeatherQueryResult
	var lastErr error
	for attempt := 0; attempt <= t.retryCount; attempt++ {
		if attempt > 0 {
			slog.Warn("retrying weather_query", "attempt", attempt)
			time.Sleep(t.retryInterval)
		}

		var err error
		result, err = t.provider.Query(ctx, req)
		if err == nil {
			break
		}
		lastErr = err
	}

	if lastErr != nil {
		slog.Error("weather_query failed", "error", lastErr)
		return &ToolResult{
			Data:   nil,
			IsMock: true,
			Source: "fallback",
			Notice: "预估数据，仅供参考",
			Error:  lastErr.Error(),
		}, nil
	}

	slog.Info("weather_query succeeded", "daily_count", len(result.Daily), "is_mock", result.IsMock)

	source := "real_api"
	notice := ""
	if result.IsMock {
		source = "mock"
		notice = "预估数据，仅供参考"
	}

	return &ToolResult{
		Data:   result,
		IsMock: result.IsMock,
		Source: source,
		Notice: notice,
	}, nil
}
