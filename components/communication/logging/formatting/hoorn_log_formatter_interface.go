package formatting

import (
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/common"
)

type HoornLogFormatterInterface interface {
	Format(hoornLog *common.HoornLog) string
}
