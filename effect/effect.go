package effect

import (
	"image/color"

	"github.com/huderlem/contest-painting-effects/canvas"
	"github.com/huderlem/contest-painting-effects/pixelq"
)

// ApplyRedChannelGrayscale performs a grayscale effect on the canvas using
// the red color channel. A delta value is added to the red channel.
func ApplyRedChannelGrayscale(c canvas.Canvas, delta int) {
	for x := 0; x < c.Width(); x++ {
		for y := 0; y < c.Height(); y++ {
			pixel := c.At(x, y)
			if pixel.A == 255 {
				// Gets the grayscale value, based on the pixel's red channel.
				// Also adds a delta to skew lighter or darker.
				grayValue := int(pixel.R)
				grayValue += delta
				if grayValue < 0 {
					grayValue = 0
				}
				if grayValue > 31 {
					grayValue = 31
				}
				c.Set(x, y, color.RGBA{uint8(grayValue), uint8(grayValue), uint8(grayValue), pixel.A})
			}
		}
	}
}

// ApplyRedChannelGrayscaleHighlight performs a grayscale highlight effect
// on the canvas using the red color channel. Brighter colors are clamped
// according to the highlight threshold.
func ApplyRedChannelGrayscaleHighlight(c canvas.Canvas, highlight int) {
	for x := 0; x < c.Width(); x++ {
		for y := 0; y < c.Height(); y++ {
			pixel := c.At(x, y)
			if pixel.A == 255 {
				grayValue := int(pixel.R)
				if grayValue > 31-highlight {
					grayValue = 31 - highlight/2
				}
				c.Set(x, y, color.RGBA{uint8(grayValue), uint8(grayValue), uint8(grayValue), pixel.A})
			}
		}
	}
}

// ApplyGrayscale performs a grayscale effect on the canvas using a specific
// weighting of each color channel.
func ApplyGrayscale(c canvas.Canvas) {
	for x := 0; x < c.Width(); x++ {
		for y := 0; y < c.Height(); y++ {
			pixel := c.At(x, y)
			if pixel.A == 255 {
				grayValue := float32(pixel.R)*0.3 + float32(pixel.G)*0.59 + float32(pixel.B)*0.1133
				c.Set(x, y, color.RGBA{uint8(grayValue), uint8(grayValue), uint8(grayValue), pixel.A})
			}
		}
	}
}

// ApplyBlur performs a blur-like effect on the canvas. This is not a gaussian
// blur.  Instead, it only considers the two pixels above and below the pixel
// in question and attempts to reconcile their RGB differences. The result is
// more of a "smudge" than a "blur".
func ApplyBlur(c canvas.Canvas) {
	for x := 0; x < c.Width(); x++ {
		prevPixel := c.At(x, 0)
		for y := 1; y < c.Height()-1; y++ {
			pixel := c.At(x, y)
			if pixel.A == 255 {
				nextPixel := c.At(x, y+1)
				blurredPixel := pixelq.Blur(prevPixel, pixel, nextPixel)
				c.Set(x, y, blurredPixel)
				prevPixel = blurredPixel
			} else {
				c.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}
}

// ApplyPersonalityColor performs a solid-color effect on the canvas for darker colors.
// Lighter colors are turned white. The personality value determines what color to use
// for darker colors. In Pokémon Emerald, this is the lower 8 bits of the mon's
// personality value.
func ApplyPersonalityColor(c canvas.Canvas, personality uint8) {
	for x := 0; x < c.Width(); x++ {
		for y := 0; y < c.Height(); y++ {
			pixel := c.At(x, y)
			if pixel.A == 255 {
				c.Set(x, y, pixelq.PersonalityColor(pixel, personality))
			}
		}
	}
}

// ApplyBlackAndWhite converts all colors to each black or white, depending on the
// pixel's average color channel value.
func ApplyBlackAndWhite(c canvas.Canvas) {
	for x := 0; x < c.Width(); x++ {
		for y := 0; y < c.Height(); y++ {
			pixel := c.At(x, y)
			if pixel.A == 255 {
				c.Set(x, y, pixelq.BlackAndWhite(pixel))
			}
		}
	}
}

// ApplyBlackOutline performs an outline effect on the canvas. All pixels that border
// transparency are changed to black.
func ApplyBlackOutline(c canvas.Canvas) {
	for y := 0; y < c.Height(); y++ {
		c.Set(0, y, pixelq.BlackOutline(c.At(0, y), c.At(1, y)))
		for x := 1; x < c.Width()-1; x++ {
			rightOutline := pixelq.BlackOutline(c.At(x, y), c.At(x+1, y))
			c.Set(x, y, rightOutline)
			leftOutline := pixelq.BlackOutline(c.At(x, y), c.At(x-1, y))
			c.Set(x, y, leftOutline)
		}
		right := c.Width() - 1
		c.Set(right, y, pixelq.BlackOutline(c.At(right, y), c.At(right-1, y)))
	}
	for x := 0; x < c.Width(); x++ {
		c.Set(x, 0, pixelq.BlackOutline(c.At(x, 0), c.At(x, 1)))
		for y := 1; y < c.Height()-1; y++ {
			c.Set(x, y, pixelq.BlackOutline(c.At(x, y), c.At(x, y+1)))
			c.Set(x, y, pixelq.BlackOutline(c.At(x, y), c.At(x, y-1)))
		}
		bottom := c.Height() - 1
		c.Set(x, bottom, pixelq.BlackOutline(c.At(x, bottom), c.At(x, bottom-1)))
	}
}

// ApplyInvert performs a negative effect on the canvas. All pixel colors are
// inverted.
func ApplyInvert(c canvas.Canvas) {
	for y := 0; y < c.Height(); y++ {
		for x := 0; x < c.Width(); x++ {
			pixel := c.At(x, y)
			if pixel.A == 255 {
				invertedPixel := pixelq.Invert(c.At(x, y))
				c.Set(x, y, invertedPixel)
			} else {
				c.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}
}

// ApplyShimmer performs a shimmering effect on the canvas. The edges become
// very light, and it sort of looks like a mirage(?).
func ApplyShimmer(c canvas.Canvas) {
	// First, invert all of the colors.
	ApplyInvert(c)

	// Blur the pixels twice.
	for x := 0; x < c.Width(); x++ {
		prevPixel := c.At(x, 0)
		for y := 1; y < c.Height()-1; y++ {
			pixel := c.At(x, y)
			if pixel.A == 255 {
				nextPixel := c.At(x, y+1)
				blurredPixel := pixelq.BlurHard(prevPixel, pixel, nextPixel)
				c.Set(x, y, blurredPixel)
				prevPixel = blurredPixel
			} else {
				c.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}
	for x := 0; x < c.Width(); x++ {
		prevPixel := c.At(x, 0)
		for y := 1; y < c.Height()-1; y++ {
			pixel := c.At(x, y)
			if pixel.A == 255 {
				nextPixel := c.At(x, y+1)
				blurredPixel := pixelq.BlurHard(prevPixel, pixel, nextPixel)
				c.Set(x, y, blurredPixel)
				prevPixel = blurredPixel
			} else {
				c.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}

	// Finally, invert colors back to the original color space.
	// The above blur causes the outline areas to darken, which makes
	// this inversion give the effect of light outlines.
	ApplyInvert(c)
}

// ApplyBlurRight performs a right-direction motion blur effect on
// the canvas. This is not a gaussian blur.  Instead, it only considers
// pixel directly to the right of the pixel in question and attempts to
// reconcile their RGB differences. The result is more of a "smudge" than
// a "blur".
func ApplyBlurRight(c canvas.Canvas) {
	for y := 0; y < c.Height(); y++ {
		prevPixel := c.At(0, y)
		for x := 1; x < c.Width()-1; x++ {
			pixel := c.At(x, y)
			if pixel.A == 255 {
				blurredPixel := pixelq.MotionBlur(prevPixel, pixel)
				c.Set(x, y, blurredPixel)
				prevPixel = blurredPixel
			}
		}
	}
}

// ApplyBlurDown performs a down-direction motion blur effect on the
// canvas. This is not a gaussian blur.  Instead, it only considers
// pixel directly beneath the pixel in question and attempts to reconcile
// their RGB differences. The result is more of a "smudge" than a "blur".
func ApplyBlurDown(c canvas.Canvas) {
	for x := 0; x < c.Width(); x++ {
		prevPixel := c.At(x, 0)
		for y := 1; y < c.Height()-1; y++ {
			pixel := c.At(x, y)
			if pixel.A == 255 {
				blurredPixel := pixelq.MotionBlur(prevPixel, pixel)
				c.Set(x, y, blurredPixel)
				prevPixel = blurredPixel
			}
		}
	}
}

// ApplyPointillism performs a pointillism effect on the canvas, so
// it looks like the image was drawn with a individual tiny dots.
func ApplyPointillism(c canvas.Canvas) {
	for i := 0; i < 3200; i++ {
		pixelq.AddPointillismPoints(c, i)
	}
}
