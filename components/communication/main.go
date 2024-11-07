package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/component-architecture-foundation/logging"
	"github.com/component-architecture-foundation/logging/common"
	"github.com/component-architecture-foundation/logging/output"
	"github.com/component-architecture-foundation/networking"
	"github.com/component-architecture-foundation/networking/coding"
	"github.com/component-architecture-foundation/shared"
)

func getLogDir() string {
	var userConfigDir, err = os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user config directory: %v", err)
	}

	var dir = filepath.Join(userConfigDir, "AppData", "Local")
	var logDir = dir + shared.RootLogDir
	return logDir
}

func getLogger() logging.HoornLogger {
	logDir := getLogDir() + "\\Communication_Layer\\"

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

	coder := coding.JsonMessageCoder{Logger: &logger}

	setupLoggingPayload := fmt.Sprintf(`{
		"action": "setup_logging",
		"args": [
			{
				"type": "string",
				"value": "%s"
			},
			{
				"type": "int",
				"value": "%d"
			},
            {
                "type": "string",
                "value": "%s"
            }
		]
	}`, strings.ReplaceAll(getLogDir()+"\\Components\\", `\`, `\\`), 5, "debug")

	router := networking.NewRouter(&logger, &wg, coder, []byte(setupLoggingPayload))
	router.Start()

	wg.Wait()
}
