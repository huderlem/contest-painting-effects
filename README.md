# Contest Painting Effects

This is an image processing library that perfectly recreates the Pokémon Contest painting effects found in Pokémon Emerald, Ruby, and Sapphire. The code is directly ported from the [pokeemerald](https://github.com/pret/pokeemerald/blob/master/src/image_processing_effects.c) project.

A full working example of using this library is found below. It processes the given image with the "Smart Contest" effect.

```go
package main

import (
	"fmt"
	"image"
	"image/png"
	_ "image/png"
	"log"
	"os"

	contestpaintingeffects "github.com/huderlem/contest-painting-effects"
	"github.com/huderlem/contest-painting-effects/canvas"
)

func loadImage() (image.Image, error) {
	imageFile, err := os.Open("dusclops.png")
	if err != nil {
		return nil, fmt.Errorf("Error opening input image file: %s", err.Error())
	}
	defer imageFile.Close()

	imageData, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, fmt.Errorf("Error decoding image file: %s", err.Error())
	}
	return imageData, nil
}

func saveImage(img image.Image) error {
	outputFile, err := os.Create("output.png")
	if err != nil {
		return fmt.Errorf("Error saving image: %s", err.Error())
	}
	defer outputFile.Close()

	png.Encode(outputFile, img)
	return nil
}

func main() {
	imageData, err := loadImage()
	if err != nil {
		log.Fatalf(err.Error())
	}

	c := canvas.FromImage(imageData)
	palette := contestpaintingeffects.ApplySmartEffect(c)
	saveImage(c.ToImage(palette))
}
```
