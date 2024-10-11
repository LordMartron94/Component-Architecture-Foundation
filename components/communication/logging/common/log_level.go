package common

// LogLevel is an "enum" representing the different types of logs.
type LogLevel int

const (
	// DEBUG represents a debug log.
	DEBUG LogLevel = iota
	// INFO represents an informational log.
	INFO
	// DEFAULT represents the default log level.
	DEFAULT
	// WARNING represents a warning log.
	WARNING
	// ERROR represents an error log.
	ERROR
	// CRITICAL represents a critical error log.
	CRITICAL
)

func GetAllLogLevels() []LogLevel {
	return []LogLevel{DEBUG, INFO, DEFAULT, WARNING, ERROR, CRITICAL}
}

func (level LogLevel) StringifyLogLevel() string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case DEFAULT:
		return "DEFAULT"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case CRITICAL:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}
