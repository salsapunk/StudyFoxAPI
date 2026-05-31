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

// CREATE

func (tR *TaskRepository) CriarUsuario(ctx context.Context, usuario *model.Usuario) (int, error) {
	row := tR.pool.QueryRow(ctx, model.CRIAR_USUARIO,
		&usuario.Email,
		&usuario.Senha_hash,
	)

	var matricula int

	err := row.Scan(&matricula)
	if err != nil {
		return 0, err
	}

	return matricula, nil
}

func (tR *TaskRepository) CriarMateria(ctx context.Context, materia *model.Materia) (int, error) {
	row := tR.pool.QueryRow(ctx, model.CRIAR_MATERIA,
		&materia.Nome,
		&materia.Matricula,
	)

	var codigo int

	err := row.Scan(&codigo)
	if err != nil {
		fmt.Println("scan")
		return 0, err
	}

	return codigo, nil
}

func (tR *TaskRepository) CriarTarefa(ctx context.Context, tarefa *model.Tarefa) (int, error) {
	row := tR.pool.QueryRow(ctx, model.CRIAR_TAREFA,
		&tarefa.Nome,
		&tarefa.Anotacao,
		&tarefa.Prazo.Time,
		&tarefa.Codigo,
	)

	var id int

	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// READ

func (tR *TaskRepository) LerUsuario(ctx context.Context, matricula int) (model.Usuario, error) {
	row := tR.pool.QueryRow(ctx, model.LER_USUARIO, matricula)

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

func (tR *TaskRepository) ListarMaterias(ctx context.Context, matricula int) ([]model.Materia, error) {
	rows, err := tR.pool.Query(ctx, model.LISTAR_MATERIAS, matricula)
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

func (tR *TaskRepository) ListarTarefas(ctx context.Context, codigo int) ([]model.Tarefa, error) {
	rows, err := tR.pool.Query(ctx, model.LISTAR_TAREFAS, codigo)
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
			&tarefa.Anotacao,
			&tarefa.Prazo.Time,
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

func (tR *TaskRepository) LerTarefa(ctx context.Context, id int) (model.Tarefa, error) {
	row := tR.pool.QueryRow(ctx, model.LER_TAREFA, id)

	var tarefa model.Tarefa

	err := row.Scan(
		&tarefa.Id,
		&tarefa.Nome,
		&tarefa.Anotacao,
		&tarefa.Prazo.Time,
		&tarefa.Codigo,
	)
	if err != nil {
		return model.Tarefa{}, err
	}

	return tarefa, nil
}
