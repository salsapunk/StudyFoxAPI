package handler

// arrumar os códigos dos erros

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/salsapunk/StudyFoxAPI/internal/model"
	"github.com/salsapunk/StudyFoxAPI/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHand(TaskServ *service.TaskService) *TaskHandler {
	return &TaskHandler{
		service: TaskServ,
	}
}

func (tH *TaskHandler) CheckHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (tH *TaskHandler) CriarUsuario(c *gin.Context) {
	ctx := c.Request.Context()

	var usuario model.Usuario

	err := c.ShouldBindJSON(&usuario)
	if err != nil {
		fmt.Println("\nshouldbind")
		model.Fail(c, http.StatusBadRequest, 1, err)
		return
	}

	id, err := tH.service.CriarUsuario(ctx, &usuario)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	model.OK(c, id)
}

func (tH *TaskHandler) CriarMateria(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("matricula")
	matricula, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, http.StatusBadRequest, err)
		return
	}

	var materia model.Materia
	err = c.ShouldBindJSON(&materia)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	materia.Matricula = matricula

	codigo, err := tH.service.CriarMateria(ctx, &materia)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	model.OK(c, codigo)
}

func (tH *TaskHandler) CriarTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	var tarefa model.Tarefa
	err = c.ShouldBindJSON(&tarefa)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	tarefa.Codigo = codigo

	id, err := tH.service.CriarTarefa(ctx, &tarefa)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	model.OK(c, id)
}

// READ

func (tH *TaskHandler) LerUsuario(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("matricula")
	matricula, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, http.StatusBadRequest, err)
		return
	}

	usuario, err := tH.service.LerUsuario(ctx, matricula)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	model.OK(c, usuario)
}

func (tH *TaskHandler) ListarMaterias(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("matricula")
	matricula, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	materias, err := tH.service.ListarMaterias(ctx, matricula)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	model.OK(c, materias)
}

func (tH *TaskHandler) LerMateria(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("matricula")
	matricula, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, http.StatusBadRequest, err)
		return
	}

	usuario, err := tH.service.LerUsuario(ctx, matricula)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	model.OK(c, usuario)
}

func (tH *TaskHandler) ListarTarefas(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 399, err)
		return
	}

	tarefas, err := tH.service.ListarTarefas(ctx, codigo)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 100, err)
		return
	}

	model.OK(c, tarefas)
}

func (tH *TaskHandler) LerTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 399, err)
		return
	}

	param = c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 399, err)
		return
	}

	tarefa, err := tH.service.LerTarefa(ctx, codigo, id)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	model.OK(c, tarefa)
}

// UPDATE

// ...

// DELETE

func (tH *TaskHandler) DeletarUsuario(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("matricula")
	matricula, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 400, err)
		return
	}

	err = tH.service.DeletarUsuario(ctx, matricula)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 400, err)
		return
	}

	model.OK(c, "Usuário deletado")
}

func (tH *TaskHandler) DeletarMateria(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("matricula")
	matricula, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 400, err)
		return
	}

	param = c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 400, err)
		return
	}

	err = tH.service.DeletarMateria(ctx, matricula, codigo)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 400, err)
		return
	}

	model.OK(c, "Matéria deletada")
}

func (tH *TaskHandler) DeletarTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	param := c.Param("codigo")
	codigo, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 400, err)
		return
	}

	param = c.Param("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 400, err)
		return
	}

	err = tH.service.DeletarTarefa(ctx, codigo, id)
	if err != nil {
		model.Fail(c, http.StatusBadRequest, 400, err)
		return
	}

	model.OK(c, "Tarefa deletada")
}
