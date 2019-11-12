package contestpaintingeffects

import (
	"image/color"

	"github.com/huderlem/contest-painting-effects/canvas"
	"github.com/huderlem/contest-painting-effects/effect"
	"github.com/huderlem/contest-painting-effects/paletteq"
)

// ApplyCoolEffect applies the effects used for Cool contest winner paintings.
func ApplyCoolEffect(c canvas.Canvas, personality uint8) []color.RGBA {
	effect.ApplyBlackOutline(c)
	effect.ApplyPersonalityColor(c, personality)
	return paletteq.ApplyStandardQuantization(c, 224)
}

// ApplyBeautyEffect applies the effects used for Beauty contest winner paintings.
func ApplyBeautyEffect(c canvas.Canvas) []color.RGBA {
	effect.ApplyShimmer(c)
	return paletteq.ApplyStandardQuantization(c, 224)
}

// ApplyCuteEffect applies the effects used for Cute contest winner paintings.
func ApplyCuteEffect(c canvas.Canvas) []color.RGBA {
	effect.ApplyPointillism(c)
	return paletteq.ApplyStandardQuantization(c, 224)
}

// ApplySmartEffect applies the effects used for Smart contest winner paintings.
func ApplySmartEffect(c canvas.Canvas) []color.RGBA {
	effect.ApplyBlackOutline(c)
	effect.ApplyBlurRight(c)
	effect.ApplyBlurDown(c)
	effect.ApplyBlackAndWhite(c)
	effect.ApplyBlur(c)
	effect.ApplyBlur(c)
	effect.ApplyRedChannelGrayscale(c, 2)
	effect.ApplyRedChannelGrayscaleHighlight(c, 4)
	return paletteq.ApplyGrayscaleQuantization(c)
}

// ApplyToughEffect applies the effects used for Tough contest winner paintings.
func ApplyToughEffect(c canvas.Canvas) []color.RGBA {
	effect.ApplyGrayscale(c)
	effect.ApplyRedChannelGrayscale(c, 3)
	return paletteq.ApplyGrayscaleQuantization(c)
}
