package model

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CRIAR_USUARIO   = "INSERT INTO usuario(matricula, nome, sobrenome, email, senha_hash, tema) VALUES($1, $2, $3, $4, $5, $6) RETURNING matricula;"
	CRIAR_MATERIA   = "INSERT INTO materia(codigo, nome, matricula) VALUES($1, $2, $3) RETURNING codigo;"
	LER_USUARIO     = "SELECT matricula, email, senha_hash, tema FROM usuario;"
	LISTAR_MATERIAS = "SELECT * FROM materia;"
	LISTAR_TAREFAS  = "SELECT * FROM tarefa;"
	LER_TAREFA      = "SELECT matricula, email, senha_hash, tema FROM tarefa;"
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

type Usuario struct {
	Matricula  int    `json:"matricula_usuario" validate:"required"`
	Nome       string `json:"nome_usuario" validate:"required"`
	Sobrenome  string `json:"sobrenome_usuario"`
	Email      string `json:"email_usuario" validate:"required"`
	Senha_hash string `json:"senha_hash" validate:"required"`
	Tema       string `json:"tema_usuario" validate:"required"`
}

type Materia struct {
	Codigo    int    `json:"codigo_materia" validate:"required"`
	Nome      string `json:"nome_materia" validate:"required"`
	Matricula int    `json:"matricula_usuario_materia" validate:"required"`
}

type Tarefa struct {
	Id     int    `json:"id_tarefa" validate:"required"`
	Nome   string `json:"nome_tarefa" validate:"required"`
	Prazo  string `json:"prazo_tarefa" time_format:"02-01-2006"`
	Codigo int    `json:"codigo_materia_tarefa" validate:"required"`
}
