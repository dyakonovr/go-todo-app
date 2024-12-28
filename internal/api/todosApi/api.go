package todosapi

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"todo-app/internal/api"
	"todo-app/internal/repo/todosRepo"
	todosusecase "todo-app/internal/usecase/todosUsecase"

	"github.com/go-chi/chi/v5"
)

type Api struct {
	todosUsecase *todosusecase.Usecase
}

func New(todosUsecase *todosusecase.Usecase) *Api {
	return &Api{
		todosUsecase: todosUsecase,
	}
}

func (a *Api) InitApi(r chi.Router) {
	r.Route("/todos", func(r chi.Router) {
		r.Get("/", a.getAllTodos)
		r.Post("/", a.createTodo)
		r.Patch("/{id}", a.updateTodo)
		r.Delete("/{id}", a.deleteTodo)
	})
}

func (a *Api) getAllTodos(w http.ResponseWriter, r *http.Request) {
	var page, perPage = 1, 10

	if qPage := r.URL.Query().Get("page"); qPage != "" {
		numPage, err := strconv.Atoi(qPage)
		if err != nil {
			api.NewErrorResponse(w, http.StatusBadRequest, api.ErrorResponse{Message: "Failed to decode 'page' query param"})
			return
		}

		if numPage <= 0 {
			page = 1
		} else {
			page = numPage
		}

	}

	if qPerPage := r.URL.Query().Get("perPage"); qPerPage != "" {
		numPerPage, err := strconv.Atoi(qPerPage)
		if err != nil {
			api.NewErrorResponse(w, http.StatusBadRequest, api.ErrorResponse{Message: "Failed to decode 'perPage' query param"})
			return
		}

		if numPerPage <= 0 {
			perPage = 1
		} else {
			perPage = numPerPage
		}
	}

	todos, count, errorStatusCode := a.todosUsecase.GetAll(r.Context(), page, perPage)

	if errorStatusCode != -1 {
		api.NewErrorResponse(w, errorStatusCode, api.ErrorResponse{})
		return
	}

	err := json.NewEncoder(w).Encode(api.DataResponse{Data: todos, TotalPages: int(math.Ceil(float64(count) / float64(perPage))), CurrentPage: page})

	if err != nil {
		api.NewErrorResponse(w, http.StatusInternalServerError, api.ErrorResponse{Message: "Failed to encode response"})
		return
	}
}

func (a *Api) createTodo(w http.ResponseWriter, r *http.Request) {
	var input createTodoInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		api.NewErrorResponse(w, http.StatusBadRequest, api.ErrorResponse{Message: "invalid input body"})
		return
	}

	newTodo, errorStatusCode := a.todosUsecase.Create(r.Context(), todosRepo.Todo{Name: input.Name})

	if errorStatusCode != -1 {
		api.NewErrorResponse(w, errorStatusCode, api.ErrorResponse{Message: "Error creating Todo. Try later"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(newTodo)

	if err != nil {
		api.NewErrorResponse(w, http.StatusInternalServerError, api.ErrorResponse{Message: "Failed to encode response"})
		return
	}
}

func (a *Api) updateTodo(w http.ResponseWriter, r *http.Request) {
	var id int
	qId := chi.URLParam(r, "id")

	if qId == "" {
		api.NewErrorResponse(w, http.StatusBadRequest, api.ErrorResponse{Message: "Invalid query param 'id'"})
		return
	}

	numId, err := strconv.Atoi(qId)
	if err != nil {
		api.NewErrorResponse(w, http.StatusBadRequest, api.ErrorResponse{Message: "Failed to decode 'id' query param"})
		return
	}

	if numId <= 0 {
		api.NewErrorResponse(w, http.StatusBadRequest, api.ErrorResponse{Message: "Query param 'id' cannot be less then 1"})
		return
	}

	id = numId
	var input updateTodoInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		api.NewErrorResponse(w, http.StatusBadRequest, api.ErrorResponse{Message: "invalid input body"})
		return
	}

	updatedTodo, errorStatusCode := a.todosUsecase.Update(r.Context(), id, todosRepo.Todo{Name: input.Name, IsCompleted: input.IsCompleted})

	if errorStatusCode != -1 {
		api.NewErrorResponse(w, errorStatusCode, api.ErrorResponse{Message: "Error updating todo. Try later"})
		return
	}

	encodeErr := json.NewEncoder(w).Encode(updatedTodo)

	if encodeErr != nil {
		api.NewErrorResponse(w, http.StatusInternalServerError, api.ErrorResponse{Message: "Failed to encode response"})
		return
	}
}

func (a *Api) deleteTodo(w http.ResponseWriter, r *http.Request) {
	qId := chi.URLParam(r, "id")

	if qId == "" {
		api.NewErrorResponse(w, http.StatusBadRequest, api.ErrorResponse{Message: "Invalid query param 'id'"})
		return
	}

	numId, err := strconv.Atoi(qId)
	if err != nil {
		api.NewErrorResponse(w, http.StatusBadRequest, api.ErrorResponse{Message: "Failed to decode 'id' query param"})
		return
	}

	if numId <= 0 {
		api.NewErrorResponse(w, http.StatusBadRequest, api.ErrorResponse{Message: "Query param 'id' cannot be less then 1"})
		return
	}

	errorStatusCode := a.todosUsecase.Delete(r.Context(), numId)

	if errorStatusCode != -1 {
		api.NewErrorResponse(w, errorStatusCode, api.ErrorResponse{Message: "Error deleting Todo. Try later"})
		return
	}
}
