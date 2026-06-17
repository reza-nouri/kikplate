package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kickplate/api/handler/handlers"
	"github.com/kickplate/api/lib"
)

type GeneratorRoutes struct {
	handler handlers.GeneratorHandler
	rh      lib.RequestHandler
}

func NewGeneratorRoutes(handler handlers.GeneratorHandler, rh lib.RequestHandler) GeneratorRoutes {
	return GeneratorRoutes{handler: handler, rh: rh}
}

func (r GeneratorRoutes) Setup() {
	r.rh.Mux.Route("/generate", func(m chi.Router) {
		m.Get("/{slug}/schema", r.handler.Schema)
		m.Post("/{slug}", r.handler.Generate)
	})
}
