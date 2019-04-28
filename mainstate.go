package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var _ = fmt.Println

const (
	selectorsX = 5
	selectorsY = 4

	selectorAntiPad = 10
	selectorW       = (screenW / selectorsX) - (selectorAntiPad * 2)
	selectorH       = (screenH / selectorsY) - (selectorAntiPad * 2)
)

type MainState struct {
	font      *ttf.Font
}

func (this *MainState) init(renderer *sdl.Renderer) {
	this.font = loadFont("DejaVuSans.ttf", 60)
}

func (this *MainState) update(events []sdl.Event) Response {
	/*
		gameState := gameStateWithLevelPath("level0.json");
		return Response{RESPONSE_PUSH, &gameState}*/
	return Response{RESPONSE_OK, nil}
}

func (this *MainState) render(renderer *sdl.Renderer) Response {
	renderer.SetDrawColor(0, 0, 0, 0xff)
	renderer.Clear()

	for row := 0; row < selectorsY; row++ {
		for col := 0; col < selectorsX; col++ {
			clicked := button(renderer, this.font,
				sdl.Rect{
					int32(col * int(selectorW + selectorAntiPad * 2)),
					int32(row * int(selectorH + selectorAntiPad * 2)),
					selectorW,
					selectorH},
				fmt.Sprintf("%d", row * selectorsX + col),
			)
			if clicked {
				path := fmt.Sprintf("level%d.json", row * selectorsX + col)
				gameState := gameStateWithLevelPath(path)
				return Response{RESPONSE_PUSH, &gameState}
			}
		}
	}

	return Response{RESPONSE_OK, nil}
}

func (this *MainState) exit() {

}
