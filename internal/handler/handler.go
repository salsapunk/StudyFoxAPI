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

func (tH *TaskHandler) ListarMaterias(c *gin.Context) {
	ctx := c.Request.Context()

	materias, err := tH.service.ListMaterias(ctx)
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
