package v1

import (
	"todo-app/internal/service"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init(r chi.Router) {
	r.Route("/v1", func(rr chi.Router) {
		h.initTodosRoutes(rr)
	})
}

