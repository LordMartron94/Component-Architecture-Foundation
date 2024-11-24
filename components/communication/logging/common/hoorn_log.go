package common

import (
	"fmt"
	"log"
	"time"
)

type HoornLog struct {
	// logTime is the timestamp of the log message.
	logTime time.Time

	// logLevel is the severity level of the log message.
	logLevel LogLevel

	// logMessage is the raw log message.
	logMessage []byte

	// FormattedMessage is a pre-formatted message for Output. Optional.
	FormattedMessage []byte

	// LogSeparator is a string representing the distinction between different types of logs.
	// Different types of log outputs might use this differently.
	// Max 30 characters for proper formatting.
	LogSeparator []byte
}

func NewHoornLog(logTime time.Time, logLevel LogLevel, logMessage string, formattedMessage string, logSeparator string) *HoornLog {
	if len(logSeparator) > 30 {
		log.Fatalln(fmt.Sprintf("Log separator exceeds maximum length of 15 characters, separator: %s", logSeparator))
	}

	logMessageBytes := []byte(logMessage)
	formattedMessageBytes := []byte(formattedMessage)
	logSeparatorBytes := []byte(logSeparator)

	return &HoornLog{
		logTime:          logTime,
		logLevel:         logLevel,
		logMessage:       logMessageBytes,
		FormattedMessage: formattedMessageBytes,
		LogSeparator:     logSeparatorBytes,
	}
}

func (log *HoornLog) GetLogLevel() LogLevel {
	return log.logLevel
}

func (log *HoornLog) GetLogLevelString() string {
	return log.logLevel.StringifyLogLevel()
}

func (log *HoornLog) GetLogTime() time.Time {
	return log.logTime
}

func (log *HoornLog) GetLogMessage() []byte {
	return log.logMessage
}

func (log *HoornLog) GetFormattedMessage() []byte {
	return log.FormattedMessage
}

func (log *HoornLog) GetLogSeparator() []byte {
	return log.LogSeparator
}
