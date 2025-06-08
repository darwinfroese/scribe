package theme

func custom() *Theme {
	return &Theme{
		Base:             customThemeKey,
		Background:       "#ffffff",
		BackgroundFocus:  "#cccccc",
		Text:             "#000000",
		TextFocus:        "#000000",
		Border:           "#000000",
		PriorityLow:      "#000000",
		PriorityMedium:   "#000000",
		PriorityHigh:     "#000000",
		PriorityCritical: "#000000",
	}
}
