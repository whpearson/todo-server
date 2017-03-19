package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"

	"github.com/whpearson/todo-server/restapi/operations"
	"github.com/whpearson/todo-server/restapi/operations/todos"
	"github.com/whpearson/todo-server/models"

)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target ../src/github.com/whpearson/todo-server --name  --spec ../go-swagger/examples/todo-list/swagger.yml

type okResp struct {
	code     int
	response interface{}
	headers  http.Header
}

func (e *okResp) WriteResponse(rw http.ResponseWriter , producer runtime.Producer) {
	for k, v := range e.headers {
		for _, val := range v {
			rw.Header().Add(k, val)
		}
	}
	if e.code > 0 {
		rw.WriteHeader(e.code)
	} else {
		rw.WriteHeader(http.StatusOK)
	}
	if err := producer.Produce(rw, e.response); err != nil {
		panic(err)
	}
}

// NotImplemented the error response when the response is not implemented
func Implemented(message string) middleware.Responder {
  var headers = make(http.Header)
  headers.Add("Content-Type", "application/json")
  return &okResp{http.StatusCreated, message, headers}
}




func configureFlags(api *operations.SimpleToDoListAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.SimpleToDoListAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

  api.SetDefaultProduces( "application/json" )

	// Applies when the "x-todolist-token" header is set
	api.KeyAuth = func(token string) (interface{}, error) {
		return nil, errors.NotImplemented("api key auth (key) x-todolist-token from header param [x-todolist-token] has not yet been implemented")
	}

	api.TodosAddOneHandler = todos.AddOneHandlerFunc(func(params todos.AddOneParams, principal interface{}) middleware.Responder {
		return todos.NewAddOneCreated().WithPayload(params.Body)
	})
	api.TodosDestroyOneHandler = todos.DestroyOneHandlerFunc(func(params todos.DestroyOneParams, principal interface{}) middleware.Responder {
	  return todos.NewDestroyOneNoContent()
	})
	var desc = "HI"
  api.TodosFindHandler = todos.FindHandlerFunc(func(params todos.FindParams, principal interface{}) middleware.Responder {
		return todos.NewFindOK().WithPayload([]*models.Item{&models.Item{
			Description:  &desc,
      ID:           1,
	}})
	})
	api.TodosUpdateOneHandler = todos.UpdateOneHandlerFunc(func(params todos.UpdateOneParams, principal interface{}) middleware.Responder {
		return todos.NewUpdateOneOK().WithPayload(params.Body)
	})

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
func configureServer(s *graceful.Server, scheme string) {
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
