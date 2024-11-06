package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/logging/common"
	"github.com/component-architecture-foundation/logging/output"
	"github.com/component-architecture-foundation/networking"
	"github.com/component-architecture-foundation/networking/authentication"
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
	logger.Info("Starting communication layer...", false, shared.MainComponentName)

	var wg sync.WaitGroup

	authenticator := &authentication.WhitelistAuthenticator{Logger: logger}
	coder := coding.JsonMessageCoder{Logger: logger}

	router := networking.NewRouter(logger, &wg, authenticator, coder)
	router.Start()

	wg.Wait()
}
