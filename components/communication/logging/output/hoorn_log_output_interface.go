package output

import (
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/common"
)

type HoornLogOutputInterface interface {
	// Output performs operations on the HoornLog according to the specific logging implementation.
	Output(hoornLog *common.HoornLog)

	// Save performs any saving operations according to the specific logging implementation. Not all implementations make use of this.
	Save()
}
