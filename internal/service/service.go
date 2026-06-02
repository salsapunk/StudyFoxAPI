package service

import (
	"context"
	"errors"
	"fmt"

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

	matricula, err := tS.repository.CriarUsuario(ctx, usuario)
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

	codigo, err := tS.repository.CriarMateria(ctx, materia)
	if err != nil {
		return 0, err
	}

	return codigo, nil
}

func (tS *TaskService) CriarTarefa(ctx context.Context, tarefa *model.Tarefa) (int, error) {
	validate := validator.New()

	err := validate.Struct(tarefa)
	if err != nil {
		return 0, err
	}

	id, err := tS.repository.CriarTarefa(ctx, tarefa)
	if err != nil {
		return 0, nil
	}

	return id, nil
}

// READ

func (tS *TaskService) LerUsuario(ctx context.Context, matricula int) (model.Usuario, error) {
	usuario, err := tS.repository.LerUsuario(ctx, matricula)
	if err != nil {
		return usuario, err
	}

	return usuario, nil
}

func (tS *TaskService) ListarMaterias(ctx context.Context, matricula int) ([]model.Materia, error) {
	materias, err := tS.repository.ListarMaterias(ctx, matricula)
	if err != nil {
		return materias, err
	}

	return materias, nil
}

func (tS *TaskService) LerMateria(ctx context.Context, matricula int, codigo int) (model.Materia, error) {
	materia, err := tS.repository.LerMateria(ctx, matricula, codigo)
	if err != nil {
		return materia, err
	}

	return materia, nil
}

func (tS *TaskService) ListarTarefas(ctx context.Context, codigo int) ([]model.Tarefa, error) {
	tarefas, err := tS.repository.ListarTarefas(ctx, codigo)
	if err != nil {
		return tarefas, err
	}

	return tarefas, nil
}

func (tS *TaskService) LerTarefa(ctx context.Context, codigo int, id int) (model.Tarefa, error) {
	tarefa, err := tS.repository.LerTarefa(ctx, codigo, id)
	if err != nil {
		return tarefa, err
	}

	return tarefa, nil
}

// UPDATE

// ...

// DELETE

func (tS *TaskService) DeletarUsuario(ctx context.Context, matricula int) error {
	_, err := tS.repository.LerUsuario(ctx, matricula)
	if err != nil {
		if errors.Is(err, model.ErrUsuarioNotFound) {
			return model.ErrUsuarioNotFound
		}
		return fmt.Errorf("erro ao buscar usuário %d: %w", matricula, err)
	}

	err = tS.repository.DeletarUsuario(ctx, matricula)
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário %d: %w", matricula, err)
	}

	return nil
}

func (tS *TaskService) DeletarMateria(ctx context.Context, matricula int, codigo int) error {
	_, err := tS.repository.LerMateria(ctx, matricula, codigo)
	if err != nil {
		if errors.Is(err, model.ErrMateriaNotFound) {
			return model.ErrMateriaNotFound
		}
		return fmt.Errorf("erro ao buscar matéria %d: %w", codigo, err)
	}

	err = tS.repository.DeletarMateria(ctx, matricula, codigo)
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário %d: %w", codigo, err)
	}

	return nil
}

func (tS *TaskService) DeletarTarefa(ctx context.Context, codigo int, id int) error {
	_, err := tS.repository.LerTarefa(ctx, codigo, id)
	if err != nil {
		if errors.Is(err, model.ErrTarefaNotFound) {
			return model.ErrTarefaNotFound
		}
		return fmt.Errorf("erro ao buscar matéria %d: %w", id, err)
	}

	err = tS.repository.DeletarTarefa(ctx, codigo, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário %d: %w", id, err)
	}

	return nil
}
