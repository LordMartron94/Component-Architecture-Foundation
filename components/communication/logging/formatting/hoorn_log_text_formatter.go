package formatting

import (
	"fmt"
	"strings"
	"time"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/common"
)

type HoornLogTextFormatter struct {
}

func getLongestLogLevelLength() int {
	var logLevels []common.LogLevel = common.GetAllLogLevels()

	var longestLogLevelLength int = 0
	for _, logLevel := range logLevels {
		if len(logLevel.StringifyLogLevel()) > longestLogLevelLength {
			longestLogLevelLength = len(logLevel.StringifyLogLevel())
		}
	}

	return longestLogLevelLength
}

func (formatter HoornLogTextFormatter) Format(log *common.HoornLog) string {
	var logLevel string = log.GetLogLevelString()

	logTime := log.GetLogTime().Format(time.RFC3339Nano)

	var formattedMessage string = "[" + fmt.Sprintf("%-*s", 35, logTime) + "] " + logLevel + " : " + log.GetLogMessage()
	formattedMessage = strings.Replace(formattedMessage, logLevel, fmt.Sprintf("%-*s", getLongestLogLevelLength(), logLevel), -1)

	return formattedMessage
}
