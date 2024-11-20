package service

import (
	"context"
	"todo-app/internal/domain"
	"todo-app/internal/repository"
)

type Todos interface {
	Create(ctx context.Context, todo domain.Todo) domain.Todo
	GetAll(ctx context.Context, page int, perPage int) []domain.Todo
	Update(ctx context.Context, id int, todo domain.Todo) domain.Todo
	Delete(ctx context.Context, id int)
}

type Services struct {
	Todos Todos
}

type Deps struct {
	Repos *repository.Repositories
}

func NewService(deps Deps) *Services {
	return &Services{
		Todos: NewTodosService(deps.Repos.Todos),
	}
}

