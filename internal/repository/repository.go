package repository

import (
	"context"
	"todo-app/internal/domain"

	"github.com/jackc/pgx/v5"
)

type Todos interface {
	Create(ctx context.Context, todo domain.Todo) (domain.Todo, error)
	GetAll(ctx context.Context, page int, perPage int) ([]domain.Todo, error)
	Update(ctx context.Context, id int, todo domain.Todo) (domain.Todo, error)
	Delete(ctx context.Context, id int) error
}

type Repositories struct {
	Todos Todos
}

func NewRepositories(connection *pgx.Conn) *Repositories {
	return &Repositories{
		Todos: NewTodosRepo(connection),
	}
}

