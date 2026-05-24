package session

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MemoryStore 基于内存的会话存储，使用 sync.RWMutex 保证并发安全
type MemoryStore struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	maxSize  int
	ttl      time.Duration
}

// NewMemoryStore 创建内存会话存储
func NewMemoryStore(maxSize int, ttl time.Duration) *MemoryStore {
	s := &MemoryStore{
		sessions: make(map[string]*Session),
		maxSize:  maxSize,
		ttl:      ttl,
	}
	// 启动后台清理过期会话的 goroutine
	go s.cleanupLoop()
	return s
}

// Get 获取会话
func (s *MemoryStore) Get(_ context.Context, id string) (*Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess, ok := s.sessions[id]
	if !ok {
		return nil, nil
	}

	// 检查 TTL
	if s.ttl > 0 && time.Since(sess.UpdatedAt) > s.ttl {
		return nil, nil
	}

	// 返回副本，避免外部修改影响存储
	copy := *sess
	return &copy, nil
}

// Save 保存会话
func (s *MemoryStore) Save(_ context.Context, session *Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查容量限制（已存在的会话不受限制）
	if _, exists := s.sessions[session.ID]; !exists && len(s.sessions) >= s.maxSize {
		return fmt.Errorf("session store is full (max %d sessions)", s.maxSize)
	}

	session.UpdatedAt = time.Now()
	// 存储副本
	copy := *session
	s.sessions[session.ID] = &copy
	return nil
}

// Delete 删除会话
func (s *MemoryStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.sessions[id]; !ok {
		return fmt.Errorf("session %s not found", id)
	}

	delete(s.sessions, id)
	return nil
}

// List 列出所有会话 ID
func (s *MemoryStore) List(_ context.Context) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := make([]string, 0, len(s.sessions))
	for id := range s.sessions {
		ids = append(ids, id)
	}
	return ids, nil
}

// cleanupLoop 定期清理过期会话
func (s *MemoryStore) cleanupLoop() {
	if s.ttl <= 0 {
		return
	}
	ticker := time.NewTicker(s.ttl / 2)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for id, sess := range s.sessions {
			if now.Sub(sess.UpdatedAt) > s.ttl {
				delete(s.sessions, id)
			}
		}
		s.mu.Unlock()
	}
}
