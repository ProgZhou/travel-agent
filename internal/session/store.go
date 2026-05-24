package session

import "context"

// Store 会话存储接口，支持未来替换为 SQLite/Redis
type Store interface {
	// Get 获取会话，不存在时返回 nil, nil
	Get(ctx context.Context, id string) (*Session, error)
	// Save 保存会话（新建或更新）
	Save(ctx context.Context, session *Session) error
	// Delete 删除会话
	Delete(ctx context.Context, id string) error
	// List 列出所有会话 ID（用于管理）
	List(ctx context.Context) ([]string, error)
}
