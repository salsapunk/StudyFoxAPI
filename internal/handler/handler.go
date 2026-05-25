package handler

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

// GET

func (tH *TaskHandler) LerUsuario(c *gin.Context) {
	ctx := c.Request.Context()

	usuario, err := tH.service.LerUsuario(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, model.Response{
			Success: false,
			Data:    usuario,
			Error: &model.ErrorInfo{
				Code:    http.StatusBadGateway,
				Message: err,
			},
		})
		return
	}

	c.JSON(http.StatusAccepted, model.Response{
		Success: true,
		Data:    usuario,
	})
}

func (tH *TaskHandler) ListarMaterias(c *gin.Context) {
	ctx := c.Request.Context()

	materias, err := tH.service.ListarMaterias(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, model.Response{
			Success: false,
			Data:    materias,
			Error: &model.ErrorInfo{
				Code:    http.StatusBadGateway,
				Message: err,
			},
		})
		return
	}

	c.JSON(http.StatusAccepted, model.Response{
		Success: true,
		Data:    materias,
	})
}

func (tH *TaskHandler) ListarTarefas(c *gin.Context) {
	ctx := c.Request.Context()

	tarefas, err := tH.service.ListarTarefas(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, model.Response{
			Success: false,
			Data:    tarefas,
			Error: &model.ErrorInfo{
				Code:    http.StatusBadGateway,
				Message: err,
			},
		})
		return
	}

	c.JSON(http.StatusAccepted, model.Response{
		Success: true,
		Data:    tarefas,
	})
}

func (tH *TaskHandler) LerTarefa(c *gin.Context) {
	ctx := c.Request.Context()

	tarefa, err := tH.service.LerTarefa(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, model.Response{
			Success: false,
			Data:    tarefa,
			Error: &model.ErrorInfo{
				Code:    http.StatusBadGateway,
				Message: err,
			},
		})
		return
	}

}
