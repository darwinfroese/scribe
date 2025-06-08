package theme

func dark() *Theme {
	return &Theme{
		Base:             darkThemeKey,
		Background:       "#222831",
		BackgroundFocus:  "#393E46",
		Text:             "#DFD0B8",
		TextFocus:        "#948979",
		SubText:          "#948979",
		InputBackground:  "#333A45",
		Border:           "#123458",
		PriorityLow:      "#309898",
		PriorityMedium:   "#FF9F00",
		PriorityHigh:     "#F4631E",
		PriorityCritical: "#CB0404",
	}
}
