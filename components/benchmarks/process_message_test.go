package benchmarks

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/benchmarks/_internal"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/coding"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/transport"
)

const defaultMessagePayload = `{
  "action": "log_debug",
  "args": [
    {
      "type": "string",
      "value": "Called log_debug function"
    },
    {
      "type": "bool",
      "value": "false"
    },
    {
      "type": "string",
      "value": ""
    },
    {
      "type": "string",
      "value": "Benchmark.Debugging"
    }
  ]
}`

func processMessage(message transport.Message, id transport.ComponentID, router *networking.Router) {
	router.ProcessMessage(message, id, nil)
}

func BenchmarkProcessMessageSuite(b *testing.B) {
	logger := _internal.GetLogger("Component-Based Architecture Foundation - Benchmarks")
	var wg sync.WaitGroup
	coder := coding.JsonMessageCoder{Logger: &logger}
	router := networking.NewRouter(&logger, &wg, &coder)
	router.Start()

	// Wait for logger registration (adapt as needed)
	for {
		if router.IsComponentRegistered("ea1973db-31e7-4fe4-bd57-e217f246f6a1") {
			logger.Info(fmt.Sprintf("Logger Registered"), false, "Benchmark.Debugging")
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	requesterID := transport.ComponentID{
		Title:             "Benchmarks",
		Version:           "1.0.0",
		Capabilities:      nil,
		ComponentUniqueID: "aea3a1e6-47e3-4d13-a308-83eb640a51f5",
	}

	payload, _ := transport.MessagePayloadFromBytes([]byte(defaultMessagePayload))

	message := transport.NewMessage(requesterID, payload, "")

	b.Run("ProcessMessage", func(b *testing.B) {
		benchmarkProcessMessage(b, router, *message, requesterID)
	})

	router.Stop()
	router.WaitUntilShutdown()
}

func benchmarkProcessMessage(b *testing.B, router *networking.Router, message transport.Message, requesterID transport.ComponentID) {
	for i := 0; i < b.N; i++ {
		processMessage(message, requesterID, router)
	}
}
