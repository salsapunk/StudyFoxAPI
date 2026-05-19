package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/log"
	"github.com/salsapunk/StudyFoxAPI/internal/handler"
	"github.com/salsapunk/StudyFoxAPI/internal/repository"
	"github.com/salsapunk/StudyFoxAPI/internal/service"
)

func main() {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Errorf("%t", err)
		return
	}
	
	taskRepository := repository.NewTaskRepo(pool)
	taskService := service.NewTaskServ(taskRepository)
	taskHandler := handler.NewTaskHand(taskService)

	router := gin.Default()
	
	router.GET("/ping", taskHandler.CheckHealth)

	router.Run(":8080")
}
