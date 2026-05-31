package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/salsapunk/StudyFoxAPI/internal/handler"
	"github.com/salsapunk/StudyFoxAPI/internal/repository"
	"github.com/salsapunk/StudyFoxAPI/internal/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Errorf("%t", err)
		return
	}

	taskRepository := repository.NewTaskRepo(pool)
	taskService := service.NewTaskServ(taskRepository)
	taskHandler := handler.NewTaskHand(taskService)

	router := gin.Default()
	router.SetTrustedProxies([]string{"192.186.0.10"})

	// CREATE
	router.POST("/api/v1/usuario", taskHandler.CriarUsuario)
	router.POST("/api/v1/usuario/:matricula/materia", taskHandler.CriarMateria)
	router.POST("/api/v1/materia/:codigo/tarefa", taskHandler.CriarTarefa)

	// READ
	router.GET("/api/v1/ping", taskHandler.CheckHealth)
	router.GET("/api/v1/usuario/:matricula", taskHandler.LerUsuario)
	router.GET("/api/v1/usuario/:matricula/materias", taskHandler.ListarMaterias)
	router.GET("/api/v1/materia/:codigo/tarefas", taskHandler.ListarTarefas)
	router.GET("/api/v1/materia/:codigo/tarefa/:id", taskHandler.LerTarefa)

	// UPDATE
	// ?

	// DELETE
	// router.DELETE("/api/v1/usuario/:matricula", taskHandler.DeletarUsuario)
	// router.DELETE("/api/v1/usuario/:matricula/materia/:codigo", taskHandler.DeletarMateria)
	// router.DELETE("/api/v1/usuario/:matricula/materia/:codigo/tarefa/:id", taskHandler.DeletarTarefa)

	router.Run(":8080")
}
