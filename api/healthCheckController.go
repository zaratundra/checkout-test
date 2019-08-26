package api

import (
	"github.com/etherlabsio/healthcheck"
	"github.com/gorilla/mux"
	"time"
)

func AddHealthCheckRoute(router *mux.Router) {

	healthCheckRouter := router.PathPrefix("/healthcheck").Subrouter()

	healthCheckRouter.Handle("/", healthcheck.Handler(
		healthcheck.WithTimeout(2*time.Second),
	))
}
