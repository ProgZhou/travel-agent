package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"travel-agent/internal/engine"
	"travel-agent/internal/pkg"
	"travel-agent/internal/session"
)

// ChatHandler 聊天处理器
type ChatHandler struct {
	engine engine.Engine
	store  session.Store
}

// NewChatHandler 创建聊天处理器
func NewChatHandler(eng engine.Engine, store session.Store) *ChatHandler {
	return &ChatHandler{
		engine: eng,
		store:  store,
	}
}

// ChatRequest 聊天请求
type ChatRequest struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

// Handle 处理聊天请求
func (h *ChatHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, pkg.NewBadRequest("invalid request body"))
		return
	}

	// 校验参数
	if req.SessionID == "" {
		writeError(w, pkg.NewBadRequest("session_id is required"))
		return
	}
	if strings.TrimSpace(req.Message) == "" {
		writeError(w, pkg.NewMessageInvalid("message cannot be empty"))
		return
	}
	if len(req.Message) > 2000 {
		writeError(w, pkg.NewMessageInvalid("message too long (max 2000 characters)"))
		return
	}

	// 获取会话
	sess, err := h.store.Get(r.Context(), req.SessionID)
	if err != nil {
		writeError(w, pkg.NewInternal(err))
		return
	}
	if sess == nil {
		writeError(w, pkg.NewSessionNotFound())
		return
	}

	// 设置 SSE 响应头
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		slog.Error("streaming not supported")
		return
	}

	// 运行 Agent
	eventCh, err := h.engine.Run(r.Context(), sess, req.Message)
	if err != nil {
		slog.Error("engine run failed", "error", err)
		h.writeSSEEvent(w, flusher, "error", map[string]interface{}{
			"code":    "20001",
			"message": "处理失败，请重试",
		})
		return
	}

	// 推送事件流
	for event := range eventCh {
		h.writeSSEEvent(w, flusher, event.Type, event.Data)
	}
}

// writeSSEEvent 写入 SSE 事件
func (h *ChatHandler) writeSSEEvent(w http.ResponseWriter, flusher http.Flusher, eventType string, data interface{}) {
	dataJSON, _ := json.Marshal(data)
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, string(dataJSON))
	flusher.Flush()
}
