package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepo(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		pool: pool,
	}
}
