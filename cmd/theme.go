package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type myTheme struct{}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNamePrimary {
		if variant == theme.VariantDark {
			return color.RGBA{0, 200, 255, 255}
		}
		return color.Black
	}

	return theme.DefaultTheme().Color(name, variant)
}
func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {

	return theme.DefaultTheme().Icon(name)
}
func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
