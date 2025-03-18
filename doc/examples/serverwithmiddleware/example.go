package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pureapi/pureapi-core/doc/examples"
	"github.com/pureapi/pureapi-core/endpoint"
	endpointtypes "github.com/pureapi/pureapi-core/endpoint/types"
	"github.com/pureapi/pureapi-core/server"
	"github.com/pureapi/pureapi-core/util"
)

// This example demonstrates how to use middlewares in your endpoints.
// It shows how to create a stacks of middleware wrappers, modify them for
// specific uses and then use them in endpoints.
//
// To test out the panic in the "public" endpoint you can set the panic=true
// query parameter.
//
// To access the "secure" endpoint you can use curl:
// curl -H "X-Auth-Token: secret" localhost:8080/secure
func main() {
	// Create a shared stack of middleware wrappers.
	loggingWrapper := endpoint.NewWrapper("logging", LoggingMiddleware)
	recoveryWrapper := endpoint.NewWrapper("recovery", RecoveryMiddleware)
	sharedStack := endpoint.NewStack(loggingWrapper, recoveryWrapper)

	// Create a new stack for the "secure" endpoint.
	authWrapper := endpoint.NewWrapper("auth", AuthMiddleware)
	authStack := sharedStack.Clone()
	authStack.InsertAfter("recovery", authWrapper)

	// Create the endpoints with the speicific middlewares
	endpoints := []endpointtypes.Endpoint{
		endpoint.NewEndpoint("/public", http.MethodGet).WithHandler(
			func(w http.ResponseWriter, r *http.Request) {
				// Simulate a panic if "panic" query parameter is set.
				if r.URL.Query().Get("panic") == "true" {
					panic("Commencing a panic!")
				}
				fmt.Fprintf(w, "Hello, Public User!")
			},
		).WithMiddlewares(sharedStack.Middlewares()), // Use the shared stack.
		endpoint.NewEndpoint("/secure", http.MethodGet).WithHandler(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello, Secure User!")
			},
		).WithMiddlewares(authStack.Middlewares()), // Use the auth stack.
		endpoint.NewEndpoint("/boring", http.MethodGet).WithHandler(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello, Boring User!")
				w.WriteHeader(http.StatusTeapot)
			},
		), // No middlewares for the boring endpoint.
	}

	// Create the server handler with logger.
	emitterLogger := util.NewEmitterLogger(
		nil,
		examples.LoggerFactoryFn(),
	)
	handler := server.NewHandler(emitterLogger)

	// Create a HTTP server.
	httpServer := server.DefaultHTTPServer(handler, 8080, endpoints)

	// Start the server.
	if err := server.StartServer(handler, httpServer, nil); err != nil {
		panic(err)
	}
}

// LoggingMiddleware logs the incoming HTTP request.
//
// Parameters:
//   - next: Next HTTP handler.
//
// Returns:
//   - http.Handler
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request and time.
		start := time.Now()
		fmt.Printf(
			"[%s] Request: %s %s\n",
			time.Now().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
		)

		// Call the next middleware.
		next.ServeHTTP(w, r)

		// Log the response.
		delta := time.Since(start)
		fmt.Printf(
			"[%s] Request completed in %s\n",
			time.Now().Format(time.RFC3339),
			delta,
		)
	})
}

// AuthMiddleware simulates a simple authentication check.
//
// Parameters:
//   - next: Next HTTP handler.
//
// Returns:
//   - http.Handler
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")
		if token != "secret" {
			http.Error(w,
				"Unauthorized. You may not pass!",
				http.StatusUnauthorized,
			)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware recovers from panics and sets HTTP 500.
//
// Parameters:
//   - next: Next HTTP handler.
//
// Returns:
//   - http.Handler
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recovered from panic:", err)
				http.Error(
					w,
					"Internal Server Error :( Try again later.",
					http.StatusInternalServerError,
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
