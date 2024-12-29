package todosRepo

import (
	"context"

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
	tx, txErr := repo.connection.Begin(ctx)

	if txErr != nil {
		return Todo{}, txErr
	}

	err := repo.connection.QueryRow(
		ctx,
		"INSERT INTO todos (name) VALUES ($1) RETURNING id, name, is_completed, created_at",
		todo.Name).Scan(&newTodo.ID, &newTodo.Name, &newTodo.IsCompleted, &newTodo.CreatedAt)

	if err != nil {
		tx.Rollback(ctx)
	} else {
		tx.Commit(ctx)
	}

	return newTodo, err
}

func (repo *Repository) GetAll(ctx context.Context, page int, perPage int) ([]Todo, int, error) {
	var todos []Todo

	tx, txErr := repo.connection.Begin(ctx)
	if txErr != nil {
		return todos, -1, txErr
	}

	rows, err := repo.connection.Query(ctx, "SELECT * FROM todoss LIMIT $1 OFFSET $2", perPage, (page-1)*perPage)
	if err != nil {
		tx.Rollback(ctx)
		return todos, -1, err
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		err = rows.Scan(&todo.ID, &todo.Name, &todo.IsCompleted, &todo.CreatedAt)
		if err != nil {
			tx.Rollback(ctx)
			return todos, -1, err
		}
		todos = append(todos, todo)
	}

	var count int
	err = repo.connection.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM todos").Scan(&count)

	if err != nil {
		tx.Rollback(ctx)
	} else {
		tx.Commit(ctx)
	}

	return todos, count, err
}

func (repo *Repository) Update(ctx context.Context, id int, todo Todo) (Todo, error) {
	var updatedTodo Todo

	tx, txErr := repo.connection.Begin(ctx)
	if txErr != nil {
		return updatedTodo, txErr
	}

	err := repo.connection.QueryRow(
		ctx,
		`
			UPDATE todos SET name = COALESCE($1, todos.name), is_completed = COALESCE($2, todos.is_completed)
			 WHERE id = $3 RETURNING id, name, is_completed, created_at
		`,
		todo.Name, todo.IsCompleted, id).Scan(&updatedTodo.ID, &updatedTodo.Name, &updatedTodo.IsCompleted, &updatedTodo.CreatedAt)

	if err != nil {
		tx.Rollback(ctx)
	} else {
		tx.Commit(ctx)
	}

	return updatedTodo, err
}

func (repo *Repository) Delete(ctx context.Context, id int) error {
	tx, txErr := repo.connection.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	_, err := repo.connection.Exec(
		ctx,
		"DELETE FROM todos WHERE id = $1",
		id)

	if err != nil {
		tx.Rollback(ctx)
	} else {
		tx.Commit(ctx)
	}

	return err
}
