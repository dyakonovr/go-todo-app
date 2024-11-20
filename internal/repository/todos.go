package repository

import (
	"context"
	"todo-app/internal/domain"
	"todo-app/pkg/logger"

	"github.com/jackc/pgx/v5"
)

type TodosRepository struct {
	connection *pgx.Conn
}

func NewTodosRepo(connection *pgx.Conn) *TodosRepository {
	connection.Exec(
		context.Background(),
		`
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			is_completed BOOLEAN DEFAULT false,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
		);
	`)

	return &TodosRepository{
		connection: connection,
	}
}

func (repo *TodosRepository) Create(ctx context.Context, todo domain.Todo) (domain.Todo, error) {
	var newTodo domain.Todo
	err := repo.connection.QueryRow(
		ctx,
		"INSERT INTO todos (name) VALUES ($1) RETURNING id, name, is_completed, created_at",
		todo.Name).Scan(&newTodo.ID, &newTodo.Name, &newTodo.IsCompleted, &newTodo.CreatedAt)

	return newTodo, err
}

func (repo *TodosRepository) GetAll(ctx context.Context, page int, perPage int) ([]domain.Todo, error) {
	var todos []domain.Todo
	rows, err := repo.connection.Query(ctx, "SELECT * FROM todos LIMIT $1 OFFSET $2", perPage, (page-1)*perPage)
	if err != nil {
		return todos, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo domain.Todo
		err = rows.Scan(&todo.ID, &todo.Name, &todo.IsCompleted, &todo.CreatedAt)
		if err != nil {
			logger.Errorf("Row scan failed in GetAll Todos: %v", err)
			return todos, err
		}
		todos = append(todos, todo)
	}

	return todos, err
}

func (repo *TodosRepository) Update(ctx context.Context, id int, todo domain.Todo) (domain.Todo, error) {
	var updatedTodo domain.Todo
	err := repo.connection.QueryRow(
		ctx,
		`
			UPDATE todos SET name = COALESCE($1, todos.name), is_completed = COALESCE($2, todos.is_completed)
			 WHERE id = $3 RETURNING id, name, is_completed, created_at
		`,
		todo.Name, todo.IsCompleted, id).Scan(&updatedTodo.ID, &updatedTodo.Name, &updatedTodo.IsCompleted, &updatedTodo.CreatedAt)

	return updatedTodo, err
}

func (repo *TodosRepository) Delete(ctx context.Context, id int) error {
	_, err := repo.connection.Exec(
		ctx,
		"DELETE FROM todos WHERE id = $1",
		id)

	return err
}


