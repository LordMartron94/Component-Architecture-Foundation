package _internal

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/common"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/output"
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

func GetLogger(applicationName string, shutdownSignal chan struct{}, wg *sync.WaitGroup) logging.HoornLogger {
	logDir := getLogDir(applicationName) + "\\Communication_Layer\\"

	return logging.NewHoornLogger(
		common.DEBUG,
		shutdownSignal,
		wg,
		&output.DefaultHoornLogOutput{},
		output.NewFileHoornLogOutput(
			logDir,
			5,
			true,
		))
}
