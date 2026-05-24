package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"travel-agent/internal/engine/prompt"
	"travel-agent/internal/llm"
	"travel-agent/internal/session"
	"travel-agent/internal/tool"
)

// Executor ReAct 执行器
type Executor struct {
	llmClient     llm.Client
	toolRegistry  tool.Registry
	promptManager *prompt.Manager
	maxIterations int
}

// NewExecutor 创建执行器
func NewExecutor(llmClient llm.Client, toolRegistry tool.Registry, promptManager *prompt.Manager, maxIterations int) *Executor {
	return &Executor{
		llmClient:     llmClient,
		toolRegistry:  toolRegistry,
		promptManager: promptManager,
		maxIterations: maxIterations,
	}
}

// Execute 执行 ReAct 循环
func (e *Executor) Execute(ctx context.Context, sess *session.Session, eventCh chan<- Event) error {
	// 构建消息历史
	messages := e.buildMessages(sess)

	// ReAct 循环
	for i := 0; i < e.maxIterations; i++ {
		slog.Debug("react iteration", "iteration", i+1, "max", e.maxIterations)

		// 调用 LLM（流式）
		req := llm.ChatRequest{
			Messages: messages,
			Tools:    e.toolRegistry.GetSchemas(),
			Stream:   true,
		}

		chunkCh, err := e.llmClient.ChatCompletionStream(ctx, req)
		if err != nil {
			return fmt.Errorf("llm stream failed: %w", err)
		}

		// 处理流式响应
		assistantMsg, hasToolCalls := e.processStream(ctx, chunkCh, eventCh)
		messages = append(messages, assistantMsg)

		// 如果没有工具调用，说明 LLM 已经给出最终回复
		if !hasToolCalls {
			break
		}

		// 执行工具调用
		toolMessages := e.executeTools(ctx, assistantMsg.ToolCalls, eventCh)
		messages = append(messages, toolMessages...)
	}

	return nil
}

// buildMessages 构建消息历史
func (e *Executor) buildMessages(sess *session.Session) []llm.LLMMessage {
	messages := []llm.LLMMessage{
		{
			Role:    "system",
			Content: e.promptManager.BuildSystemPrompt(sess),
		},
	}

	// 添加历史消息（最近 10 条）
	start := 0
	if len(sess.ChatHistory) > 10 {
		start = len(sess.ChatHistory) - 10
	}
	for _, msg := range sess.ChatHistory[start:] {
		messages = append(messages, llm.LLMMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return messages
}

// processStream 处理流式响应
func (e *Executor) processStream(ctx context.Context, chunkCh <-chan llm.StreamChunk, eventCh chan<- Event) (llm.LLMMessage, bool) {
	var assistantMsg llm.LLMMessage
	assistantMsg.Role = "assistant"
	var contentBuilder strings.Builder
	var reasoningBuilder strings.Builder
	hasToolCalls := false

	for {
		select {
		case <-ctx.Done():
			return assistantMsg, hasToolCalls
		case chunk, ok := <-chunkCh:
			if !ok {
				assistantMsg.Content = contentBuilder.String()
				return assistantMsg, hasToolCalls
			}

			if chunk.Error != nil {
				eventCh <- Event{Type: "error", Data: ErrorEvent{Code: "20002", Message: chunk.Error.Message}}
				return assistantMsg, hasToolCalls
			}

			if len(chunk.Choices) == 0 {
				continue
			}

			delta := chunk.Choices[0].Delta

			// 思考内容
			if delta.ReasoningContent != "" {
				reasoningBuilder.WriteString(delta.ReasoningContent)
				eventCh <- Event{Type: "thinking", Data: ThinkingEvent{Content: delta.ReasoningContent}}
			}

			// 文本内容
			if delta.Content != "" {
				contentBuilder.WriteString(delta.Content)
				eventCh <- Event{Type: "text", Data: TextEvent{Content: delta.Content, Delta: true}}
			}

			// 工具调用
			if len(delta.ToolCalls) > 0 {
				hasToolCalls = true
				assistantMsg.ToolCalls = append(assistantMsg.ToolCalls, delta.ToolCalls...)
				for _, tc := range delta.ToolCalls {
					eventCh <- Event{Type: "tool_call", Data: ToolCallEvent{
						Name:   tc.Function.Name,
						Args:   json.RawMessage(tc.Function.Arguments),
						Status: "running",
					}}
				}
			}
		}
	}
}

// executeTools 执行工具调用
func (e *Executor) executeTools(ctx context.Context, toolCalls []llm.ToolCall, eventCh chan<- Event) []llm.LLMMessage {
	messages := []llm.LLMMessage{}

	for _, tc := range toolCalls {
		t, ok := e.toolRegistry.Get(tc.Function.Name)
		if !ok {
			slog.Warn("tool not found", "name", tc.Function.Name)
			continue
		}

		result, err := t.Execute(ctx, json.RawMessage(tc.Function.Arguments))
		if err != nil {
			slog.Error("tool execution failed", "tool", tc.Function.Name, "error", err)
			eventCh <- Event{Type: "tool_call", Data: ToolCallEvent{
				Name:   tc.Function.Name,
				Status: "error",
				Result: err.Error(),
			}}
			continue
		}

		eventCh <- Event{Type: "tool_call", Data: ToolCallEvent{
			Name:   tc.Function.Name,
			Status: "done",
			Result: result,
		}}

		// 将工具结果添加到消息历史
		resultJSON, _ := json.Marshal(result)
		messages = append(messages, llm.LLMMessage{
			Role:       "tool",
			Content:    string(resultJSON),
			ToolCallID: tc.ID,
		})
	}

	return messages
}
