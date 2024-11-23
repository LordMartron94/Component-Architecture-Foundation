package logging

import (
	"sync"
	"time"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/common"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/factory"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/output"
)

type HoornLogger struct {
	// outputs is a list of HoornLogOutputInterface objects that will handle the logging.
	outputs []output.HoornLogOutputInterface
	// minLevel is the minimum log level required for a message to be logged.
	minLevel common.LogLevel

	// hoornLogFactory is a HoornLogFactoryInterface object that will be used to create new HoornLog objects.
	hoornLogFactory factory.HoornLogFactory

	// shutdownSignal is a receive-only channel that will be closed when the logger should stop logging.
	shutdownSignal chan struct{}

	// waitGroup is a synchronization primitive that will be used to wait for all goroutines to finish.
	waitGroup *sync.WaitGroup
}

func NewHoornLogger(minLevel common.LogLevel, shutdownSignal chan struct{}, wg *sync.WaitGroup, outputs ...output.HoornLogOutputInterface) HoornLogger {
	if len(outputs) == 0 {
		outputs = []output.HoornLogOutputInterface{&output.DefaultHoornLogOutput{}}
	}

	logger := HoornLogger{
		minLevel:        minLevel,
		outputs:         outputs,
		hoornLogFactory: factory.HoornLogFactory{},
		shutdownSignal:  shutdownSignal,
		waitGroup:       wg,
	}

	go logger.ListenForShutdown()

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				logger.save()
			case <-shutdownSignal:
				return
			}
		}
	}()

	return logger
}

func (hL *HoornLogger) canOutput(level common.LogLevel) bool {
	return level >= hL.minLevel
}

func (hL *HoornLogger) log(level common.LogLevel, message string, forceShow bool, separator string) {
	if !hL.canOutput(level) && !forceShow {
		return
	}

	var hoornLog = hL.hoornLogFactory.CreateHoornLog(level, message, separator)

	for _, outputMethod := range hL.outputs {
		outputMethod.Output(hoornLog)
	}
}

func (hL *HoornLogger) ListenForShutdown() {
	hL.waitGroup.Add(1)

	<-hL.shutdownSignal

	hL.save()

	hL.waitGroup.Done()
}

func (hL *HoornLogger) save() {
	for _, outputMethod := range hL.outputs {
		outputMethod.Save()
	}
}

func (hL *HoornLogger) SetMinLevel(level common.LogLevel) {
	hL.minLevel = level
}

func (hL *HoornLogger) Debug(message string, forceShow bool, separator string) {
	hL.log(common.DEBUG, message, forceShow, separator)
}

func (hL *HoornLogger) Info(message string, forceShow bool, separator string) {
	hL.log(common.INFO, message, forceShow, separator)
}

func (hL *HoornLogger) Warn(message string, forceShow bool, separator string) {
	hL.log(common.WARNING, message, forceShow, separator)
}

func (hL *HoornLogger) Error(message string, forceShow bool, separator string) {
	hL.log(common.ERROR, message, forceShow, separator)
}

func (hL *HoornLogger) Critical(message string, forceShow bool, separator string) {
	hL.log(common.CRITICAL, message, forceShow, separator)
}
