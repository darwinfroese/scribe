package theme

func prairie() *Theme {
	return &Theme{
		Base:             prairieThemeKey,
		Background:       "#FFF0D1",
		BackgroundFocus:  "#FFE5B3",
		Text:             "#1A0B00",
		TextFocus:        "#1A0B00",
		SubText:          "#A7A7A7",
		InputBackground:  "#F0E1C2",
		Border:           "#0065AD",
		PriorityLow:      "#0065AD",
		PriorityMedium:   "#3D7F2E",
		PriorityHigh:     "#AFA53C",
		PriorityCritical: "#FF3336",
	}
}
