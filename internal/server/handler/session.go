package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"travel-agent/internal/pkg"
	"travel-agent/internal/session"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// SessionHandler 会话处理器
type SessionHandler struct {
	store session.Store
}

// NewSessionHandler 创建会话处理器
func NewSessionHandler(store session.Store) *SessionHandler {
	return &SessionHandler{store: store}
}

// Create 创建会话
func (h *SessionHandler) Create(w http.ResponseWriter, r *http.Request) {
	sess := &session.Session{
		ID:           uuid.New().String(),
		Phase:        session.PhaseCollecting,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ChatHistory:  []session.ChatMessage{},
		BookingItems: []session.BookingItem{},
	}

	if err := h.store.Save(r.Context(), sess); err != nil {
		writeError(w, pkg.NewInternal(err))
		return
	}

	writeSuccess(w, map[string]interface{}{
		"session_id": sess.ID,
		"phase":      sess.Phase,
		"created_at": sess.CreatedAt.Format(time.RFC3339),
	})
}

// Get 获取会话
func (h *SessionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, pkg.NewBadRequest("session_id is required"))
		return
	}

	sess, err := h.store.Get(r.Context(), id)
	if err != nil {
		writeError(w, pkg.NewInternal(err))
		return
	}

	if sess == nil {
		writeError(w, pkg.NewSessionNotFound())
		return
	}

	writeSuccess(w, sess)
}

// Delete 删除会话
func (h *SessionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, pkg.NewBadRequest("session_id is required"))
		return
	}

	if err := h.store.Delete(r.Context(), id); err != nil {
		writeError(w, pkg.NewInternal(err))
		return
	}

	writeSuccess(w, map[string]string{"message": "session deleted"})
}

// writeSuccess 写入成功响应
func writeSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

// writeError 写入错误响应
func writeError(w http.ResponseWriter, err *pkg.AppError) {
	w.Header().Set("Content-Type", "application/json")
	statusCode := http.StatusInternalServerError
	if err.Code >= 10000 && err.Code < 20000 {
		statusCode = http.StatusBadRequest
	}
	if err.Code == pkg.ErrCodeSessionNotFound {
		statusCode = http.StatusNotFound
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    err.Code,
		"message": err.Message,
		"data":    nil,
	})
}
