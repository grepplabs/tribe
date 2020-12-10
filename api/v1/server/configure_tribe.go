// This file is safe to edit. Once it exists it will not be overwritten

package server

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/grepplabs/tribe/api/v1/server/restapi"
	"github.com/grepplabs/tribe/api/v1/server/restapi/healthz"
	"github.com/grepplabs/tribe/api/v1/server/restapi/realms"
	"github.com/grepplabs/tribe/api/v1/server/restapi/users"
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
	if api.RealmsCreateRealmHandler == nil {
		api.RealmsCreateRealmHandler = realms.CreateRealmHandlerFunc(func(params realms.CreateRealmParams) middleware.Responder {
			return middleware.NotImplemented("operation realms.CreateRealm has not yet been implemented")
		})
	}
	if api.UsersCreateUserHandler == nil {
		api.UsersCreateUserHandler = users.CreateUserHandlerFunc(func(params users.CreateUserParams) middleware.Responder {
			return middleware.NotImplemented("operation users.CreateUser has not yet been implemented")
		})
	}
	if api.RealmsGetRealmHandler == nil {
		api.RealmsGetRealmHandler = realms.GetRealmHandlerFunc(func(params realms.GetRealmParams) middleware.Responder {
			return middleware.NotImplemented("operation realms.GetRealm has not yet been implemented")
		})
	}
	if api.UsersGetUserHandler == nil {
		api.UsersGetUserHandler = users.GetUserHandlerFunc(func(params users.GetUserParams) middleware.Responder {
			return middleware.NotImplemented("operation users.GetUser has not yet been implemented")
		})
	}
	if api.UsersGetUsersHandler == nil {
		api.UsersGetUsersHandler = users.GetUsersHandlerFunc(func(params users.GetUsersParams) middleware.Responder {
			return middleware.NotImplemented("operation users.GetUsers has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
