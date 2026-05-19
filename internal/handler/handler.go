package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func CheckHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
