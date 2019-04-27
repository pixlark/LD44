package main

import "github.com/veandco/go-sdl2/sdl"

func inRect(mx, my int32, rect sdl.Rect) bool {
	return (&sdl.Point{mx, my}).InRect(&rect)
}
