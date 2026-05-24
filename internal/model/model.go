package model

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Responses

type ErrorInfo struct {
	Code    int   `json:"code"`
	Message error `json:"message"`
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

func Fail(c *gin.Context, status int, code int, message error) {
	c.JSON(status, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

// Modelos

type User struct {
	Matricula  int    `json:"matricula_usuario"`
	Nome       string `json:"nome_usuario"`
	Sobrenome  string `json:"sobrenome_usuario"`
	Email      string `json:"email_usuario"`
	Senha_hash string `json:"senha_hash"`
	Tema       string `json:"tema_usuario"`
}

type Materia struct {
	Codigo    int    `json:"codigo_materia"`
	Nome      string `json:"nome_materia"`
	Matricula int    `json:"matricula_usuario_materia"`
}

type Tarefa struct {
	Id     int       `json:"id_tarefa"`
	Nome   string    `json:"nome_tarefa"`
	Prazo  time.Time `json:"prazo_tarefa"`
	Codigo int       `json:"codigo_materia_tarefa"`
}
