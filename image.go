package main

import (
	"errors"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	threshold     int = 50
	minBrightness int = 0
	maxBrightness int = 200
)

func colorDifference(col1 color.Color, col2 color.Color, threshold int) bool {
	c1 := col1.(color.NRGBA)
	c2 := col2.(color.NRGBA)

	rDiff := abs(int(c1.R) - int(c2.R))
	gDiff := abs(int(c1.G) - int(c2.G))
	bDiff := abs(int(c1.B) - int(c2.B))

	total := rDiff + gDiff + bDiff
	return total >= threshold
}

func colorsFromImage(filename string) ([]color.Color, error) {
	fuzzyness := 5
	img := loadImage(filename)
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	colors := make([]color.Color, 0, w*h)
	for x := 0; x < w; x += fuzzyness {
		for y := 0; y < h; y += fuzzyness {
			col := color.NRGBAModel.Convert(img.At(x, y))
			colors = append(colors, col)
		}
	}
	distinctColors := getDistinctColors(colors, threshold, minBrightness, maxBrightness)

	count := 0
	for len(distinctColors) < 16 {
		count++
		distinctColors = append(distinctColors, getDistinctColors(colors, threshold-count, minBrightness, maxBrightness)...)
		if count == threshold {
			return nil, errors.New("Could not get colors from image with settings specified. Aborting.\n")
		}
	}

	if len(distinctColors) > 16 {
		distinctColors = distinctColors[:16]
	}
	return distinctColors, nil
}

func getDistinctColors(colors []color.Color, threshold int, minBrightness, maxBrightness int) []color.Color {
	distinctColors := make([]color.Color, 0)
	for _, c := range colors {
		same := false
		if !colorDifference(c, color.NRGBAModel.Convert(color.Black), minBrightness*3) {
			continue
		}
		if !colorDifference(c, color.NRGBAModel.Convert(color.White), (255-maxBrightness)*3) {
			continue
		}
		for _, k := range distinctColors {
			if !colorDifference(c, k, threshold) {
				same = true
				break
			}
		}
		if !same {
			distinctColors = append(distinctColors, c)
		}
	}
	return distinctColors
}

func loadImage(filepath string) image.Image {
	infile, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	src, _, err := image.Decode(infile)
	if err != nil {
		log.Fatal(err)
	}
	return src
}

func randomImage(colors []color.Color, w int, h int) image.Image {
	rand.Seed(time.Now().UnixNano())
	switch rand.Intn(3) {
	case 0:
		return Circles(colors, w, h, rand.Intn(w/2), rand.Intn(w/2), randBool(), randBool(), randBool(), rand.Intn(20))
	case 1:
		return Rays(colors, w, h, rand.Intn(h/32)+1, rand.Intn(h/32), randBool(), randBool(), randBool())
	case 2:
		return Lines(colors, w, h, rand.Intn(h/32)+1, rand.Intn(h/32), randBool(), randBool(), rand.Intn(h/32), rand.Intn(h/2)+1)
	}
	return nil
}
