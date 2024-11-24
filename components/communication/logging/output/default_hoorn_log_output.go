package output

import (
	"fmt"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/common"
)

type DefaultHoornLogOutput struct {
}

func (o *DefaultHoornLogOutput) Output(log *common.HoornLog) {
	fmt.Println(fmt.Sprintf("[%-30s] %s", log.LogSeparator, log.GetFormattedMessage()))
}

func (o *DefaultHoornLogOutput) Save() {
	// Nothing, we don't need to save anything.
}
