package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/salsapunk/StudyFoxAPI/internal/model"
)

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepo(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		pool: pool,
	}
}

// GET

func (tR *TaskRepository) LerUsuario(ctx context.Context) (model.Usuario, error) {
	row := tR.pool.QueryRow(ctx, "SELECT matricula, email, senha_hash, tema FROM usuario;")

	var usuario model.Usuario

	err := row.Scan(
		&usuario.Matricula,
		&usuario.Email,
		&usuario.Senha_hash,
		&usuario.Tema,
	)
	if err != nil {
		return model.Usuario{}, err
	}

	return usuario, nil
}

func (tR *TaskRepository) ListarMaterias(ctx context.Context) ([]model.Materia, error) {
	rows, err := tR.pool.Query(ctx, "SELECT * FROM materia;")
	if err != nil {
		fmt.Printf("%v", err)
		return []model.Materia{}, err
	}

	var materias []model.Materia
	var materia model.Materia

	for rows.Next() {
		err = rows.Scan(
			&materia.Codigo,
			&materia.Nome,
			&materia.Matricula,
		)
		if err != nil {
			fmt.Printf("%v", err)
			return []model.Materia{}, err
		}

		materias = append(materias, materia)
	}

	rows.Close()

	return materias, nil
}

func (tR *TaskRepository) ListarTarefas(ctx context.Context) ([]model.Tarefa, error) {
	rows, err := tR.pool.Query(ctx, "SELECT * FROM tarefa;")
	if err != nil {
		fmt.Printf("%v", err)
		return []model.Tarefa{}, err
	}

	var tarefas []model.Tarefa
	var tarefa model.Tarefa

	for rows.Next() {
		err = rows.Scan(
			&tarefa.Id,
			&tarefa.Nome,
			&tarefa.Prazo,
			&tarefa.Codigo,
		)
		if err != nil {
			fmt.Printf("%v", err)
			return []model.Tarefa{}, err
		}

		tarefas = append(tarefas, tarefa)
	}

	rows.Close()

	return tarefas, nil
}

func (tR *TaskRepository) LerTarefa(ctx context.Context) (model.Tarefa, error) {
	row := tR.pool.QueryRow(ctx, "SELECT matricula, email, senha_hash, tema FROM tarefa;")

	var tarefa model.Tarefa

	err := row.Scan(
		&tarefa.Id,
		&tarefa.Nome,
		&tarefa.Prazo,
		&tarefa.Codigo,
	)
	if err != nil {
		return model.Tarefa{}, err
	}

	return tarefa, nil
}
