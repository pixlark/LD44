package main

import (
	"fmt"
	//"math"
	"github.com/veandco/go-sdl2/sdl"
)

var _ = fmt.Println

type MainState struct {
	t float32
}

func (this *MainState) init(renderer *sdl.Renderer) {
	
}

func (this *MainState) update(events []sdl.Event) Response {
	/*
	for _, event := range events {
		switch event := event.(type) {
		case *sdl.KeyboardEvent:
			switch event.Type {
			case sdl.KEYDOWN:
				return Response{RESPONSE_PUSH, &GameState{}}
			}
		}
	}
	
	this.t += globalState.deltaTime
	return Response{RESPONSE_OK, nil}*/

	return Response{RESPONSE_PUSH, &GameState{}}
}

func (this *MainState) render(renderer *sdl.Renderer) Response {
	/*
	hsv := HSV{float32(math.Mod(float64(this.t * 12), 360.0)), 1.0, 1.0}
	color := hsv.Rgba()
	
	renderer.SetDrawColor(
		uint8(color.R * 255.0),
		uint8(color.G * 255.0),
		uint8(color.B * 255.0), 0xff)
	renderer.Clear()*/
	
	return Response{RESPONSE_OK, nil}
}

func (this *MainState) exit() {
	
}
