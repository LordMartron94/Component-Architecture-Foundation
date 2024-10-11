package utils

import (
	"fmt"
	"strconv"
)

type ColorHelper struct{}

func (c ColorHelper) convertHexToRGB(hexColor string) (r, g, b int64) {
	if len(hexColor) < 7 { // Add length check
		return 0, 0, 0 // Or handle the error appropriately
	}

	r, _ = strconv.ParseInt(hexColor[1:3], 16, 64)
	g, _ = strconv.ParseInt(hexColor[3:5], 16, 64)
	b, _ = strconv.ParseInt(hexColor[5:7], 16, 64)
	return
}

func (c ColorHelper) ColorizeString(s, textColorHex, backgroundColorHex string) string {
	if len(s) == 0 {
		return ""
	}

	rText, gText, bText := c.convertHexToRGB(textColorHex)
	closestAnsiColorCode := fmt.Sprintf("\x1b[38;2;%d;%d;%dm", rText, gText, bText)

	if len(backgroundColorHex) != 0 {
		rBackground, gBackground, bBackground := c.convertHexToRGB(backgroundColorHex)
		closestAnsiColorCodeBackground := fmt.Sprintf("\x1b[48;2;%d;%d;%dm", rBackground, gBackground, bBackground)
		closestAnsiColorCode = closestAnsiColorCodeBackground + closestAnsiColorCode
	}

	return closestAnsiColorCode + s
}
