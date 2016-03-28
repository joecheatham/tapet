package main

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"sort"
	)

type Ray struct {
	col   color.Color
	x, y  int
	angle int
	size  int
}

type rayBySize []Ray

func (a rayBySize) Len() int { return len(a) }
func (a rayBySize) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a rayBySize) Less(i, j int) bool { return a[i].size < a[j].size }

func Rays(colors []color.Color, w int, h int, size int, sizevar int, evendist bool, centered bool, large2small bool) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))

	rays := make([]Ray, 0)

	spacing := 360 / len(colors)
	current_angle := 0

	xpos := w / 2
	ypos := h / 2

	bg := colors[0]

	for _, c := range colors {
		if c == bg {
			continue
		}
		var ray Ray
		if !centered {
			xpos = rand.Intn(w)
			ypos = rand.Intn(h)
		}
		if !evendist {
			current_angle = rand.Intn(360)
		}
		ray = Ray{c, xpos, ypos, current_angle, randMinMax(size - sizevar, size + sizevar)}

		if evendist {
			current_angle += spacing + ray.size
		}
		rays = append(rays, ray)
	}

	if large2small {
		sort.Sort(sort.Reverse(rayBySize(rays)))
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, bg)
			for _, r := range rays {
				deltaX := float64(x - r.x)
				deltaY := float64(y - r.y)
				angle := math.Atan(deltaY / deltaX) * 180 / math.Pi
				if angle < 0 {
					angle += 360
				}
				if int(math.Abs(float64(int(angle) - r.angle))) < r.size {
					img.Set(x, y, r.col)
				}
			}
		}
	}
	return img
}