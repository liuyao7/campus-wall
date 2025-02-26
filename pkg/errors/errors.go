package errors

import "fmt"

type AppError struct {
    Code    int
    Message string
    Err     error
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

// 定义错误码
const (
    ErrCodeBadRequest          = 400
    ErrCodeUnauthorized        = 401
    ErrCodeForbidden          = 403
    ErrCodeNotFound           = 404
    ErrCodeInternalServer     = 500
)

// 创建错误的便捷方法
func NewBadRequestError(message string) *AppError {
    return &AppError{
        Code:    ErrCodeBadRequest,
        Message: message,
    }
}

func NewUnauthorizedError(message string) *AppError {
    return &AppError{
        Code:    ErrCodeUnauthorized,
        Message: message,
    }
}

func NewForbiddenError(message string) *AppError {
    return &AppError{
        Code:    ErrCodeForbidden,
        Message: message,
    }
}

func NewNotFoundError(message string) *AppError {
    return &AppError{
        Code:    ErrCodeNotFound,
        Message: message,
    }
}

func NewInternalServerError(message string, err error) *AppError {
    return &AppError{
        Code:    ErrCodeInternalServer,
        Message: message,
        Err:     err,
    }
}