package formatting

import (
	"bytes"
	"fmt"
	"time"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/common"
)

type HoornLogTextFormatter struct {
	longestLogLevelLength int
}

func NewHoornLogTextFormatter() *HoornLogTextFormatter {
	return &HoornLogTextFormatter{
		longestLogLevelLength: getLongestLogLevelLength(),
	}
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

func (formatter *HoornLogTextFormatter) formatLogTime(log *common.HoornLog) string {
	maxLen := 35
	var buffer bytes.Buffer
	buffer.Grow(maxLen)

	logTime := log.GetLogTime().Format(time.RFC3339Nano)
	buffer.WriteString(logTime)

	diff := maxLen - len(logTime)
	buffer.Write(bytes.Repeat([]byte{' '}, diff))

	return buffer.String()
}

func (formatter *HoornLogTextFormatter) Format(log *common.HoornLog) string {
	var buffer bytes.Buffer

	logLevel := log.GetLogLevelString()
	logTime := formatter.formatLogTime(log)

	buffer.WriteByte('[')
	buffer.WriteString(logTime)
	buffer.WriteString("] ")

	buffer.WriteString(fmt.Sprintf("%-*s", formatter.longestLogLevelLength, logLevel))

	buffer.WriteString(" : ")
	buffer.WriteString(log.GetLogMessage())

	return buffer.String()
}
