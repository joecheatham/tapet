package main

import (
	"image"
	"image/color"
	"math/rand"
)

type Line struct {
	col      color.Color
	position int
	size     int
}

func Lines(colors []color.Color, w int, h int, size int, sizevar int, horizontal bool, equalspacing bool, spacingsize int, offset int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	var maxsize int
	if horizontal {
		maxsize = h
	} else {
		maxsize = w
	}

	currentposition := offset
	spacing := spacingsize

	lines := make([]Line, 0)
	bg := colors[0]

	for _, c := range colors {
		if c == bg {
			continue
		}
		line := Line{c, currentposition, randMinMax(size-sizevar, size+sizevar)}
		lines = append(lines, line)
		if !equalspacing {
			spacing = rand.Intn(maxsize / 16)
		}
		currentposition += line.size + spacing
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, bg)
			for _, l := range lines {
				var pixelpos int
				if horizontal {
					pixelpos = y
				} else {
					pixelpos = x
				}

				if pixelpos > l.position && pixelpos < l.position+l.size {
					img.Set(x, y, l.col)
				}
			}
		}
	}

	return img
}
