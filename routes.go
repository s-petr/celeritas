package celeritas

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (c *Celeritas) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Recoverer)
	mux.Use(c.NoSurf)
	mux.Use(c.SessionLoad)
	mux.Use(c.CheckForMaintenanceMode)

	if c.Debug {
		mux.Use(middleware.Logger)
	}

	return mux
}
