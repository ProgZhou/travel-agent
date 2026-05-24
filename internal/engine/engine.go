package engine

import (
	"context"
	"log/slog"
	"travel-agent/internal/engine/prompt"
	"travel-agent/internal/llm"
	"travel-agent/internal/session"
	"travel-agent/internal/tool"
)

// Engine Agent 引擎接口
type Engine interface {
	// Run 运行 Agent，返回事件流 channel
	Run(ctx context.Context, sess *session.Session, userMessage string) (<-chan Event, error)
}

// DefaultEngine 默认 Agent 引擎实现
type DefaultEngine struct {
	llmClient          llm.Client
	toolRegistry       tool.Registry
	sessionStore       session.Store
	promptManager      *prompt.Manager
	maxReactIterations int
}

// NewEngine 创建 Agent 引擎
func NewEngine(
	llmClient llm.Client,
	toolRegistry tool.Registry,
	sessionStore session.Store,
	maxReactIterations int,
) *DefaultEngine {
	return &DefaultEngine{
		llmClient:          llmClient,
		toolRegistry:       toolRegistry,
		sessionStore:       sessionStore,
		promptManager:      prompt.NewManager(),
		maxReactIterations: maxReactIterations,
	}
}

// Run 运行 Agent
func (e *DefaultEngine) Run(ctx context.Context, sess *session.Session, userMessage string) (<-chan Event, error) {
	eventCh := make(chan Event, 10)

	go func() {
		defer close(eventCh)

		// 添加用户消息到历史
		sess.ChatHistory = append(sess.ChatHistory, session.ChatMessage{
			ID:        generateID(),
			Role:      "user",
			Content:   userMessage,
			Blocks:    []session.ContentBlock{{Type: "text", Content: userMessage}},
			Timestamp: nowMillis(),
		})

		// 执行 ReAct 循环
		executor := NewExecutor(e.llmClient, e.toolRegistry, e.promptManager, e.maxReactIterations)
		if err := executor.Execute(ctx, sess, eventCh); err != nil {
			slog.Error("executor failed", "error", err)
			eventCh <- Event{Type: "error", Data: ErrorEvent{Code: "20001", Message: "处理失败，请重试"}}
		}

		// 保存会话
		if err := e.sessionStore.Save(ctx, sess); err != nil {
			slog.Error("save session failed", "error", err)
		}

		// 发送完成事件
		eventCh <- Event{Type: "done", Data: map[string]interface{}{}}
	}()

	return eventCh, nil
}
