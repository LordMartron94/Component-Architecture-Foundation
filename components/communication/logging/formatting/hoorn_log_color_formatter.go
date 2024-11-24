package formatting

import (
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/common"
	"github.com/LordMartron94/Component-Architecture-Foundation/components/communication/logging/utils"
)

type Color struct {
	Text       string
	Background *string
}

type HoornLogColorFormatter struct {
	colorDict   map[common.LogLevel]Color
	colorHelper utils.ColorHelper
}

func getPointer(str string) *string {
	return &str
}

func NewHoornLogColorFormatter() *HoornLogColorFormatter {
	return &HoornLogColorFormatter{
		colorDict: map[common.LogLevel]Color{
			common.DEBUG:    {Text: "#079B00", Background: nil},
			common.INFO:     {Text: "#9B9B9B", Background: nil},
			common.WARNING:  {Text: "#FFA300", Background: nil},
			common.ERROR:    {Text: "#FF0000", Background: nil},
			common.CRITICAL: {Text: "#FFFFFF", Background: getPointer("#FF0000")},
			common.DEFAULT:  {Text: "#9B9B9B", Background: nil},
		},
	}
}

func (f HoornLogColorFormatter) Format(hoornLog *common.HoornLog) []byte {
	var logLevel = hoornLog.GetLogLevel()

	color, found := f.colorDict[logLevel]
	if !found {
		logLevel = common.DEFAULT
	}

	var textColorHex = color.Text
	var backgroundColorHex = color.Background

	var colorizedString []byte
	if backgroundColorHex != nil {
		colorizedString = f.colorHelper.ColorizeString(hoornLog.GetFormattedMessage(), textColorHex, *backgroundColorHex)
	} else {
		colorizedString = f.colorHelper.ColorizeString(hoornLog.GetFormattedMessage(), textColorHex, "")
	}

	return append(colorizedString, []byte("\x1b[0m")...)
}
