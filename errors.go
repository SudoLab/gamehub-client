package gamehub

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("GameHub API Error [%d]: %s", e.Code, e.Message)
}

var (
	ErrUnauthorized      = &Error{Code: http.StatusUnauthorized, Message: "Unauthorized", Type: "auth_error"}
	ErrInsufficientCoins = &Error{Code: http.StatusBadRequest, Message: "Insufficient coins", Type: "insufficient_funds"}
	ErrUserNotFound      = &Error{Code: http.StatusNotFound, Message: "User not found", Type: "user_not_found"}
	ErrInvalidAPIKey     = &Error{Code: http.StatusUnauthorized, Message: "Invalid API key", Type: "invalid_api_key"}
	ErrRateLimited       = &Error{Code: http.StatusTooManyRequests, Message: "Rate limit exceeded", Type: "rate_limited"}
	ErrServerError       = &Error{Code: http.StatusInternalServerError, Message: "Internal server error", Type: "server_error"}
)

func NewError(code int, message, errorType string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Type:    errorType,
	}
}
