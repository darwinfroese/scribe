package theme

import (
	"reflect"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

const (
	darkThemeKey    = "dark"
	lightThemeKey   = "light"
	prairieThemeKey = "prairie"

	customThemeKey = "custom"

	fallbackColor = 0x000000
)

type Theme struct {
	Base string

	Background      string
	BackgroundFocus string

	Text      string
	TextFocus string

	SubText         string
	InputBackground string

	Border string

	PriorityLow      string
	PriorityMedium   string
	PriorityHigh     string
	PriorityCritical string
}

func Load(theme *Theme) *Theme {
	base := getBaseTheme(theme.Base)

	override(base, theme)

	return base
}

func Color(color string) tcell.Color {
	if len(color) == 7 && color[0] == '#' {
		if v, e := strconv.ParseInt(color[1:], 16, 32); e == nil {
			return tcell.NewHexColor(int32(v))
		}
	}

	return tcell.NewHexColor(fallbackColor)
}

func getBaseTheme(key string) *Theme {
	switch key {
	case customThemeKey:
		return custom()
	case lightThemeKey:
		return light()
	case darkThemeKey:
		return dark()
	case prairieThemeKey:
		return prairie()
	}

	return dark()
}

func override(base, theme *Theme) {
	themeElem := reflect.ValueOf(theme).Elem()
	baseElem := reflect.ValueOf(base).Elem()

	for i := range themeElem.NumField() {
		themeField := themeElem.Field(i)

		if themeField.IsZero() {
			continue
		}

		baseField := baseElem.Field(i)
		if baseField.CanSet() {
			baseField.Set(themeField)
		}
	}
}
