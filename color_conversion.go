package main

import (
	"image/color"
	"math"
)

func HSVtoRGBNorm(h, s, v float64) (float64, float64, float64) {
	h1 := h / 60.0
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h1, 2)-1))
	switch {
	case 0 <= h1 && h1 < 1:
		return c, x, 0
	case 1 <= h1 && h1 < 3:
		return x, c, 0
	case 2 <= h1 && h1 < 3:
		return 0, c, x
	case 3 <= h1 && h1 < 4:
		return 0, x, c
	case 4 <= h1 && h1 < 5:
		return x, 0, c
	case 5 <= h1 && h1 < 6:
		return c, 0, x
	}
	return 0, 0, 0
}

// Decorate/Wrap the HSV func to return a color.RGBA
func HSVtoRGBA(f func(h, s, v float64) (float64, float64, float64)) func(h, s, v float64) color.RGBA {
	return func(h, s, v float64) color.RGBA {
		r, g, b := f(h, 1, 1)
		r *= 255
		g *= 255
		b *= 255
		return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
	}
}
