package paletteq

import (
	"image/color"

	"github.com/huderlem/contest-painting-effects/pixelq"

	"github.com/huderlem/contest-painting-effects/canvas"
)

// ApplyStandardQuantization generates a quantized palette for the Canvas pixels, and
// assigns canvas pixels to each color in the quantized palette.
func ApplyStandardQuantization(c canvas.Canvas, maxColors int) []color.RGBA {
	palette := make([]color.RGBA, maxColors)
	for i := 0; i < maxColors-1; i++ {
		palette[i] = color.RGBA{0, 0, 0, 0}
	}
	palette[maxColors-1] = color.RGBA{15, 15, 15, 255}
	for y := 0; y < c.Height(); y++ {
		for x := 0; x < c.Width(); x++ {
			pixel := c.At(x, y)
			if pixel.A != 255 {
				c.SetColorIndex(x, y, 0)
			} else {
				quantizedPixel := quantizePixelStandard(pixel)
				success := false
				for curIndex := 1; curIndex < maxColors-1; curIndex++ {
					curColor := palette[curIndex]
					if curColor.R == 0 && curColor.G == 0 && curColor.B == 0 && curColor.A == 0 {
						// The quantized color does not match any existing color in the
						// palette, so we add it to the palette.
						// This if block seems pointless because the below while loop handles
						// this same logic.
						palette[curIndex] = quantizedPixel
						c.SetColorIndex(x, y, curIndex)
						success = true
						break
					} else if curColor.R == quantizedPixel.R && curColor.G == quantizedPixel.G &&
						curColor.B == quantizedPixel.B && curColor.A == quantizedPixel.A {
						// The quantized color matches this existing color in the
						// palette, so we use this existing color for the pixel.
						c.SetColorIndex(x, y, curIndex)
						success = true
						break
					}
				}
				if !success {
					// The entire palette's colors are already in use, which means
					// the base image has too many colors to handle. This error is handled
					// by marking such pixels as gray color.
					c.SetColorIndex(x, y, maxColors-1)
				}
			}
		}
	}
	return palette
}

func quantizePixelStandard(pixel color.RGBA) color.RGBA {
	r, g, b := pixel.R, pixel.G, pixel.B

	// Quantize color channels to muliples of 4, rounding up.
	if r&3 != 0 {
		r = (r & 0x1c) + 4
	}
	if g&3 != 0 {
		g = (g & 0x1c) + 4
	}
	if b&3 != 0 {
		b = (b & 0x1c) + 4
	}

	// Clamp channels to [6, 30].
	if r < 6 {
		r = 6
	}
	if r > 30 {
		r = 30
	}
	if g < 6 {
		g = 6
	}
	if g > 30 {
		g = 30
	}
	if b < 6 {
		b = 6
	}
	if b > 30 {
		b = 30
	}
	return color.RGBA{r, g, b, pixel.A}
}

// ApplyPrimaryColorsQuantization generates a quantized palette for the Canvas pixels, which
// is basd on a preset list of bright primary colors.
func ApplyPrimaryColorsQuantization(c canvas.Canvas) []color.RGBA {
	palette := make([]color.RGBA, 16)
	palette[0] = color.RGBA{0, 0, 0, 0}
	palette[1] = color.RGBA{R: 6, G: 6, B: 6, A: 255}
	palette[2] = color.RGBA{R: 29, G: 29, B: 29, A: 255}
	palette[3] = color.RGBA{R: 11, G: 11, B: 11, A: 255}
	palette[4] = color.RGBA{R: 29, G: 6, B: 6, A: 255}
	palette[5] = color.RGBA{R: 6, G: 29, B: 6, A: 255}
	palette[6] = color.RGBA{R: 6, G: 6, B: 29, A: 255}
	palette[7] = color.RGBA{R: 29, G: 29, B: 6, A: 255}
	palette[8] = color.RGBA{R: 29, G: 6, B: 29, A: 255}
	palette[9] = color.RGBA{R: 6, G: 29, B: 29, A: 255}
	palette[10] = color.RGBA{R: 29, G: 11, B: 6, A: 255}
	palette[11] = color.RGBA{R: 11, G: 29, B: 6, A: 255}
	palette[12] = color.RGBA{R: 6, G: 11, B: 29, A: 255}
	palette[13] = color.RGBA{R: 29, G: 6, B: 11, A: 255}
	palette[14] = color.RGBA{R: 6, G: 29, B: 11, A: 255}
	palette[15] = color.RGBA{R: 11, G: 6, B: 29, A: 255}

	for x := 0; x < c.Width(); x++ {
		for y := 0; y < c.Height(); y++ {
			pixel := c.At(x, y)
			if pixel.A != 255 {
				c.SetColorIndex(x, y, 0)
			} else {
				c.SetColorIndex(x, y, quantizePixelPrimaryColorsIndex(pixel))
			}
		}
	}

	return palette
}

func quantizePixelPrimaryColorsIndex(pixel color.RGBA) int {
	if pixel.R < 12 && pixel.G < 11 && pixel.B < 11 {
		return 1
	}
	if pixel.R > 19 && pixel.G > 19 && pixel.B > 19 {
		return 2
	}

	if pixel.R > 19 {
		if pixel.G > 19 {
			if pixel.B > 14 {
				return 2
			}
			return 7
		} else if pixel.B > 19 {
			if pixel.G > 14 {
				return 2
			}
			return 8
		}
	}

	if pixel.G > 19 && pixel.B > 19 {
		if pixel.R > 14 {
			return 2
		}
		return 9
	}

	if pixel.R > 19 {
		if pixel.G > 11 {
			if pixel.B > 11 {
				if pixel.G < pixel.B {
					return 8
				}
				return 7
			}
			return 10
		} else if pixel.B > 11 {
			return 13
		} else {
			return 4
		}
	}

	if pixel.G > 19 {
		if pixel.R > 11 {
			if pixel.B > 11 {
				if pixel.R < pixel.B {
					return 9
				}
				return 7
			}
			return 11
		}
		if pixel.B > 11 {
			return 14
		}
		return 5
	}

	if pixel.B > 19 {
		if pixel.R > 11 {
			if pixel.G > 11 {
				if pixel.R < pixel.G {
					return 9
				}
				return 8
			}
		} else if pixel.G > 11 {
			return 12
		}

		if pixel.B > 11 {
			return 15
		}
		return 6
	}

	return 3
}

// ApplyGrayscaleQuantization generates a quantized palette for grayscale
// (32) colors.
func ApplyGrayscaleQuantization(c canvas.Canvas) []color.RGBA {
	palette := make([]color.RGBA, 33)
	palette[0] = color.RGBA{0, 0, 0, 0}
	for i := uint8(0); i < 32; i++ {
		palette[i+1] = color.RGBA{i, i, i, 255}
	}

	for x := 0; x < c.Width(); x++ {
		for y := 0; y < c.Height(); y++ {
			pixel := c.At(x, y)
			if pixel.A != 255 {
				c.SetColorIndex(x, y, 0)
			} else {
				c.SetColorIndex(x, y, quantizePixelGrayscale(pixel))
			}
		}
	}

	return palette
}

func quantizePixelGrayscale(pixel color.RGBA) int {
	avg := (int(pixel.R) + int(pixel.G) + int(pixel.B)) / 3
	return avg + 1
}

// ApplyGrayscaleSmallQuantization generates a quantized palette for grayscale
// (16) colors.
func ApplyGrayscaleSmallQuantization(c canvas.Canvas) []color.RGBA {
	palette := make([]color.RGBA, 16)
	palette[0] = color.RGBA{0, 0, 0, 0}
	palette[1] = color.RGBA{0, 0, 0, 255}
	for i := uint8(0); i < 14; i++ {
		grayValue := 2 * (i + 2)
		palette[i+2] = color.RGBA{grayValue, grayValue, grayValue, 255}
	}

	for x := 0; x < c.Width(); x++ {
		for y := 0; y < c.Height(); y++ {
			pixel := c.At(x, y)
			if pixel.A != 255 {
				c.SetColorIndex(x, y, 0)
			} else {
				c.SetColorIndex(x, y, quantizePixelGrayscaleSmall(pixel))
			}
		}
	}

	return palette
}

func quantizePixelGrayscaleSmall(pixel color.RGBA) int {
	avg := uint8((int(pixel.R) + int(pixel.G) + int(pixel.B)) / 3)
	avg = avg & 0x1E
	if avg == 0 {
		return 1
	}
	return int(avg) / 2
}

// ApplyBlackAndWhiteQuantization generates a quantized palette for black
// and white colors.
func ApplyBlackAndWhiteQuantization(c canvas.Canvas) []color.RGBA {
	palette := make([]color.RGBA, 3)
	palette[0] = color.RGBA{0, 0, 0, 0}
	palette[1] = color.RGBA{0, 0, 0, 255}
	palette[2] = color.RGBA{31, 31, 31, 255}

	for x := 0; x < c.Width(); x++ {
		for y := 0; y < c.Height(); y++ {
			pixel := c.At(x, y)
			if pixel.A != 255 {
				c.SetColorIndex(x, y, 0)
			} else {
				qp := pixelq.BlackAndWhite(pixel)
				if qp.R == 0 && qp.G == 0 && qp.B == 0 {
					c.SetColorIndex(x, y, 1)
				} else {
					c.SetColorIndex(x, y, 2)
				}
			}
		}
	}

	return palette
}
