package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/logging/common"
	"github.com/component-architecture-foundation/logging/output"
	common2 "github.com/component-architecture-foundation/protocols/common"
	"github.com/component-architecture-foundation/protocols/payloads"
	"github.com/component-architecture-foundation/protocols/serialization/json"
	"github.com/component-architecture-foundation/shared"
)

func getLogger() logging.HoornLogger {
	var userConfigDir, err = os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user config directory: %v", err)
	}

	var dir = filepath.Join(userConfigDir, "AppData", "Local")
	var logDir = dir + "\\Component Architecture Foundation\\logs\\communication_layer\\"

	return logging.NewHoornLogger(
		common.DEBUG,
		output.DefaultHoornLogOutput{},
		output.NewFileHoornLogOutput(
			logDir,
			5,
			true,
		))
}

func main() {
	logger := getLogger()
	logger.Info("Starting communication layer...", false, shared.MainComponentName)

	testPayload := payloads.EventPayload{
		Sender: common2.ComponentID{
			Name:         "Test",
			Language:     "Go",
			Version:      "0.0.0",
			Capabilities: nil,
		},
		Type: "Notification",
		Data: "Oh man, what a great test string!",
	}

	serializer := json.JSONSerializer{Logger: logger}
	serializedPayload, err := serializer.Serialize(&testPayload)

	if err != nil {
		return
	}

	logger.Info("Serialized payload: "+string(serializedPayload), false, shared.MainComponentName)

	deserializedPayload, err := serializer.Deserialize(serializedPayload)

	if err != nil {
		return
	}

	fmt.Printf("Deserialized payload: %+v\n", deserializedPayload)
}
