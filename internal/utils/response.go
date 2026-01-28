package utils

import "github.com/gin-gonic/gin"

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

func Success(c *gin.Context, data interface{}, message string) {
	c.JSON(200, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func Created(c *gin.Context, data interface{}, message string) {
	c.JSON(201, SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func Error(c *gin.Context, status int, code string, message string) {
	c.JSON(status, ErrorResponse{
		Success: false,
		Code:    code,
		Message: message,
	})
}
