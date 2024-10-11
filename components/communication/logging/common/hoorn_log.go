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
	logMessage string

	// FormattedMessage is a pre-formatted string for Output. Optional.
	FormattedMessage string

	// LogSeparator is a string representing the distinction between different types of logs.
	// Different types of log outputs might use this differently.
	// Max 30 characters for proper formatting.
	LogSeparator string
}

func NewHoornLog(logTime time.Time, logLevel LogLevel, logMessage string, formattedMessage string, logSeparator string) HoornLog {
	if len(logSeparator) > 30 {
		log.Fatalln(fmt.Sprintf("Log separator exceeds maximum length of 15 characters, separator: %s", logSeparator))
	}

	return HoornLog{
		logTime:          logTime,
		logLevel:         logLevel,
		logMessage:       logMessage,
		FormattedMessage: formattedMessage,
		LogSeparator:     logSeparator,
	}
}

func (log HoornLog) GetLogLevel() LogLevel {
	return log.logLevel
}

func (log HoornLog) GetLogLevelString() string {
	return log.logLevel.StringifyLogLevel()
}

func (log HoornLog) GetLogTime() time.Time {
	return log.logTime
}

func (log HoornLog) GetLogMessage() string {
	return log.logMessage
}

func (log HoornLog) GetFormattedMessage() string {
	return log.FormattedMessage
}

func (log HoornLog) GetLogSeparator() string {
	return log.LogSeparator
}
