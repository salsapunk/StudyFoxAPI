package service

import (
	"context"

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

func (tS *TaskService) ListMaterias(ctx context.Context) ([]model.Materia, error) {
	materias, err := tS.repository.ListMaterias(ctx)
	if err != nil {
		return materias, err
	}

	return materias, nil
}
