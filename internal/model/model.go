package model

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CRIAR_USUARIO = "INSERT INTO usuario(email, senha_hash) VALUES($1, $2) RETURNING matricula;"
	CRIAR_MATERIA = "INSERT INTO materia(nome, matricula) VALUES($1, $2) RETURNING codigo;"
	CRIAR_TAREFA  = "INSERT INTO tarefa(nome, anotacao, prazo, codigo) VALUES($1, $2, $3, $4) RETURNING id;"

	LER_USUARIO     = "SELECT matricula, email, senha_hash, tema FROM usuario WHERE matricula = $1;"
	LISTAR_MATERIAS = "SELECT codigo, nome, m.matricula FROM materia m INNER JOIN usuario u ON m.matricula = u.matricula AND u.matricula = $1;" //materias de um usuario
	LER_MATERIA     = "SELECT codigo, m.nome, u.matricula FROM materia m INNER JOIN usuario u ON u.matricula = $1 AND codigo = $2;"
	LISTAR_TAREFAS  = "SELECT id, t.nome, prazo, anotacao, t.codigo FROM tarefa t INNER JOIN materia m ON t.codigo = m.codigo AND m.codigo = $1;" // tarefas em uma matéria
	LER_TAREFA      = "SELECT id, t.nome, anotacao, prazo, t.codigo FROM tarefa t INNER JOIN materia m ON m.codigo = $1 AND id = $2;"

	MUDAR_SENHA_USR = "UPDATE TABLE usuario SET senha_hash = $1 WHERE matricula = $2"
	MUDAR_EMAIL_USR = "UPDATE TABLE usuario SET email = $1 WHERE matricula = $2"
	MUDAR_NOME_MAT  = "UPDATE TABLE materia SET nome = $1 WHERE matricula = $2"
	MUDAR_NOME_TAR  = "UPDATE TABLE tarefa SET nome = $1 WHERE codigo = $2"
	MUDAR_PRAZ_TAR  = "UPDATE TABLE tarefa SET prazo = $1 WHERE codigo = $2"
	MUDAR_ANOT_TAR  = "UPDATE TABLE tarefa SET anotacao = $1 WHERE codigo = $2"

	DELETAR_USUARIO = "DELETE FROM usuario WHERE matricula = $1"
	DELETAR_MATERIA = "DELETE FROM materia WHERE matricula = $1 AND codigo = $2"
	DELETAR_TAREFA  = "DELETE FROM tarefa WHERE codigo = $1 AND id = $2"
)

func ErrNotFound(fragmento string) error {
	err := fmt.Sprintf("%s não encontrado", fragmento)
	Err := errors.New(err)
	return Err
}

var (
	ErrUsuarioNotFound = ErrNotFound("Usuário")
	ErrMateriaNotFound = ErrNotFound("Matéria")
	ErrTarefaNotFound  = ErrNotFound("Tarefa")
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
	Matricula  int    `json:"matricula_usuario"`
	Email      string `json:"email" validate:"required"`
	Senha_hash string `json:"senha" validate:"required"`
	Tema       string `json:"tema"`
}

type Materia struct {
	Codigo    int    `json:"codigo_materia"`
	Nome      string `json:"nome" validate:"required"`
	Matricula int    `json:"matricula_usuario" validate:"required"`
}

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

type Tarefa struct {
	Id       int    `json:"id_tarefa"`
	Nome     string `json:"nome" validate:"required"`
	Prazo    Date   `json:"prazo"`
	Anotacao string `json:"anotacao"`
	Codigo   int    `json:"codigo_materia" validate:"required"`
}
