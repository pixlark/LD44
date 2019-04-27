package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

var _ = fmt.Println

type MainState struct {
	t float32
}

func (this *MainState) init(renderer *sdl.Renderer) {
	
}

func (this *MainState) update(events []sdl.Event) Response {
	gameState := gameStateWithLevelPath("level0.json");
	return Response{RESPONSE_PUSH, &gameState}
}

func (this *MainState) render(renderer *sdl.Renderer) Response {
	return Response{RESPONSE_OK, nil}
}

func (this *MainState) exit() {
	
}
