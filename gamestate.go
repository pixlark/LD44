package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

var _ = fmt.Println

type GameState struct {
	
}

func (this *GameState) Init() {
	fmt.Println("Entering GameState")
}

func (this *GameState) Update(events []sdl.Event) Response {
	for _, event := range events {
		switch event := event.(type) {
		case *sdl.KeyboardEvent:
			switch event.Type {
			case sdl.KEYDOWN:
				fmt.Println("Hello!");
				return Response{RESPONSE_POP, nil}
			}
		}
	}
	
	return Response{RESPONSE_OK, nil}
}

func (this *GameState) Render(renderer *sdl.Renderer) Response {
	renderer.SetDrawColor(0, 0, 0, 0xff)
	renderer.Clear()
	
	return Response{RESPONSE_OK, nil}
}

func (this *GameState) Exit() {
	fmt.Println("Exiting GameState")
}
