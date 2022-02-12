package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"github.com/dez11de/cryptodb"
)

type myTheme struct{}

var _ fyne.Theme = (*myTheme)(nil)

// TODO: remove this function and define values directly
func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	return
}

// colors based on https://github.com/EdenEast/nightfox.nvim/blob/main/extra/nightfox/nightfox_alacritty.yml
// TODO: ok for now, but could use some optimization
func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == theme.ColorNameBackground {
		rgbacolor, _ := ParseHexColor("#192330")
		return rgbacolor
	}

	if name == theme.ColorNameForeground {
		rgbacolor, _ := ParseHexColor("#cdcecf")
		return rgbacolor
	}

	return theme.DefaultTheme().Color(name, variant)
}

func DirectionColor(d cryptodb.Direction) color.Color {
	switch d {
	case cryptodb.DirectionLong:
		rgbacolor, _ := ParseHexColor("#81b29a")
		return rgbacolor
	case cryptodb.DirectionShort:
		rgbacolor, _ := ParseHexColor("#c94f6d")
		return rgbacolor
	}
	return theme.PrimaryColor()
}

func StatusColor(s cryptodb.Status) color.Color {
    // TODO: define colors for all statussus
	switch s {
	case cryptodb.StatusPlanned:
		rgbacolor, _ := ParseHexColor("#84cee4")
		return rgbacolor
	case cryptodb.StatusOrdered:
		rgbacolor, _ := ParseHexColor("#81b29a")
		return rgbacolor
	case cryptodb.StatusCancelled:
		rgbacolor, _ := ParseHexColor("#dbc074")
		return rgbacolor
	}
	return theme.PrimaryColor()
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
