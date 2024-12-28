package todosRepo

import (
	"context"
	"todo-app/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	connection *pgxpool.Pool
}

func New(connection *pgxpool.Pool) *Repository {
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

	return &Repository{
		connection: connection,
	}
}

func (repo *Repository) Create(ctx context.Context, todo Todo) (Todo, error) {
	var newTodo Todo
	err := repo.connection.QueryRow(
		ctx,
		"INSERT INTO todos (name) VALUES ($1) RETURNING id, name, is_completed, created_at",
		todo.Name).Scan(&newTodo.ID, &newTodo.Name, &newTodo.IsCompleted, &newTodo.CreatedAt)

	return newTodo, err
}

func (repo *Repository) GetAll(ctx context.Context, page int, perPage int) ([]Todo, int, error) {
	var todos []Todo
	rows, err := repo.connection.Query(ctx, "SELECT * FROM todoss LIMIT $1 OFFSET $2", perPage, (page-1)*perPage)
	if err != nil {
		return todos, -1, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.ID, &todo.Name, &todo.IsCompleted, &todo.CreatedAt)
		if err != nil {
			logger.Errorf("Row scan failed in GetAll Todos: %v", err)
			return todos, -1, err
		}
		todos = append(todos, todo)
	}

	var count int
	err = repo.connection.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM todos").Scan(&count)

	return todos, count, err
}

func (repo *Repository) Update(ctx context.Context, id int, todo Todo) (Todo, error) {
	var updatedTodo Todo
	err := repo.connection.QueryRow(
		ctx,
		`
			UPDATE todos SET name = COALESCE($1, todos.name), is_completed = COALESCE($2, todos.is_completed)
			 WHERE id = $3 RETURNING id, name, is_completed, created_at
		`,
		todo.Name, todo.IsCompleted, id).Scan(&updatedTodo.ID, &updatedTodo.Name, &updatedTodo.IsCompleted, &updatedTodo.CreatedAt)

	return updatedTodo, err
}

func (repo *Repository) Delete(ctx context.Context, id int) error {
	_, err := repo.connection.Exec(
		ctx,
		"DELETE FROM todos WHERE id = $1",
		id)

	return err
}
