package main

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/log"
	"github.com/salsapunk/StudyFoxAPI/internal/handler"
	"github.com/salsapunk/StudyFoxAPI/internal/repository"
	"github.com/salsapunk/StudyFoxAPI/internal/service"
)

func initPool() *pgxpool.Pool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Errorf("%w", err)
		return nil
	}

	config.MaxConns = 50
	config.MaxConnLifetime = time.Hour * 2
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute * 15
	
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Errorf("%t", err)
		return nil
	}

	return pool
}

func main() {
	pool := initPool()

	taskRepository := repository.NewTaskRepo(pool)
	taskService := service.NewTaskServ(taskRepository)
	taskHandler := handler.NewTaskHand(taskService)

	router := gin.Default()
	router.SetTrustedProxies([]string{"192.186.0.10"})

	// CREATE
	router.GET("/api/v1/validate", taskHandler.RequireAuth, taskHandler.Validate)
	router.POST("/api/v1/singup", taskHandler.SingUp)
	router.POST("/api/v1/login", taskHandler.Login)
	router.POST("/api/v1/materia", taskHandler.RequireAuth, taskHandler.CriarMateria)
	router.POST("/api/v1/materia/:codigo/tarefa", taskHandler.RequireAuth, taskHandler.CriarTarefa)

	// READ
	router.GET("/api/v1/ping", taskHandler.RequireAuth, taskHandler.CheckHealth)
	router.GET("/api/v1/materias", taskHandler.RequireAuth, taskHandler.ListarMaterias)
	router.GET("/api/v1/materia/:codigo", taskHandler.RequireAuth, taskHandler.LerMateria)
	router.GET("/api/v1/materia/:codigo/tarefas", taskHandler.RequireAuth, taskHandler.ListarTarefas)
	router.GET("/api/v1/materia/:codigo/tarefa/:id", taskHandler.RequireAuth, taskHandler.LerTarefa)

	// UPDATE
	router.PUT("/api/v1/email", taskHandler.RequireAuth, taskHandler.MudarEmail)
	router.PUT("/api/v1/senha", taskHandler.RequireAuth, taskHandler.MudarSenha)
	router.PUT("/api/v1/tema", taskHandler.RequireAuth, taskHandler.MudarTema)
	router.PUT("/api/v1/materia/:codigo", taskHandler.RequireAuth, taskHandler.MudarNomeMateria)
	router.PUT("/api/v1/materia/:codigo/tarefa/:id/nome", taskHandler.RequireAuth, taskHandler.MudarNomeTarefa)
	router.PUT("/api/v1/materia/:codigo/tarefa/:id/prazo", taskHandler.RequireAuth, taskHandler.MudarPrazoTarefa)
	router.PUT("/api/v1/materia/:codigo/tarefa/:id/anotacao", taskHandler.RequireAuth, taskHandler.MudarAnotacaoTarefa)
	router.PUT("/api/v1/materia/:codigo/tarefa/:id/status", taskHandler.RequireAuth, taskHandler.MudarStatusTarefa)

	// DELETE
	router.DELETE("/api/v1/usuario", taskHandler.RequireAuth, taskHandler.DeletarUsuario)
	router.DELETE("/api/v1/materia/:codigo", taskHandler.RequireAuth, taskHandler.DeletarMateria)
	router.DELETE("/api/v1/materia/:codigo/tarefa/:id", taskHandler.RequireAuth, taskHandler.DeletarTarefa)

	router.Run(":33249")
}
