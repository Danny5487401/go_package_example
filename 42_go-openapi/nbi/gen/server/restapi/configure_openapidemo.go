// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/Danny5487401/go_package_example/42_go-openapi/handlers"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/Danny5487401/go_package_example/42_go-openapi/nbi/gen/server/restapi/operations"
)

//go:generate swagger generate server --target ../../server --name Openapidemo --spec ../../../nbi-swagger.yaml --principal any

func configureFlags(api *operations.OpenapidemoAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
	_ = api
}

func configureAPI(api *operations.OpenapidemoAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...any)
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.AddPhoneBookEntryHandler = operations.AddPhoneBookEntryHandlerFunc(func(params operations.AddPhoneBookEntryParams) middleware.Responder {
		return handlers.AddPhoneBookEntry(params)
	})
	api.AddPostObjectHandler = operations.AddPostObjectHandlerFunc(func(params operations.AddPostObjectParams) middleware.Responder {
		return handlers.AddPostObject(params)
	})
	api.GetHostInfoHandler = operations.GetHostInfoHandlerFunc(func(params operations.GetHostInfoParams) middleware.Responder {
		return handlers.GetHostInfo(params)
	})
	api.GetPhoneBookHandler = operations.GetPhoneBookHandlerFunc(func(params operations.GetPhoneBookParams) middleware.Responder {
		return handlers.GetPhoneBook(params)
	})
	api.GetPhoneBookEntryHandler = operations.GetPhoneBookEntryHandlerFunc(func(params operations.GetPhoneBookEntryParams) middleware.Responder {
		return handlers.GetPhoneBookEntry(params)
	})
	api.GetPostTitlesHandler = operations.GetPostTitlesHandlerFunc(func(params operations.GetPostTitlesParams) middleware.Responder {
		return handlers.GetPostTitles(params)
	})
	api.GetPostsByUserHandler = operations.GetPostsByUserHandlerFunc(func(params operations.GetPostsByUserParams) middleware.Responder {
		return handlers.GetPostsByUser(params)
	})

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
	_ = tlsConfig
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(server *http.Server, scheme, addr string) {
	_ = server
	_ = scheme
	_ = addr
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
