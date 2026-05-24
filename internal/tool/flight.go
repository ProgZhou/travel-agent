package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
	"travel-agent/internal/provider"
)

// FlightTool 航班查询工具
type FlightTool struct {
	provider      provider.FlightProvider
	timeout       time.Duration
	retryCount    int
	retryInterval time.Duration
}

// NewFlightTool 创建航班查询工具
func NewFlightTool(p provider.FlightProvider, timeout time.Duration, retryCount int, retryInterval time.Duration) *FlightTool {
	return &FlightTool{
		provider:      p,
		timeout:       timeout,
		retryCount:    retryCount,
		retryInterval: retryInterval,
	}
}

// Name 工具名称
func (t *FlightTool) Name() string {
	return "flight_search"
}

// Description 工具描述
func (t *FlightTool) Description() string {
	return "搜索从出发城市到目的城市的航班或火车票信息，返回可用班次列表（含时间、价格、座位等级）"
}

// Parameters 参数 Schema
func (t *FlightTool) Parameters() json.RawMessage {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"depart_city": map[string]interface{}{
				"type":        "string",
				"description": "出发城市名称，如\"北京\"",
			},
			"arrive_city": map[string]interface{}{
				"type":        "string",
				"description": "到达城市名称，如\"三亚\"",
			},
			"date": map[string]interface{}{
				"type":        "string",
				"description": "出发日期，格式 YYYY-MM-DD",
			},
			"type": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"flight", "train", "any"},
				"description": "交通类型：flight=飞机, train=火车, any=不限",
			},
		},
		"required": []string{"depart_city", "arrive_city", "date"},
	}
	data, _ := json.Marshal(schema)
	return data
}

// Execute 执行工具调用
func (t *FlightTool) Execute(ctx context.Context, args json.RawMessage) (*ToolResult, error) {
	var req provider.FlightSearchRequest
	if err := json.Unmarshal(args, &req); err != nil {
		return nil, fmt.Errorf("parse args: %w", err)
	}

	slog.Info("flight_search called", "depart", req.DepartCity, "arrive", req.ArriveCity, "date", req.Date, "type", req.Type)

	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	// 执行查询（带重试）
	var result *provider.FlightSearchResult
	var lastErr error
	for attempt := 0; attempt <= t.retryCount; attempt++ {
		if attempt > 0 {
			slog.Warn("retrying flight_search", "attempt", attempt)
			time.Sleep(t.retryInterval)
		}

		var err error
		result, err = t.provider.Search(ctx, req)
		if err == nil {
			break
		}
		lastErr = err
	}

	if lastErr != nil {
		slog.Error("flight_search failed", "error", lastErr)
		return &ToolResult{
			Data:   nil,
			IsMock: true,
			Source: "fallback",
			Notice: "预估数据，仅供参考",
			Error:  lastErr.Error(),
		}, nil
	}

	slog.Info("flight_search succeeded", "tickets", len(result.Tickets), "is_mock", result.IsMock)

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
