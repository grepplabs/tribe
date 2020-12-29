// This file is safe to edit. Once it exists it will not be overwritten

package server

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"os"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/gorilla/handlers"

	"github.com/grepplabs/tribe/api/v1/server/restapi"
	"github.com/grepplabs/tribe/api/v1/server/restapi/healthz"
)

//go:generate swagger generate server --target ../../v1 --name Tribe --spec ../openapi.yaml --api-package restapi --server-package server --principal interface{} --exclude-main

func configureFlags(api *restapi.TribeAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *restapi.TribeAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.HealthzGetHealthyHandler == nil {
		api.HealthzGetHealthyHandler = healthz.GetHealthyHandlerFunc(func(params healthz.GetHealthyParams) middleware.Responder {
			return middleware.NotImplemented("operation healthz.GetHealthy has not yet been implemented")
		})
	}
	if api.HealthzGetReadyHandler == nil {
		api.HealthzGetReadyHandler = healthz.GetReadyHandlerFunc(func(params healthz.GetReadyParams) middleware.Responder {
			return middleware.NotImplemented("operation healthz.GetReady has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {
		// canceling server context
		serverCancel()
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

var (
	// ServerCtx and ServerCancel
	ServerCtx, serverCancel = context.WithCancel(context.Background())
)

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
	s.BaseContext = func(_ net.Listener) context.Context {
		return ServerCtx
	}
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {

	//TODO: see also https://github.com/didip/tollbooth
	return handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(handlers.LoggingHandler(os.Stdout, handler))
}
