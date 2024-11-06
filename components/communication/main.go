package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/logging/common"
	"github.com/component-architecture-foundation/logging/output"
	"github.com/component-architecture-foundation/networking"
	"github.com/component-architecture-foundation/networking/coding"
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
	logger.Info(fmt.Sprintf("Starting server '%s@%s'...", shared.ServerName, shared.ServerVersion), false, shared.MainComponentName)

	var wg sync.WaitGroup

	coder := coding.JsonMessageCoder{Logger: logger}

	router := networking.NewRouter(logger, &wg, coder)
	router.Start()

	wg.Wait()
}
