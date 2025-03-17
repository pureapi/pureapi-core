package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/pureapi/pureapi-core/endpoint"
	"github.com/pureapi/pureapi-core/endpoint/types"
	"github.com/pureapi/pureapi-core/server"
	"github.com/pureapi/pureapi-core/util"
)

// LoggingMiddleware logs the incoming HTTP request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s\n", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware simulates a simple authentication check.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")
		if token != "secret-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware recovers from panics and returns HTTP 500.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recovered from panic:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func RunMiddleware() {
	eventEmitter := setupEventEmitter()
	emitterLogger := util.NewEmitterLogger(eventEmitter, nil)
	handler := server.NewHandler(emitterLogger)

	loggingWrapper := endpoint.NewWrapper("logging", LoggingMiddleware)
	recoveryWrapper := endpoint.NewWrapper("recovery", RecoveryMiddleware)
	commonStack := endpoint.NewStack(loggingWrapper, recoveryWrapper)

	authWrapper := endpoint.NewWrapper("auth", AuthMiddleware)
	authStack := commonStack.Clone()
	authStack.InsertBefore("recovery", authWrapper)

	endpoints := []types.Endpoint{
		endpoint.NewEndpoint("/public", http.MethodGet).WithHandler(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello, Public User!")
			},
		),
		endpoint.NewEndpoint("/secure", http.MethodGet).WithHandler(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello, Secure User!")
			},
		),
	}

	instance := server.DefaultHTTPServer(handler, 8080, endpoints)

	if err := server.StartServer(handler, instance, nil); err != nil {
		panic(err)
	}
}
