package formatting

import (
	"fmt"
	"strings"
	"time"

	"github.com/component-architecture-foundation/logging/common"
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

func (formatter HoornLogTextFormatter) Format(log common.HoornLog) string {
	var logLevel string = log.GetLogLevelString()

	var formattedMessage string = "[" + log.GetLogTime().Format(time.RFC3339Nano) + "] " + logLevel + " : " + log.GetLogMessage()
	formattedMessage = strings.Replace(formattedMessage, logLevel, fmt.Sprintf("%-*s", getLongestLogLevelLength(), logLevel), -1)

	return formattedMessage
}
