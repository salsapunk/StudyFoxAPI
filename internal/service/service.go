package service

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/salsapunk/StudyFoxAPI/internal/model"
	"github.com/salsapunk/StudyFoxAPI/internal/repository"
)

type TaskService struct {
	repository *repository.TaskRepository
}

func NewTaskServ(TaskRepo *repository.TaskRepository) *TaskService {
	return &TaskService{
		repository: TaskRepo,
	}
}

func (tS *TaskService) CriarUsuario(ctx context.Context, usuario *model.Usuario) (int, error) {
	validate := validator.New()

	err := validate.Struct(usuario)
	if err != nil {
		return 0, err
	}

	matricula, err := tS.CriarUsuario(ctx, usuario)
	if err != nil {
		return 0, err
	}

	return matricula, nil
}

func (tS *TaskService) CriarMateria(ctx context.Context, materia *model.Materia) (int, error) {
	validate := validator.New()

	err := validate.Struct(materia)
	if err != nil {
		return 0, err
	}

	codigo, err := tS.CriarMateria(ctx, materia)
	if err != nil {
		return 0, err
	}

	return codigo, nil
}

func (tS *TaskService) LerUsuario(ctx context.Context) (model.Usuario, error) {
	usuario, err := tS.repository.LerUsuario(ctx)
	if err != nil {
		return usuario, err
	}

	return usuario, nil
}

func (tS *TaskService) ListarMaterias(ctx context.Context) ([]model.Materia, error) {
	materias, err := tS.repository.ListarMaterias(ctx)
	if err != nil {
		return materias, err
	}

	return materias, nil
}

func (tS *TaskService) ListarTarefas(ctx context.Context) ([]model.Tarefa, error) {
	tarefas, err := tS.repository.ListarTarefas(ctx)
	if err != nil {
		return tarefas, err
	}

	return tarefas, nil
}

func (tS *TaskService) LerTarefa(ctx context.Context) (model.Tarefa, error) {
	tarefa, err := tS.repository.LerTarefa(ctx)
	if err != nil {
		return tarefa, err
	}

	return tarefa, nil
}
