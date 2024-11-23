package factory

import (
	"time"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/common"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/formatting"
)

type HoornLogFactory struct {
}

func (factory HoornLogFactory) CreateHoornLog(level common.LogLevel, message string, logSeparator string) common.HoornLog {
	var currentTime time.Time = time.Now()

	var formatters []formatting.HoornLogFormatterInterface = []formatting.HoornLogFormatterInterface{
		formatting.HoornLogTextFormatter{},
		formatting.NewHoornLogColorFormatter(),
	}

	var hoornLog = common.NewHoornLog(
		currentTime,
		level,
		message,
		message,
		logSeparator,
	)

	for _, formatter := range formatters {
		hoornLog.FormattedMessage = formatter.Format(hoornLog)
	}

	return hoornLog
}
