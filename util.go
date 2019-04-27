package main

import (
	"math"
	"github.com/veandco/go-sdl2/sdl"
)

func inRect(mx, my int32, rect sdl.Rect) bool {
	return (&sdl.Point{mx, my}).InRect(&rect)
}

func distance(x0, y0, x1, y1 int32) float32 {
	return float32(math.Sqrt(float64((x1 - x0) * (x1 - x0) + (y1 - y0) * (y1 - y0))))
}
