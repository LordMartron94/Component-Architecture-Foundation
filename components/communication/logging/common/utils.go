package common

func GetLongestLogLevelLength() int {
	var logLevels []LogLevel = GetAllLogLevels()

	var longestLogLevelLength int = 0
	for _, logLevel := range logLevels {
		if len(logLevel.StringifyLogLevel()) > longestLogLevelLength {
			longestLogLevelLength = len(logLevel.StringifyLogLevel())
		}
	}

	return longestLogLevelLength
}
