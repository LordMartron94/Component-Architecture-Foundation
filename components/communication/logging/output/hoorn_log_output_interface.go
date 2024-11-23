package output

import (
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/common"
)

type HoornLogOutputInterface interface {
	Output(hoornLog common.HoornLog)
}
