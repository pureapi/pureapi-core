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

// This example demonstrates how to use the server package to start a server.
func main() {
	// Create the endpoints.
	endpoints := []endpointtypes.Endpoint{
		endpoint.NewEndpoint("/hello", http.MethodGet).WithHandler(
			func(w http.ResponseWriter, r *http.Request) {
				log.Println("Incoming request")
				fmt.Fprintf(w, "Hello, PureAPI!")
			}),
	}

	// Create the server handler.
	port := 8080
	eventEmitter := SetupEventEmitter(port)
	emitterLogger := util.NewEmitterLogger(eventEmitter, nil)
	handler := server.NewHandler(emitterLogger)

	// Create a HTTP server.
	httpServer := server.DefaultHTTPServer(handler, port, endpoints)

	// Start the server.
	if err := server.StartServer(handler, httpServer, nil); err != nil {
		panic(fmt.Errorf("server panic: %w", err))
	}
}

// SetupEventEmitter sets up an event emitter for the server. It demonstrates
// how to register event listeners. For server there are more events available.
// See the server package for more information.
//
// Parameters:
//   - port: Port for the server.
//
// Returns:
//   - utiltypes.EventEmitter: The event emitter.
func SetupEventEmitter(port int) utiltypes.EventEmitter {
	eventEmitter := util.NewEventEmitter()
	eventEmitter.
		RegisterListener(
			server.EventStart,
			func(event *utiltypes.Event) {
				// Using event message directly for logging.
				log.Printf("Event: %s, port: %d\n", event.Message, port)
			},
		).
		RegisterListener(
			server.EventRegisterURL,
			func(event *utiltypes.Event) {
				// Using event data to log the path and methods.
				data := event.Data.(map[string]any)
				log.Printf(
					"Event: Registering URL %s %v\n",
					data["path"],
					data["methods"],
				)
			},
		)
	return eventEmitter
}
