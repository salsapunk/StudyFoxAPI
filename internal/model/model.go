package model

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

type Response struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
	Meta    *Meta      `json:"meta,omitempty"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func Fail(c *gin.Context, status int, code, message string) {
	c.JSON(status, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}
