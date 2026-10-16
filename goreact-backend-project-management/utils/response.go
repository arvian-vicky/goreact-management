package utils

import "github.com/gofiber/fiber/v2"

type Response struct {
	Status       string      `json:"status"` // "success" or "error"
	ResponseCode int         `json:"response_code"`
	Message      string      `json:"message,omitempty"` // Optional
	Data         interface{} `json:"data,omitempty"`    // Bisa user, list, dll
	Error        string      `json:"error,omitempty"`   // Kalau error
}

type PaginationMeta struct {
	Page      int    `json:"page" example:"1"`
	Limit     int    `json:"limit" example:"10"`
	Total     int    `json:"total" example:"100"`
	TotalPage int    `json:"total_pages" example:"10"`
	Filter    string `json:"filter,omitempty" example:"nama"`
	Sort      string `json:"sort,omitempty" example:"email desc"`
}

type PaginatedResponse struct {
	Status       string         `json:"status"`
	ResponseCode int            `json:"response_code"`
	Message      string         `json:"message"`
	Data         interface{}    `json:"data"`
	Meta         PaginationMeta `json:"meta"`
}

// Success Response
func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Status:       "success",
		ResponseCode: fiber.StatusOK,
		Message:      message,
		Data:         data,
	})
}

// Created Response (201)
func Created(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Status:       "success",
		ResponseCode: fiber.StatusCreated,
		Message:      message,
		Data:         data,
	})
}

// BadRequest Response (400)
func BadRequest(c *fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Status:       "error",
		ResponseCode: fiber.StatusBadRequest,
		Message:      message,
		Error:        err,
	})
}

// Unauthorized Response (401)
func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Status:       "error",
		ResponseCode: fiber.StatusUnauthorized,
		Message:      message,
		Error:        "Authentication required",
	})
}

// Forbidden Response (403)
func Forbidden(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(Response{
		Status:       "error",
		ResponseCode: fiber.StatusForbidden,
		Message:      message,
		Error:        "Insufficient permissions",
	})
}

// NotFound Response (404)
func NotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Status:       "error",
		ResponseCode: fiber.StatusNotFound,
		Message:      message,
		Error:        "Resource not found",
	})
}

// Conflict Response (409)
func Conflict(c *fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusConflict).JSON(Response{
		Status:       "error",
		ResponseCode: fiber.StatusConflict,
		Message:      message,
		Error:        err,
	})
}

// ValidationError Response (422)
func ValidationError(c *fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{
		Status:       "error",
		ResponseCode: fiber.StatusUnprocessableEntity,
		Message:      message,
		Error:        err,
	})
}

// TooManyRequests Response (429)
func TooManyRequests(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(Response{
		Status:       "error",
		ResponseCode: fiber.StatusTooManyRequests,
		Message:      message,
		Error:        "Too many requests",
	})
}

// InternalServerError Response (500)
func InternalServerError(c *fiber.Ctx, message string, err string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Status:       "error",
		ResponseCode: fiber.StatusInternalServerError,
		Message:      message,
		Error:        err,
	})
}

// ServiceUnavailable Response (503)
func ServiceUnavailable(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusServiceUnavailable).JSON(Response{
		Status:       "error",
		ResponseCode: fiber.StatusServiceUnavailable,
		Message:      message,
		Error:        "Service temporarily unavailable",
	})
}

// Generic Error Response
func Error(c *fiber.Ctx, code int, message string, err string) error {
	return c.Status(code).JSON(Response{
		Status:       "error",
		ResponseCode: code,
		Message:      message,
		Error:        err,
	})
}

//paginate response

func PaginatedSuccess(c *fiber.Ctx, message string, data interface{}, meta PaginationMeta) error {
	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Status:       "success",
		ResponseCode: fiber.StatusOK,
		Message:      message,
		Data:         data,
		Meta:         meta,
	})
}
func PaginatedNotFound(c *fiber.Ctx, message string, data interface{}, meta PaginationMeta) error {
	return c.Status(fiber.StatusNotFound).JSON(PaginatedResponse{
		Status:       "not found",
		ResponseCode: fiber.StatusNotFound,
		Message:      message,
		Meta:         meta,
		Data:         data,
	})
}
