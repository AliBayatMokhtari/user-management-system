package router

import (
	"ums/handler"

	"github.com/go-chi/chi/v5"
)

func NewRouter(userHandler *handler.UserHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.Create)
		r.Get("/", userHandler.List)
		r.Get("/{id}", userHandler.GetByID)
		r.Put("/{id}", userHandler.Update)
		r.Delete("/{id}", userHandler.Delete)
	})

	return r
}
