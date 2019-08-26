package server

import (
	"context"
	"fmt"
	"github.com/alfcope/checkouttest/api"
	"github.com/alfcope/checkouttest/config"
	"github.com/alfcope/checkouttest/datasource"
	"github.com/alfcope/checkouttest/pkg/logging"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type checkoutApi struct {
	routes *mux.Router

	controller *api.CheckoutController
	service    *api.CheckoutService
}

// Creates an instance of the api endpoints
func NewCheckoutApi(configuration config.Configuration) (*checkoutApi, error) {

	ds, err := datasource.InitInMemoryDatasource(configuration.Data)
	if err != nil {
		fmt.Println("Error initiating datasource: ", err.Error())
		return nil, err
	}

	checkoutService := api.NewCheckoutService(ds)

	routes := mux.NewRouter()
	apiRoute := routes.PathPrefix("/api/v1").Subrouter().StrictSlash(true)

	api.AddHealthCheckRoute(apiRoute)

	return &checkoutApi{
		routes:     apiRoute,
		controller: api.NewCheckoutController(apiRoute, checkoutService),
		service:    &checkoutService,
	}, nil
}

// Start the http server
func (c checkoutApi) RunServer(port int) {

	var server = &http.Server{
		Addr:           fmt.Sprintf(":%v", port),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20, // Max header of 1MB,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// sigterm signal sent from kubernetes
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		// We received an interrupt signal, shut down.
		if err := server.Shutdown(context.Background()); err != nil {
			logging.Logger.Errorf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	corsHandler := handlers.CORS(
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With", "Authorization"}))

	server.Handler = corsHandler(c.routes)

	logging.Logger.Info("Starting HTTP service at ", port)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logging.Logger.Errorf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
