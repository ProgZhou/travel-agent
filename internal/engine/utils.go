package engine

import (
	"time"

	"github.com/google/uuid"
)

// generateID 生成唯一 ID
func generateID() string {
	return uuid.New().String()
}

// nowMillis 获取当前时间戳（毫秒）
func nowMillis() int64 {
	return time.Now().UnixMilli()
}
