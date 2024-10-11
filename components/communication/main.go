package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/logging/common"
	"github.com/component-architecture-foundation/logging/output"
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
	logger.Info("Test", false, shared.MainComponentName)
}
