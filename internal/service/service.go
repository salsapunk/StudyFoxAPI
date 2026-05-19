package service

import "github.com/salsapunk/StudyFoxAPI/internal/repository"

type TaskService struct {
	repository *repository.TaskRepository
}

func NewTaskServ(TaskRepo *repository.TaskRepository) *TaskService {
	return &TaskService{
		repository: TaskRepo,
	}
}
