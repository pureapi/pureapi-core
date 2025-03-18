package examples

import (
	"log"

	"github.com/pureapi/pureapi-core/server"
	"github.com/pureapi/pureapi-core/util"
	utiltypes "github.com/pureapi/pureapi-core/util/types"
)

// SetupEventEmitter sets up an event emitter for the server. It logs some
// server events for demonstration.
//
// Returns:
//   - utiltypes.EventEmitter: The event emitter.
func SetupEventEmitter() utiltypes.EventEmitter {
	eventEmitter := util.NewEventEmitter()
	eventEmitter.
		RegisterListener(
			server.EventStart,
			func(event *utiltypes.Event) {
				log.Println(event.Message)
			},
		).
		RegisterListener(
			server.EventRegisterURL,
			func(event *utiltypes.Event) {
				log.Println(event.Message)
			},
		)
	return eventEmitter
}
