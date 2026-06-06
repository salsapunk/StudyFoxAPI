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
	router.GET("/api/v1/validate", taskHandler.RequireAuth, taskHandler.Validate)
	router.POST("/api/v1/singup", taskHandler.SingUp)
	router.POST("/api/v1/login", taskHandler.Login)
	router.POST("/api/v1/usuario/:matricula/materia", taskHandler.RequireAuth, taskHandler.CriarMateria)
	router.POST("/api/v1/materia/:codigo/tarefa", taskHandler.RequireAuth, taskHandler.CriarTarefa)

	// READ
	router.GET("/api/v1/ping", taskHandler.RequireAuth, taskHandler.CheckHealth)
	// router.GET("/api/v1/usuario/:matricula", taskHandler.LerUsuario)
	router.GET("/api/v1/usuario/:matricula/materias", taskHandler.RequireAuth, taskHandler.ListarMaterias)
	router.GET("/api/v1/usuario/:matricula/materia/:codigo", taskHandler.RequireAuth, taskHandler.LerMateria)
	router.GET("/api/v1/materia/:codigo/tarefas", taskHandler.RequireAuth, taskHandler.ListarTarefas)
	router.GET("/api/v1/materia/:codigo/tarefa/:id", taskHandler.RequireAuth, taskHandler.LerTarefa)

	// UPDATE
	router.PUT("/api/v1/usuario/:matricula/email", taskHandler.RequireAuth, taskHandler.MudarEmail)
	router.PUT("/api/v1/usuario/:matricula/senha", taskHandler.RequireAuth, taskHandler.MudarSenha)
	router.PUT("/api/v1/usuario/:matricula/tema", taskHandler.RequireAuth, taskHandler.MudarTema)
	router.PUT("/api/v1/usuario/:matricula/materia/:codigo", taskHandler.RequireAuth, taskHandler.MudarNomeMateria)
	router.PUT("/api/v1/materia/:codigo/tarefa/:id/nome", taskHandler.RequireAuth, taskHandler.MudarNomeTarefa)
	router.PUT("/api/v1/materia/:codigo/tarefa/:id/prazo", taskHandler.RequireAuth, taskHandler.MudarPrazoTarefa)
	router.PUT("/api/v1/materia/:codigo/tarefa/:id/anotacao", taskHandler.RequireAuth, taskHandler.MudarAnotacaoTarefa)
	router.PUT("/api/v1/materia/:codigo/tarefa/:id/status", taskHandler.RequireAuth, taskHandler.MudarStatusTarefa)

	// DELETE
	router.DELETE("/api/v1/usuario/:matricula", taskHandler.RequireAuth, taskHandler.DeletarUsuario)
	router.DELETE("/api/v1/usuario/:matricula/materia/:codigo", taskHandler.RequireAuth, taskHandler.DeletarMateria)
	router.DELETE("/api/v1/materia/:codigo/tarefa/:id", taskHandler.RequireAuth, taskHandler.DeletarTarefa)

	router.Run(":8080")
}
