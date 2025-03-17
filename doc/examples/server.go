package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pureapi/pureapi-core/endpoint"
	endpointtypes "github.com/pureapi/pureapi-core/endpoint/types"
	"github.com/pureapi/pureapi-core/server"
	"github.com/pureapi/pureapi-core/util"
	utiltypes "github.com/pureapi/pureapi-core/util/types"
)

func Server() {
	eventEmitter := setupEventEmitter()
	emitterLogger := util.NewEmitterLogger(eventEmitter, nil)
	handler := server.NewHandler(emitterLogger)

	endpoints := []endpointtypes.Endpoint{
		endpoint.NewEndpoint("/hello", http.MethodGet).WithHandler(
			func(w http.ResponseWriter, r *http.Request) {
				log.Println("Incoming request")
				fmt.Fprintf(w, "Hello, PureAPI!")
			}),
	}

	instance := server.DefaultHTTPServer(handler, 8080, endpoints)

	if err := server.StartServer(handler, instance, nil); err != nil {
		panic(err)
	}
}

func setupEventEmitter() utiltypes.EventEmitter {
	eventEmitter := util.NewEventEmitter()
	eventEmitter.
		RegisterListener(
			server.EventStart,
			func(event *utiltypes.Event) {
				log.Printf("Event: %s\n", event.Message)
			},
		).
		RegisterListener(
			server.EventRegisterURL,
			func(event *utiltypes.Event) {
				log.Printf("Event: %s\n", event.Message)
			},
		)
	return eventEmitter
}
