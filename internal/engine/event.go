package engine

import "encoding/json"

// Event Agent Engine 向外推送的事件
type Event struct {
	Type string      `json:"type"` // thinking | text | tool_call | card | phase_change | done | error
	Data interface{} `json:"data"`
}

// TextEvent 文本事件数据
type TextEvent struct {
	Content string `json:"content"`
	Delta   bool   `json:"delta"`
}

// ThinkingEvent 思考事件数据
type ThinkingEvent struct {
	Content string `json:"content"`
}

// ToolCallEvent 工具调用事件数据
type ToolCallEvent struct {
	Name   string          `json:"name"`
	Args   json.RawMessage `json:"args"`
	Status string          `json:"status"` // running | done | error
	Result interface{}     `json:"result"`
}

// CardEvent 卡片事件数据
type CardEvent struct {
	CardType string      `json:"card_type"`
	Data     interface{} `json:"data"`
}

// PhaseChangeEvent 阶段变更事件数据
type PhaseChangeEvent struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// ErrorEvent 错误事件数据
type ErrorEvent struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
