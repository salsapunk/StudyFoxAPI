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

func (tR *TaskRepository) ListMaterias(ctx context.Context) ([]model.Materia, error) {
	rows, err := tR.pool.Query(ctx, "SELECT * FROM materia;")
	if err != nil {
		fmt.Printf("%v", err)
		return []model.Materia{}, err
	}

	var materias []model.Materia
	var materiaModel model.Materia

	for rows.Next() {
		err = rows.Scan(
			&materiaModel.Codigo,
			&materiaModel.Nome,
			&materiaModel.Matricula,
		)
		if err != nil {
			fmt.Printf("%v", err)
			return []model.Materia{}, err
		}

		materias = append(materias, materiaModel)
	}

	rows.Close()

	return materias, nil
}
