// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/rs/cors"

	"github.com/BNPrashanth/poc-docker-go-1/restapi/operations"
)

//go:generate swagger generate server --target ../../poc-docker-go-1 --name GoOne --spec ../spec.yml

func configureFlags(api *operations.GoOneAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.GoOneAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.DefaultHandler == nil {
		api.DefaultHandler = operations.DefaultHandlerFunc(func(params operations.DefaultParams) middleware.Responder {

			client := &http.Client{
				Timeout: 5 * time.Second,
			}
			resp := &operations.DefaultOKBody{}

			request, err := http.NewRequest("GET", "http://localhost:8082/", nil)
			request.Header.Set("Content-type", "application/json")

			if err != nil {
				fmt.Println("Error requesting provider service: " + err.Error())
				operations.NewDefaultOK().WithPayload(&operations.DefaultOKBody{
					Success: false,
					Data:    "Error",
				})
			}

			response, err := client.Do(request)
			if err != nil {
				fmt.Println("Error requesting provider service: " + err.Error())
				return operations.NewDefaultOK().WithPayload(&operations.DefaultOKBody{
					Success: false,
					Data:    "Error",
				})
			}

			data, _ := ioutil.ReadAll(response.Body)
			if json.Unmarshal(data, &resp) != nil {
				fmt.Println("Error trying to unmarshal response")
				return operations.NewDefaultOK().WithPayload(&operations.DefaultOKBody{
					Success: false,
					Data:    "Error",
				})
			}

			return operations.NewDefaultOK().WithPayload(&operations.DefaultOKBody{
				Success: true,
				Data:    resp.Data,
			})
		})
	}

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
	handleCORS := cors.AllowAll().Handler

	return handleCORS(handler)
}
