package formatting

import (
	"github.com/component-architecture-foundation/logging/common"
)

type HoornLogFormatterInterface interface {
	Format(hoornLog common.HoornLog) string
}
