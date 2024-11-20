package delivery

import (
	"net/http"
	v1 "todo-app/internal/delivery/http/v1"
	"todo-app/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Init() *chi.Mux {
	r := chi.NewRouter()
	// TODO: set my 'pkg/logger'
	r.Use(middleware.Logger)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	h.InitAPI(r)

	return r
}

func (h *Handler) InitAPI(router *chi.Mux) {
	handlerV1 := v1.NewHandler(h.services)

	router.Route("/api", func(r chi.Router) {
		handlerV1.Init(r)
	})
}

