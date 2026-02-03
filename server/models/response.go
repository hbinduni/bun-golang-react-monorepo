package models

// ApiResponse is the standard API response wrapper
type ApiResponse[T any] struct {
	Success bool   `json:"success"`
	Data    *T     `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse creates a successful API response
func SuccessResponse[T any](data T) ApiResponse[T] {
	return ApiResponse[T]{
		Success: true,
		Data:    &data,
	}
}

// ErrorResponse creates an error API response
func ErrorResponse(message string) ApiResponse[any] {
	return ApiResponse[any]{
		Success: false,
		Error:   message,
	}
}
