package pkg

import "fmt"

// AppError 应用错误，包含错误码和用户友好的消息
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code=%d msg=%s err=%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code=%d msg=%s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// 预定义错误码
const (
	ErrCodeSuccess         = 0
	ErrCodeBadRequest      = 10001
	ErrCodeSessionNotFound = 10002
	ErrCodeMessageInvalid  = 10003
	ErrCodeInternal        = 20001
	ErrCodeLLMUnavailable  = 20002
)

// NewBadRequest 请求参数校验失败
func NewBadRequest(msg string) *AppError {
	return &AppError{Code: ErrCodeBadRequest, Message: msg}
}

// NewSessionNotFound 会话不存在
func NewSessionNotFound() *AppError {
	return &AppError{Code: ErrCodeSessionNotFound, Message: "session not found"}
}

// NewMessageInvalid 消息内容无效
func NewMessageInvalid(msg string) *AppError {
	return &AppError{Code: ErrCodeMessageInvalid, Message: msg}
}

// NewInternal 服务器内部错误
func NewInternal(err error) *AppError {
	return &AppError{Code: ErrCodeInternal, Message: "internal server error", Err: err}
}

// NewLLMUnavailable LLM 服务不可用
func NewLLMUnavailable(err error) *AppError {
	return &AppError{Code: ErrCodeLLMUnavailable, Message: "AI 服务暂时不可用，请稍后重试", Err: err}
}
