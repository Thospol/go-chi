package routes

import (
	"compress/flate"

	"saaa-api/internal/core/config"
	"saaa-api/internal/handlers/middlewares"
	"saaa-api/internal/pkg/guest"
	"saaa-api/internal/pkg/healthcheck"

	"github.com/NYTimes/gziphandler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter new router
func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(gziphandler.MustNewGzipLevelHandler(flate.BestSpeed))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/", func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Use(middlewares.AcceptLanguage)
			r.Use(middlewares.Request)

			healthCheckEndpoint := healthcheck.NewEndpoint(config.CF, config.RR)
			guestEndpoint := guest.NewEndpoint(config.CF, config.RR)

			r.Get("/healthz", healthCheckEndpoint.HealthCheck())
			r.Route("/v1", func(r chi.Router) {
				if config.CF.SwaggerInfo.Enable {
					r.Get("/swagger/*", httpSwagger.Handler(
						httpSwagger.URL("doc.json"),
					))
				}

				// c is general client.
				r.Route("/c", func(r chi.Router) {
					r.Route("/guest", func(r chi.Router) {
						r.With(middlewares.Transaction).Post("/refreshToken", guestEndpoint.RefreshToken())
					})
				})
			})
		})
	})

	return r
}
