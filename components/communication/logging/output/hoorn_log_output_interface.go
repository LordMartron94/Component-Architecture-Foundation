package output

import (
	"github.com/component-architecture-foundation/logging/common"
)

type HoornLogOutputInterface interface {
	Output(hoornLog common.HoornLog)
}
