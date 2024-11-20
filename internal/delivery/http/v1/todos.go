package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"todo-app/internal/domain"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) initTodosRoutes(r chi.Router) {
	r.Route("/todos", func(r chi.Router) {
		r.Get("/", h.getAllTodos)
		r.Post("/", h.createTodo)
		r.Patch("/{id}", h.updateTodo)
		r.Delete("/{id}", h.deleteTodo)
	})
}

func (h *Handler) getAllTodos(w http.ResponseWriter, r *http.Request) {
	var page, perPage = 1, 10

	if qPage := r.URL.Query().Get("page"); qPage != "" {
		numPage, err := strconv.Atoi(qPage)
		if err != nil {
			newErrorResponse(w, http.StatusInternalServerError, errorResponse{Message: "Failed to decode 'page' query param"})
			return
		}

		if numPage <= 0 {
			newErrorResponse(w, http.StatusBadRequest, errorResponse{Message: "Query param 'page' cannot be less then 1"})
			return
		}

		page = numPage
	}

	if qPerPage := r.URL.Query().Get("perPage"); qPerPage != "" {
		numPerPage, err := strconv.Atoi(qPerPage)
		if err != nil {
			newErrorResponse(w, http.StatusInternalServerError, errorResponse{Message: "Failed to decode 'perPage' query param"})
			return
		}

		if numPerPage <= 0 {
			newErrorResponse(w, http.StatusBadRequest, errorResponse{Message: "Query param 'perPage' cannot be less then 1"})
			return
		}

		perPage = numPerPage
	}

	todos := h.services.Todos.GetAll(context.Background(), page, perPage)
	// TODO: Set global?
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(dataResponse{Data: todos, Count: len(todos), CurrentPage: page})

	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, errorResponse{Message: "Failed to encode response"})
		return
	}
}

type createTodoInput struct {
	Name string
}

func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) {
	var input createTodoInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, errorResponse{Message: "invalid input body"})
		return
	}

	newTodo := h.services.Todos.Create(context.Background(), domain.Todo{Name: input.Name})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(newTodo)

	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, errorResponse{Message: "Failed to encode response"})
		return
	}
}

type updateTodoInput struct {
	Name        string
	IsCompleted bool
}

func (h *Handler) updateTodo(w http.ResponseWriter, r *http.Request) {
	var id int
	qId := chi.URLParam(r, "id")

	if qId == "" {
		newErrorResponse(w, http.StatusBadRequest, errorResponse{Message: "Invalid query param 'id'"})
		return
	}

	numId, err := strconv.Atoi(qId)
	if err != nil {
		newErrorResponse(w, http.StatusInternalServerError, errorResponse{Message: "Failed to decode 'id' query param"})
		return
	}

	if numId <= 0 {
		newErrorResponse(w, http.StatusBadRequest, errorResponse{Message: "Query param 'id' cannot be less then 1"})
		return
	}

	id = numId
	var input updateTodoInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		newErrorResponse(w, http.StatusBadRequest, errorResponse{Message: "invalid input body"})
		return
	}

	updatedTodo := h.services.Todos.Update(context.Background(), id, domain.Todo{Name: input.Name, IsCompleted: input.IsCompleted})
	w.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(w).Encode(updatedTodo)

	if encodeErr != nil {
		newErrorResponse(w, http.StatusInternalServerError, errorResponse{Message: "Failed to encode response"})
		return
	}
}

func (h *Handler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	qId := chi.URLParam(r, "id")

	if qId == "" {
		newErrorResponse(w, http.StatusBadRequest, errorResponse{Message: "Invalid query param 'id'"})
		return
	}

	numId, err := strconv.Atoi(qId)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, errorResponse{Message: "Failed to decode 'id' query param"})
		return
	}

	if numId <= 0 {
		newErrorResponse(w, http.StatusBadRequest, errorResponse{Message: "Query param 'id' cannot be less then 1"})
		return
	}

	h.services.Todos.Delete(context.Background(), numId)
}

