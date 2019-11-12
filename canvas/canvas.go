package canvas

import (
	"image"
	"image/color"
)

// Canvas is a representation of the image processing canvas.
// Holds an internal state of pixels.
type Canvas struct {
	pixels       []color.RGBA
	pixelIndexes []int
	width        int
	height       int
}

// FromImage builds a new Canvas from a given image.
func FromImage(imageData image.Image) Canvas {
	bounds := imageData.Bounds()
	c := New(bounds.Max.X, bounds.Max.Y)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			r, g, b, a := imageData.At(x, y).RGBA()
			convertedRed := uint8(r) / 8
			convertedGreen := uint8(g) / 8
			convertedBlue := uint8(b) / 8
			convertedAlpha := uint8(a)
			if convertedAlpha != 255 {
				convertedAlpha = 0
			}
			c.Set(x, y, color.RGBA{convertedRed, convertedGreen, convertedBlue, convertedAlpha})
		}
	}
	return c
}

// New creates a new pixel Canvas.
func New(width, height int) Canvas {
	return Canvas{
		pixels:       make([]color.RGBA, width*height),
		pixelIndexes: make([]int, width*height),
		width:        width,
		height:       height,
	}
}

// Width returns the canvas width in pixels.
func (c *Canvas) Width() int {
	return c.width
}

// Height returns the canvas height in pixels.
func (c *Canvas) Height() int {
	return c.height
}

// At returns the color of a pixel.
func (c *Canvas) At(x, y int) color.RGBA {
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return color.RGBA{}
	}
	rowOffset := y * c.width
	return c.pixels[rowOffset+x]
}

// AtColorIndex returns the color index of a pixel.
func (c *Canvas) AtColorIndex(x, y int) int {
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return 0
	}
	rowOffset := y * c.width
	return c.pixelIndexes[rowOffset+x]
}

// Set a pixel color.
func (c *Canvas) Set(x, y int, color color.RGBA) {
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return
	}
	rowOffset := y * c.width
	c.pixels[rowOffset+x] = color
}

// SetColorIndex the color index for a pixel.
func (c *Canvas) SetColorIndex(x, y, index int) {
	if x < 0 || x >= c.width || y < 0 || y >= c.height {
		return
	}
	rowOffset := y * c.width
	c.pixelIndexes[rowOffset+x] = index
}

// ToImage returns an image representation of the Canvas.
func (c *Canvas) ToImage(palette []color.RGBA) image.Image {
	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: c.width, Y: c.height},
	})
	for x := 0; x < c.width; x++ {
		for y := 0; y < c.height; y++ {
			pixelIndex := c.AtColorIndex(x, y)
			if pixelIndex >= 0 && pixelIndex < len(palette) {
				paletteColor := palette[pixelIndex]
				r := paletteColor.R * 8
				g := paletteColor.G * 8
				b := paletteColor.B * 8
				a := paletteColor.A
				if a != 255 {
					a = 0
				}
				img.Set(x, y, color.RGBA{r, g, b, a})
			} else {
				img.Set(x, y, color.RGBA{})
			}
		}
	}
	return img
}
