package output

import (
	"fmt"

	"github.com/component-architecture-foundation/logging/common"
)

type DefaultHoornLogOutput struct {
}

func (o DefaultHoornLogOutput) Output(log common.HoornLog) {
	fmt.Println(fmt.Sprintf("[%-30s] %s", log.LogSeparator, log.GetFormattedMessage()))
}
