package service

import (
	"context"
	"todo-app/internal/domain"
	"todo-app/internal/repository"
)

type TodosService struct {
	repo repository.Todos
}

func NewTodosService(repo repository.Todos) *TodosService {
	return &TodosService{
		repo: repo,
	}
}

func (s *TodosService) Create(ctx context.Context, todo domain.Todo) domain.Todo {
	todo, _ = s.repo.Create(ctx, todo)

	// TODO: pls help
	// if err != nil {
	// 	return nil, error{Message: "SQL-error while creating Todo"}
	// }

	return todo
}

func (s *TodosService) GetAll(ctx context.Context, page int, perPage int) []domain.Todo {
	todos, _ := s.repo.GetAll(ctx, page, perPage)
	// TODO: handle error
	return todos
}

func (s *TodosService) Update(ctx context.Context, id int, todo domain.Todo) domain.Todo {
	todo, _ = s.repo.Update(ctx, id, todo)
	// TODO: handle error
	return todo
}

func (s *TodosService) Delete(ctx context.Context, id int) {
	s.repo.Delete(ctx, id)
	// TODO: handle error
}

