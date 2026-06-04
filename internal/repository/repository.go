package repository

import (
	"context"
	"fmt"
	"time"

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
	fmt.Println(usuario)

	row := tR.pool.QueryRow(ctx, model.CRIAR_USUARIO,
		&usuario.Email,
		&usuario.Senha_hash,
	)

	var matricula int

	err := row.Scan(&matricula)
	if err != nil {
		return 0, fmt.Errorf("erro ao criar usuário: %s", err.Error())
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
		&tarefa.Prazo,
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

func (tR *TaskRepository) LerMateria(ctx context.Context, matricula int, codigo int) (model.Materia, error) {
	row := tR.pool.QueryRow(ctx, model.LER_MATERIA, matricula, codigo)

	var materia model.Materia

	err := row.Scan(
		&materia.Codigo,
		&materia.Nome,
		&materia.Matricula,
	)
	if err != nil {
		return model.Materia{}, err
	}

	return materia, nil
}

func (tR *TaskRepository) ListarTarefas(ctx context.Context, codigo int) ([]model.Tarefa, error) {
	rows, err := tR.pool.Query(ctx, model.LISTAR_TAREFAS, codigo)
	if err != nil {
		fmt.Println(err)
		return []model.Tarefa{}, err
	}

	var tarefas []model.Tarefa
	var tarefa model.Tarefa

	for rows.Next() {
		err = rows.Scan(
			&tarefa.Id,
			&tarefa.Nome,
			&tarefa.Anotacao,
			&tarefa.Prazo,
			&tarefa.Codigo,
		)
		if err != nil {
			fmt.Println(err)
			return []model.Tarefa{}, err
		}

		tarefas = append(tarefas, tarefa)
	}

	rows.Close()

	return tarefas, nil
}

func (tR *TaskRepository) LerTarefa(ctx context.Context, codigo int, id int) (model.Tarefa, error) {
	row := tR.pool.QueryRow(ctx, model.LER_TAREFA, codigo, id)

	var tarefa model.Tarefa

	err := row.Scan(
		&tarefa.Id,
		&tarefa.Nome,
		&tarefa.Anotacao,
		&tarefa.Prazo,
		&tarefa.Codigo,
	)
	if err != nil {
		return model.Tarefa{}, err
	}

	return tarefa, nil
}

// UPDATE

func (tR *TaskRepository) MudarEmail(ctx context.Context, email string, matricula int) error {
	_, err := tR.pool.Exec(ctx, model.MUDAR_EMAIL_USR, email, matricula)
	if err != nil {
		return err
	}

	return nil
}

func (tR *TaskRepository) MudarSenha(ctx context.Context, senha_hash string, matricula int) error {
	_, err := tR.pool.Exec(ctx, model.MUDAR_SENHA_USR, senha_hash, matricula)
	if err != nil {
		return err
	}

	return nil
}

func (tR *TaskRepository) MudarNomeMateria(ctx context.Context, nome string, matricula int, codigo int) error {
	_, err := tR.pool.Exec(ctx, model.MUDAR_NOME_MAT, nome, matricula, codigo)
	if err != nil {
		return err
	}
	return nil
}

func (tR *TaskRepository) MudarNomeTarefa(ctx context.Context, nome string, id int, codigo int) error {
	_, err := tR.pool.Exec(ctx, model.MUDAR_NOME_TAR, nome, codigo, id)
	if err != nil {
		return err
	}
	return nil
}

func (tR *TaskRepository) MudarPrazoTarefa(ctx context.Context, prazo time.Time, codigo int, id int) error {
	_, err := tR.pool.Exec(ctx, model.MUDAR_PRAZ_TAR, prazo, codigo, id)
	if err != nil {
		return err
	}
	return nil
}

func (tR *TaskRepository) MudarAnotacaoTarefa(ctx context.Context, anotacao string, codigo int, id int) error {
	_, err := tR.pool.Exec(ctx, model.MUDAR_ANOT_TAR, anotacao, codigo, id)
	if err != nil {
		return err
	}
	return nil
}

// DELETE

func (tR *TaskRepository) DeletarUsuario(ctx context.Context, matricula int) error {
	_, err := tR.pool.Exec(ctx, model.DELETAR_USUARIO, matricula)
	if err != nil {
		return err
	}

	return nil
}

func (tR *TaskRepository) DeletarMateria(ctx context.Context, matricula int, codigo int) error {
	_, err := tR.pool.Exec(ctx, model.DELETAR_MATERIA, matricula, codigo)
	if err != nil {
		return err
	}

	return nil
}

func (tR *TaskRepository) DeletarTarefa(ctx context.Context, codigo int, id int) error {
	_, err := tR.pool.Exec(ctx, model.DELETAR_TAREFA, codigo, id)
	if err != nil {
		return err
	}

	return nil
}
