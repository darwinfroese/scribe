package theme

func light() *Theme {
	return &Theme{
		Base:             lightThemeKey,
		Background:       "#F6F1DE",
		BackgroundFocus:  "#BEDAD5",
		Text:             "#474747",
		TextFocus:        "#0A0A0A",
		SubText:          "#0A0A0A",
		InputBackground:  "#E8E3CF",
		Border:           "#0072BB",
		PriorityLow:      "#33B1FF",
		PriorityMedium:   "#F4AC45",
		PriorityHigh:     "#FD6035",
		PriorityCritical: "#A61C3C",
	}
}
