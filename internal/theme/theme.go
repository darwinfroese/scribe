package theme

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
)

const (
	darkThemeKey   = "dark"
	lightThemeKey  = "light"
	customThemeKey = "custom"

	fallbackColor = 0x000000
)

type Theme struct {
	Base string

	Background      string
	BackgroundFocus string

	Text      string
	TextFocus string

	SubText string

	Border string

	PriorityLow      string
	PriorityMedium   string
	PriorityHigh     string
	PriorityCritical string
}

// theme.PrimitiveBackgroundColor = tcell.NewHexColor(0xfff0d1)
// theme.BorderColor = tcell.NewHexColor(0x0065ad)
// theme.PrimaryTextColor = tcell.NewHexColor(0x1a0b00)
// theme.SecondaryTextColor = tcell.NewHexColor(0x1a0b00)
// theme.TertiaryTextColor = tcell.NewHexColor(0x1a0b00)
// theme.TitleColor = tcell.NewHexColor(0x1a0b00)
// theme.BackgroundFocus = tcell.NewHexColor(0xffe5b3)

func Load(key string) *Theme {
	switch key {
	case customThemeKey:
		return custom()
	case lightThemeKey:
		return light()
	case darkThemeKey:
		return dark()
	}

	return dark()
}

func Color(color string) tcell.Color {
	if len(color) == 7 && color[0] == '#' {
		if v, e := strconv.ParseInt(color[1:], 16, 32); e == nil {
			return tcell.NewHexColor(int32(v))
		}
	}

	return tcell.NewHexColor(fallbackColor)
}
