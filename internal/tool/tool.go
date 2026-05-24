package tool

import (
	"context"
	"encoding/json"
	"travel-agent/internal/llm"
)

// Tool 工具接口
type Tool interface {
	// Name 工具名称
	Name() string
	// Description 工具描述
	Description() string
	// Parameters 返回 JSON Schema 格式的参数描述
	Parameters() json.RawMessage
	// Execute 执行工具调用
	Execute(ctx context.Context, args json.RawMessage) (*ToolResult, error)
}

// ToolResult 工具执行结果
type ToolResult struct {
	Data   interface{} `json:"data"`
	IsMock bool        `json:"is_mock"`
	Source string      `json:"source"` // "real_api" | "mock" | "fallback"
	Notice string      `json:"notice"`
	Error  string      `json:"error"`
}

// Registry 工具注册表
type Registry interface {
	// Register 注册工具
	Register(tool Tool)
	// Get 获取工具
	Get(name string) (Tool, bool)
	// List 列出所有工具
	List() []Tool
	// GetSchemas 获取所有工具的 Schema（用于 LLM Function Calling）
	GetSchemas() []llm.ToolSchema
}
