package formatting

import (
	"bytes"
	"fmt"

	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/common"
)

type HoornLogCombinedTextFormatter struct {
	textFormatter HoornLogTextFormatter
}

func NewHoornLogCombinedTextFormatter(textFormatter HoornLogTextFormatter) *HoornLogCombinedTextFormatter {
	return &HoornLogCombinedTextFormatter{
		textFormatter: textFormatter,
	}
}

func (formatter *HoornLogCombinedTextFormatter) Format(hoornLog *common.HoornLog) []byte {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("[%-30s] ", hoornLog.LogSeparator))
	formattedText := formatter.textFormatter.Format(hoornLog)

	for _, fmtByte := range formattedText {
		buffer.WriteByte(fmtByte)
	}

	return buffer.Bytes()
}
