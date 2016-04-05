package main

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"sort"
)

type Circle struct {
	col  color.Color
	x, y int
	size int
}

type circleBySize []Circle

func (a circleBySize) Len() int           { return len(a) }
func (a circleBySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a circleBySize) Less(i, j int) bool { return a[i].size < a[j].size }

func Circles(colors []color.Color, w int, h int, size int, sizevar int, overlap bool, large2small bool, filled bool, bordersize int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))

	circles := make([]Circle, 0)
	bg := colors[0]

	for _, c := range colors {
		if c == bg {
			continue
		}
		circle := Circle{c, rand.Intn(w), rand.Intn(h), randMinMax(size-sizevar, size+sizevar)}
		circles = append(circles, circle)
	}

	if large2small {
		sort.Sort(sort.Reverse(circleBySize(circles)))
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, bg)
			for _, c := range circles {
				a := float64((x - c.x) * (x - c.x))
				b := float64((y - c.y) * (y - c.y))

				if filled {
					if int(math.Sqrt(a+b)) < c.size {
						img.Set(x, y, c.col)
					}
				} else {
					if int(math.Sqrt(a+b)) < c.size && int(math.Sqrt(a+b)) > (c.size-bordersize) {
						img.Set(x, y, c.col)
					}
				}
			}
		}
	}
	return img
}
