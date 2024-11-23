package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"sync"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/common"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/output"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/networking/coding"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/shared"
)

func getLogDir(applicationName string) string {
	var userConfigDir, err = os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user config directory: %v", err)
	}

	var dir = filepath.Join(userConfigDir, "AppData", "Local")
	var logDir = dir + "\\" + applicationName + "\\logs"
	return logDir
}

func getLogger(applicationName string) logging.HoornLogger {
	logDir := getLogDir(applicationName) + "\\Communication_Layer\\"

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
	applicationArgs := os.Args

	if len(applicationArgs) < 2 {
		log.Fatalf("Usage: %s <application_name>", os.Args[0])
		return
	}

	applicationName := applicationArgs[1]

	logger := getLogger(applicationName)
	logger.Info(fmt.Sprintf("Starting server '%s@%s'...", shared.ServerName, shared.ServerVersion), false, shared.MainComponentName)

	go func() {
		logger.Info(fmt.Sprintf("Starting pprof network listener on localhost:6060"), false, shared.MainComponentName)

		err := http.ListenAndServe("localhost:6060", nil)

		if err != nil {
			logger.Critical(fmt.Sprintf("Something went wrong while starting the pprof network Listener: '%s'", err.Error()), false, shared.MainComponentName)
			return
		}
	}()

	var wg sync.WaitGroup

	coder := coding.JsonMessageCoder{Logger: &logger}

	router := networking.NewRouter(&logger, &wg, &coder)
	router.Start()

	wg.Wait()
}
