package main

import (
	"math"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"github.com/veandco/go-sdl2/img"
)

func inRect(mx, my int32, rect sdl.Rect) bool {
	return (&sdl.Point{mx, my}).InRect(&rect)
}

func distance(x0, y0, x1, y1 int32) float32 {
	return float32(math.Sqrt(float64((x1 - x0) * (x1 - x0) + (y1 - y0) * (y1 - y0))))
}

func loadTexture(renderer *sdl.Renderer, path string) *sdl.Texture {
	texture, err := img.LoadTexture(renderer, path)
	if err != nil {
		fatal("Could not open texture")
	}
	return texture
}

func loadFont(path string, size int) *ttf.Font {
	font, err := ttf.OpenFont(path, size)
	if err != nil {
		fatal("Could not open font")
	}
	return font
}

func fontRender(renderer *sdl.Renderer, font *ttf.Font, text string, color sdl.Color) *sdl.Texture {
	surface, _ := font.RenderUTF8Blended(text, color)
	defer surface.Free()
	texture, _ := renderer.CreateTextureFromSurface(surface)
	return texture
}

func centerRectInRect(inner, outer sdl.Rect) sdl.Rect {
	ret := outer
	
	offsetX := (outer.W - inner.W) / 2
	ret.X += offsetX
	ret.W = inner.W
	
	offsetY := (outer.H - inner.H) / 2
	ret.Y += offsetY
	ret.H = inner.H

	return ret
}
