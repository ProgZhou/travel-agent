package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
	"travel-agent/internal/provider"
)

// HotelTool 酒店查询工具
type HotelTool struct {
	provider      provider.HotelProvider
	timeout       time.Duration
	retryCount    int
	retryInterval time.Duration
}

// NewHotelTool 创建酒店查询工具
func NewHotelTool(p provider.HotelProvider, timeout time.Duration, retryCount int, retryInterval time.Duration) *HotelTool {
	return &HotelTool{
		provider:      p,
		timeout:       timeout,
		retryCount:    retryCount,
		retryInterval: retryInterval,
	}
}

// Name 工具名称
func (t *HotelTool) Name() string {
	return "hotel_search"
}

// Description 工具描述
func (t *HotelTool) Description() string {
	return "搜索目的城市的酒店信息，返回可用酒店列表（含名称、价格、位置、房型）"
}

// Parameters 参数 Schema
func (t *HotelTool) Parameters() json.RawMessage {
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"city": map[string]interface{}{
				"type":        "string",
				"description": "城市名称，如\"三亚\"",
			},
			"check_in": map[string]interface{}{
				"type":        "string",
				"description": "入住日期，格式 YYYY-MM-DD",
			},
			"check_out": map[string]interface{}{
				"type":        "string",
				"description": "离店日期，格式 YYYY-MM-DD",
			},
			"level": map[string]interface{}{
				"type":        "string",
				"enum":        []string{"economy", "comfort", "luxury"},
				"description": "酒店档次：economy=经济型, comfort=舒适型, luxury=豪华型",
			},
		},
		"required": []string{"city", "check_in", "check_out"},
	}
	data, _ := json.Marshal(schema)
	return data
}

// Execute 执行工具调用
func (t *HotelTool) Execute(ctx context.Context, args json.RawMessage) (*ToolResult, error) {
	var req provider.HotelSearchRequest
	if err := json.Unmarshal(args, &req); err != nil {
		return nil, fmt.Errorf("parse args: %w", err)
	}

	slog.Info("hotel_search called", "city", req.City, "check_in", req.CheckIn, "check_out", req.CheckOut, "level", req.Level)

	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	var result *provider.HotelSearchResult
	var lastErr error
	for attempt := 0; attempt <= t.retryCount; attempt++ {
		if attempt > 0 {
			slog.Warn("retrying hotel_search", "attempt", attempt)
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
		slog.Error("hotel_search failed", "error", lastErr)
		return &ToolResult{
			Data:   nil,
			IsMock: true,
			Source: "fallback",
			Notice: "预估数据，仅供参考",
			Error:  lastErr.Error(),
		}, nil
	}

	slog.Info("hotel_search succeeded", "hotels", len(result.Hotels), "is_mock", result.IsMock)

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
