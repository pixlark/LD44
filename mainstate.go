package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var _ = fmt.Println

const (
	selectorsX = 5
	selectorsY = 2

	selectorAntiPad = 20
	selectorW       = (screenW / selectorsX) - (selectorAntiPad) - (selectorAntiPad / selectorsX)
	selectorH       = (screenH / selectorsY) - (selectorAntiPad) - (selectorAntiPad / selectorsY)
)

type MainState struct {
	font      *ttf.Font
	completed []bool
}

func (this *MainState) init(renderer *sdl.Renderer) {
	this.font = loadFont("DejaVuSans.ttf", 60)
	this.completed = make([]bool, selectorsX * selectorsY)
}

func (this *MainState) update(events []sdl.Event) Response {
	return Response{RESPONSE_OK, nil}
}

func (this *MainState) render(renderer *sdl.Renderer) Response {
	renderer.SetDrawColor(0, 0, 0, 0xff)
	renderer.Clear()

	for row := 0; row < selectorsY; row++ {
		for col := 0; col < selectorsX; col++ {
			index := row * selectorsX + col
			var color sdl.Color
			if this.completed[index] {
				color = sdl.Color{0xcc, 0xcc, 0xcc, 0xff}
			} else {
				color = white
			}
			clicked := button(renderer, this.font,
				sdl.Rect{
					int32(col * int(selectorW + selectorAntiPad) + selectorAntiPad),
					int32(row * int(selectorH + selectorAntiPad) + selectorAntiPad),
					selectorW,
					selectorH,
				},
				fmt.Sprintf("%d", index),
				color,
			)
			if clicked && !this.completed[index] {
				path := fmt.Sprintf("level%d.json", index)
				gameState := gameStateWithLevelPath(path)
				// Since the only way to get back to the main menu is
				// to win, we're just doing a hack where we set this
				// level to complete as soon as you start it up
				this.completed[index] = true
				return Response{RESPONSE_PUSH, &gameState}
			}
		}
	}

	return Response{RESPONSE_OK, nil}
}

func (this *MainState) exit() {

}
