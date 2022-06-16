package application

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/sioncojp/famili-api/utils/config"
	"github.com/sioncojp/famili-api/utils/log"
)

// NewRouter...routerを初期化して返す
func (s *HttpHandler) NewRouter() {
	r := chi.NewRouter()
	newMiddlewares(r, s.AppConfig)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/todos", func(r chi.Router) {
			r.Get("/", s.Router.V1.TodosHandler.List)
			r.Post("/", s.Router.V1.TodosHandler.Create)
			r.Route("/{id}", func(r chi.Router) {
				r.Use(s.Router.V1.TodosHandler.Ctx)
				r.Put("/", s.Router.V1.TodosHandler.Update)
				r.Delete("/", s.Router.V1.TodosHandler.Delete)
			})
		})
	})

	s.ServeMux = r
}

// newMiddlewares...routerで利用するミドルウェア. https://github.com/go-chi/chi#core-middlewares
func newMiddlewares(r *chi.Mux, c *config.AppConfig) {
	r.Use(middleware.GetHead)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(log.NewChiLogger(c.Server.Name, c.Service.Env))

	// healthcheckはhealthz。k8sを想定してこのネーミングにしている
	r.Use(middleware.Heartbeat("/healthz"))
}
