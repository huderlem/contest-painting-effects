package pixelq

import (
	"image/color"

	"github.com/huderlem/contest-painting-effects/canvas"
)

// Invert returns the pixel color. Also called the "negative".
func Invert(pixel color.RGBA) color.RGBA {
	return color.RGBA{
		R: 31 - pixel.R,
		G: 31 - pixel.G,
		B: 31 - pixel.B,
		A: pixel.A,
	}
}

// Blur returns a blurred version of the given 3 pixels.
func Blur(prevPixel, curPixel, nextPixel color.RGBA) color.RGBA {
	if prevPixel == curPixel && nextPixel == curPixel {
		return curPixel
	}

	prevAvg := (int(prevPixel.R) + int(prevPixel.G) + int(prevPixel.B)) / 3
	curAvg := (int(curPixel.R) + int(curPixel.G) + int(curPixel.B)) / 3
	nextAvg := (int(nextPixel.R) + int(nextPixel.G) + int(nextPixel.B)) / 3

	if prevAvg == curAvg && nextAvg == curAvg {
		return curPixel
	}

	prevDiff := int(curAvg - prevAvg)
	if curAvg < prevAvg {
		prevDiff = int(prevAvg - curAvg)
	}

	nextDiff := int(curAvg - nextAvg)
	if curAvg < nextAvg {
		nextDiff = int(nextAvg - curAvg)
	}

	diff := nextDiff
	if prevDiff > nextDiff {
		diff = prevDiff
	}

	factor := 31 - diff/2
	return color.RGBA{
		R: uint8((int(curPixel.R) * factor) / 31),
		G: uint8((int(curPixel.G) * factor) / 31),
		B: uint8((int(curPixel.B) * factor) / 31),
		A: curPixel.A,
	}
}

// BlurHard returns a harder-blurred version of the given 3 pixels.
func BlurHard(prevPixel, curPixel, nextPixel color.RGBA) color.RGBA {
	if prevPixel == curPixel && nextPixel == curPixel {
		return curPixel
	}

	prevAvg := (int(prevPixel.R) + int(prevPixel.G) + int(prevPixel.B)) / 3
	curAvg := (int(curPixel.R) + int(curPixel.G) + int(curPixel.B)) / 3
	nextAvg := (int(nextPixel.R) + int(nextPixel.G) + int(nextPixel.B)) / 3

	if prevAvg == curAvg && nextAvg == curAvg {
		return curPixel
	}

	prevDiff := int(curAvg - prevAvg)
	if curAvg < prevAvg {
		prevDiff = int(prevAvg - curAvg)
	}

	nextDiff := int(curAvg - nextAvg)
	if curAvg < nextAvg {
		nextDiff = int(nextAvg - curAvg)
	}

	diff := nextDiff
	if prevDiff > nextDiff {
		diff = prevDiff
	}

	factor := 31 - diff
	return color.RGBA{
		R: uint8((int(curPixel.R) * factor) / 31),
		G: uint8((int(curPixel.G) * factor) / 31),
		B: uint8((int(curPixel.B) * factor) / 31),
		A: curPixel.A,
	}
}

// MotionBlur returns a blurred version of the given 2 pixels.
func MotionBlur(prevPixel, curPixel color.RGBA) color.RGBA {
	if prevPixel == curPixel {
		return curPixel
	}

	// Don't blur light colors.
	if prevPixel.R > 25 && prevPixel.G > 25 && prevPixel.B > 25 {
		return curPixel
	}
	if curPixel.R > 25 && curPixel.G > 25 && curPixel.B > 25 {
		return curPixel
	}

	// Find the largest diff of any of the color channels.
	redDiff := int(prevPixel.R) - int(curPixel.R)
	if redDiff < 0 {
		redDiff = -redDiff
	}
	greenDiff := int(prevPixel.G) - int(curPixel.G)
	if greenDiff < 0 {
		greenDiff = -greenDiff
	}
	blueDiff := int(prevPixel.B) - int(curPixel.B)
	if blueDiff < 0 {
		blueDiff = -blueDiff
	}
	diff := redDiff
	if greenDiff > diff {
		diff = greenDiff
	}
	if blueDiff > diff {
		diff = blueDiff
	}

	factor := 31 - diff/2
	return color.RGBA{
		R: uint8((int(curPixel.R) * factor) / 31),
		G: uint8((int(curPixel.G) * factor) / 31),
		B: uint8((int(curPixel.B) * factor) / 31),
		A: curPixel.A,
	}
}

// BlackOutline returns a border-colored pixel if pixelA is a border pixel.
func BlackOutline(pixelA, pixelB color.RGBA) color.RGBA {
	if pixelA.R == 0 && pixelA.G == 0 && pixelA.B == 0 && pixelA.A == 255 {
		return pixelA
	}
	if pixelA.A != 255 {
		return color.RGBA{0, 0, 0, 0}
	}
	if pixelB.A != 255 {
		return color.RGBA{0, 0, 0, 255}
	}
	return pixelA
}

// PersonalityColor returns a color determined by the personality value of a pokemon. (only uses lower 8 bytes, which is why)
// personality is a uint8. Returns white if the pixel is light.
func PersonalityColor(pixel color.RGBA, personality uint8) color.RGBA {
	if pixel.R < 17 && pixel.G < 17 && pixel.B < 17 {
		return getColorFromPersonality(personality)
	}
	return color.RGBA{31, 31, 31, 255}
}

func getColorFromPersonality(personality uint8) color.RGBA {
	var red, green, blue uint8 = 0, 0, 0
	strength := (personality / 6) % 3

	switch personality % 6 {
	case 0:
		// Teal color
		green = 21 - strength
		blue = green
		red = 0
	case 1:
		// Yellow color
		blue = 0
		red = 21 - strength
		green = red
	case 2:
		// Purple color
		blue = 21 - strength
		green = 0
		red = blue
	case 3:
		// Red color
		blue = 0
		green = 0
		red = 23 - strength
	case 4:
		// Blue color
		blue = 23 - strength
		green = 0
		red = 0
	case 5:
		// Green color
		blue = 0
		green = 23 - strength
		red = 0
	}
	return color.RGBA{red, green, blue, 255}
}

// BlackAndWhite converts a pixel to black or white.
func BlackAndWhite(pixel color.RGBA) color.RGBA {
	if pixel.R < 17 && pixel.G < 17 && pixel.B < 17 {
		return color.RGBA{0, 0, 0, 255}
	}
	return color.RGBA{31, 31, 31, 255}
}

type pointillismPoint struct {
	column int
	row    int
	delta  uint8
}

// AddPointillismPoints splats dots onto the canvas to give
// a pointillism effect.
func AddPointillismPoints(c canvas.Canvas, point int) {
	for cx := 0; cx < c.Width()/64; cx++ {
		for cy := 0; cy < c.Height()/64; cy++ {
			index := point * 3
			points := make([]pointillismPoint, 6)
			points[0].column = int(pointillism[index]) + cx*64
			points[0].row = int(pointillism[index+1]) + cy*64
			points[0].delta = (pointillism[index+2] >> 3) & 7

			colorType := (pointillism[index+2] >> 1) & 3
			offsetDownLeft := pointillism[index+2] & 1
			for i := uint8(1); i < points[0].delta; i++ {
				if offsetDownLeft == 0 {
					points[i].column = points[0].column - int(i)
					points[i].row = points[0].row + int(i)
				} else {
					points[i].column = points[0].column + 1
					points[i].row = points[0].row - 1
				}
				if int(points[i].column) > c.Width()-1 || int(points[i].row) > c.Height()-1 {
					points[0].delta = i - 1
					break
				}

				points[i].delta = points[0].delta - i
			}

			for i := uint8(0); i < points[0].delta; i++ {
				x, y := int(points[i].column), int(points[i].row)
				pixel := c.At(x, y)
				if pixel.A == 255 {
					red := pixel.R
					green := pixel.G
					blue := pixel.B
					switch colorType {
					case 0, 1:
						switch ((pointillism[index+2] >> 3) & 7) % 3 {
						case 0:
							if red >= points[i].delta {
								red -= points[i].delta
							} else {
								red = 0
							}
						case 1:
							if green >= points[i].delta {
								green -= points[i].delta
							} else {
								green = 0
							}
						case 2:
							if blue >= points[i].delta {
								blue -= points[i].delta
							} else {
								blue = 0
							}
						}
					case 2, 3:
						red += points[i].delta
						green += points[i].delta
						blue += points[i].delta
						if red > 31 {
							red = 31
						}
						if green > 31 {
							green = 31
						}
						if blue > 31 {
							blue = 31
						}
					}

					c.Set(x, y, color.RGBA{red, green, blue, pixel.A})
				}
			}
		}
	}
}
