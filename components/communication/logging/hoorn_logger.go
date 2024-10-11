package logging

import (
	"github.com/component-architecture-foundation/logging/common"
	"github.com/component-architecture-foundation/logging/factory"
	"github.com/component-architecture-foundation/logging/output"
)

type HoornLogger struct {
	// outputs is a list of HoornLogOutputInterface objects that will handle the logging.
	outputs []output.HoornLogOutputInterface
	// minLevel is the minimum log level required for a message to be logged.
	minLevel common.LogLevel

	// hoornLogFactory is a HoornLogFactoryInterface object that will be used to create new HoornLog objects.
	hoornLogFactory factory.HoornLogFactory
}

func NewHoornLogger(minLevel common.LogLevel, outputs ...output.HoornLogOutputInterface) HoornLogger {
	if len(outputs) == 0 {
		outputs = []output.HoornLogOutputInterface{output.DefaultHoornLogOutput{}}
	}

	return HoornLogger{
		minLevel:        minLevel,
		outputs:         outputs,
		hoornLogFactory: factory.HoornLogFactory{},
	}
}

func (hL HoornLogger) canOutput(level common.LogLevel) bool {
	return level >= hL.minLevel
}

func (hL HoornLogger) log(level common.LogLevel, message string, forceShow bool, separator string) {
	if !hL.canOutput(level) && !forceShow {
		return
	}

	var hoornLog = hL.hoornLogFactory.CreateHoornLog(level, message, separator)

	for _, outputMethod := range hL.outputs {
		outputMethod.Output(hoornLog)
	}
}

func (hL HoornLogger) SetMinLevel(level common.LogLevel) {
	hL.minLevel = level
}

func (hL HoornLogger) Debug(message string, forceShow bool, separator string) {
	hL.log(common.DEBUG, message, forceShow, separator)
}

func (hL HoornLogger) Info(message string, forceShow bool, separator string) {
	hL.log(common.INFO, message, forceShow, separator)
}

func (hL HoornLogger) Warn(message string, forceShow bool, separator string) {
	hL.log(common.WARNING, message, forceShow, separator)
}

func (hL HoornLogger) Error(message string, forceShow bool, separator string) {
	hL.log(common.ERROR, message, forceShow, separator)
}

func (hL HoornLogger) Critical(message string, forceShow bool, separator string) {
	hL.log(common.CRITICAL, message, forceShow, separator)
}
