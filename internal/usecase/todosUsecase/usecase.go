package todosusecase

import (
	"context"
	"net/http"
	"todo-app/internal/repo/todosRepo"
	"todo-app/internal/usecase/helpers"
	"todo-app/pkg/logger"
)

type TodosRepository interface {
	Create(ctx context.Context, todo todosRepo.Todo) (todosRepo.Todo, error)
	GetAll(ctx context.Context, page int, perPage int) ([]todosRepo.Todo, int, error)
	Update(ctx context.Context, id int, todo todosRepo.Todo) (todosRepo.Todo, error)
	Delete(ctx context.Context, id int) error
}

type Usecase struct {
	repo TodosRepository
}

func New(todosRepo TodosRepository) *Usecase {
	return &Usecase{
		repo: todosRepo,
	}
}

func (u *Usecase) Create(ctx context.Context, todo todosRepo.Todo) (todosRepo.Todo, int) {
	todo, dbErr := u.repo.Create(ctx, todo)
	status, err := helpers.TranslatePgError(dbErr)

	if status != -1 && err != "" {
		logger.Errorf("Error while create Todo: %v", dbErr.Error())
	}

	return todo, status
}

func (u *Usecase) GetAll(ctx context.Context, page int, perPage int) ([]todosRepo.Todo, int, int) {
	todos, count, dbErr := u.repo.GetAll(ctx, page, perPage)
	
	errorStatusCode := -1
	if (dbErr != nil) {
		logger.Errorf("Error while Get all Todos: %v", dbErr.Error())
		errorStatusCode = http.StatusInternalServerError
	}

	return todos, count, errorStatusCode
}

func (u *Usecase) Update(ctx context.Context, id int, todo todosRepo.Todo) (todosRepo.Todo, int) {
	todo, dbErr := u.repo.Update(ctx, id, todo)
	status, err := helpers.TranslatePgError(dbErr)

	if status != -1 && err != "" {
		logger.Errorf("Error while update Todo: %v", dbErr.Error())
	}

	return todo, status
}

func (u *Usecase) Delete(ctx context.Context, id int) int {
	dbErr := u.repo.Delete(ctx, id)
	status, err := helpers.TranslatePgError(dbErr)
	
	if status != -1 && err != "" {
		logger.Errorf("Error while update Todo: %v", dbErr.Error())
	}

	return status
}
