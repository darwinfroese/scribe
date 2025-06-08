package theme

func custom() *Theme {
	return &Theme{
		Base:             customThemeKey,
		Background:       "",
		BackgroundFocus:  "",
		Text:             "",
		TextFocus:        "",
		Border:           "",
		PriorityLow:      "",
		PriorityMedium:   "",
		PriorityHigh:     "",
		PriorityCritical: "",
	}
}
