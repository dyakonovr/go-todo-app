package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ApiInterface interface {
	InitApi(r chi.Router)
}

func Init(apis []ApiInterface) *chi.Mux {
	r := chi.NewRouter()
	
	r.Use(helperMiddleware)
	r.Use(loggerMiddleware)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	initAPI(r, apis)
	return r
}

func initAPI(router *chi.Mux, apis []ApiInterface) {
	router.Route("/api/v1", func(r chi.Router) {
		for _, api := range apis {
			api.InitApi(r)
		}
	})
}
