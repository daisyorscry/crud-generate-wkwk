package responses

import (
	"daisy/pkg/pkgErr"

	"github.com/gofiber/fiber/v2"
)

type ErrorDetailWrapper struct {
	Detail string `json:"detail"`
}

type BaseResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    any                 `json:"data,omitempty"`
	Error   *ErrorDetailWrapper `json:"error,omitempty"`
	Errors  pkgErr.FieldErrors  `json:"errors,omitempty"`
	Total   int64               `json:"total,omitempty"` // ← tambahkan ini
}

// --- Handler-level Success ---
func Success(c *fiber.Ctx, code int, message string, data any) error {
	return c.Status(code).JSON(BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// --- Handler-level Error ---
func Error(c *fiber.Ctx, code int, message string, err error) error {
	return c.Status(code).JSON(BaseResponse{
		Code:    code,
		Message: message,
		Error:   wrapError(err),
	})
}

// --- Handler-level Validation ---
func ValidationError(c *fiber.Ctx, code int, message string, errs map[string]string) error {
	return c.Status(code).JSON(BaseResponse{
		Code:    code,
		Message: message,
		Errors:  errs,
	})
}

// --- Handler-level Pagination ---
func Paginate(c *fiber.Ctx, code int, message string, data any, total int64) error {
	return c.Status(code).JSON(BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
		Total:   total,
	})
}

// --- Service-level Success ---
func ResponseSuccess(code int, message string, data any) *BaseResponse {
	return &BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// --- Service-level Error ---
func ResponseError(code int, message string, err error) *BaseResponse {
	return &BaseResponse{
		Code:    code,
		Message: message,
		Error:   wrapError(err),
	}
}

// --- Service-level Validation ---
func ResponseValidation(code int, message string, errs pkgErr.FieldErrors) *BaseResponse {
	return &BaseResponse{
		Code:    code,
		Message: message,
		Errors:  errs,
	}
}

// --- Service-level Pagination ---
func ResponsePagination(code int, message string, data any, total int64) *BaseResponse {
	return &BaseResponse{
		Code:    code,
		Message: message,
		Data:    data,
		Total:   total,
	}
}

// --- Private helper to wrap error ---
func wrapError(err error) *ErrorDetailWrapper {
	if err == nil {
		return nil
	}
	return &ErrorDetailWrapper{
		Detail: err.Error(),
	}
}
