package formatting

import (
	"bytes"
	"time"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/common"
)

type HoornLogTextFormatter struct {
	longestLogLevelLength int
}

func NewHoornLogTextFormatter() *HoornLogTextFormatter {
	return &HoornLogTextFormatter{
		longestLogLevelLength: common.GetLongestLogLevelLength(),
	}
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

func (formatter *HoornLogTextFormatter) formatLogLevel(logLevel string) string {
	maxLen := formatter.longestLogLevelLength
	var buffer bytes.Buffer
	buffer.Grow(maxLen)

	buffer.WriteString(logLevel)
	diff := maxLen - len(logLevel)

	buffer.Write(bytes.Repeat([]byte{' '}, diff))

	return buffer.String()
}

func (formatter *HoornLogTextFormatter) Format(log *common.HoornLog) []byte {
	var buffer bytes.Buffer

	logLevel := log.GetLogLevelString()
	logTime := formatter.formatLogTime(log)

	buffer.WriteByte('[')
	buffer.WriteString(logTime)
	buffer.WriteString("] ")

	buffer.WriteByte('[')
	buffer.WriteString(formatter.formatLogLevel(logLevel))
	buffer.WriteString("] ")

	buffer.WriteString(" : ")

	message := log.GetLogMessage()
	for _, msgByte := range message {
		buffer.WriteByte(msgByte)
	}

	return buffer.Bytes()
}
