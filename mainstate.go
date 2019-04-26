package main

import (
	"math"
	"github.com/veandco/go-sdl2/sdl"
)

type MainState struct {
	t float32
}

func (this *MainState) Init() {
	
}

func (this *MainState) Update() Response {
	this.t += globalState.deltaTime
	
	return Response{RESPONSE_OK, nil}
}

func (this *MainState) Render(renderer *sdl.Renderer) Response {
	hsv := HSV{float32(math.Mod(float64(this.t * 12), 360.0)), 1.0, 1.0}
	color := hsv.Rgba()
	
	renderer.SetDrawColor(
		uint8(color.R * 255.0),
		uint8(color.G * 255.0),
		uint8(color.B * 255.0), 0xff)
	renderer.Clear()
	
	return Response{RESPONSE_OK, nil}
}

func (this *MainState) Exit() {
	
}
