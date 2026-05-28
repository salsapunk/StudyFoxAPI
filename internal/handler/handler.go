package handler

// arrumar os códigos dos erros

import (
	"net/http"

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
		model.Fail(c, http.StatusBadRequest, 400, err)
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

	var materia model.Materia
	err := c.ShouldBindJSON(&materia)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	codigo, err := tH.service.CriarMateria(ctx, &materia)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	model.OK(c, codigo)
}

func (tH *TaskHandler) LerUsuario(c *gin.Context) {
	ctx := c.Request.Context()

	usuario, err := tH.service.LerUsuario(ctx)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	model.OK(c, usuario)
}

func (tH *TaskHandler) ListarMaterias(c *gin.Context) {
	ctx := c.Request.Context()

	materias, err := tH.service.ListarMaterias(ctx)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	model.OK(c, materias)
}

func (tH *TaskHandler) ListarTarefas(c *gin.Context) {
	ctx := c.Request.Context()

	tarefas, err := tH.service.ListarTarefas(ctx)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	model.OK(c, tarefas)
}

func (tH *TaskHandler) LerTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	tarefa, err := tH.service.LerTarefa(ctx)
	if err != nil {
		model.Fail(c, http.StatusBadGateway, 400, err)
		return
	}

	model.OK(c, tarefa)
}
