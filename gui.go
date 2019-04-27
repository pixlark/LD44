package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var _ = fmt.Println

func button(renderer *sdl.Renderer, font *ttf.Font, rect sdl.Rect, text string) bool {
	// Render
	renderer.SetDrawColor(0xff, 0xff, 0xff, 0xff)
	renderer.FillRect(&rect)
	texture := fontRender(renderer, font, text)
	defer texture.Destroy()
	renderer.Copy(texture, nil, &rect)

	// Check for click
	mx, my, _ := sdl.GetMouseState()
	if (globalState.leftClick) {
		if mx > rect.X && mx < rect.X+rect.W &&
			my > rect.Y && my < rect.Y+rect.H {
			return true
		}
	}
	return false
}
