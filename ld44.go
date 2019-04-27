package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var _ = fmt.Println

const (
	RESPONSE_OK   = iota
	RESPONSE_PUSH = iota
	RESPONSE_POP  = iota
)

type Response struct {
	code  int
	state State
}

type State interface {
	Init(renderer *sdl.Renderer)
	Update(events []sdl.Event) Response
	Render(renderer *sdl.Renderer) Response
	Exit()
}

const (
	screenW int32 = 800
	screenH int32 = 600
)

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)
	sdl.StopTextInput() // This is enabled by default for some bizarre reason
	ttf.Init()

	window, _ := sdl.CreateWindow(
		"LD44", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenW, screenH, 0)
	renderer, _ := sdl.CreateRenderer(window, -1, 0)

	states := make([]State, 0)
	states = append(states, &MainState{})
	states[len(states)-1].Init(renderer)

	globalState.Init()

mainloop:
	for len(states) > 0 {
		events := make([]sdl.Event, 0)

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			// Some events are handled globally, no matter what state frame we're in
			switch event := event.(type) {
			case *sdl.QuitEvent:
				break mainloop
			case *sdl.KeyboardEvent:
				switch event.Type {
				case sdl.KEYDOWN:
					if event.Keysym.Scancode == sdl.SCANCODE_ESCAPE {
						break mainloop
					}
				}
			case *sdl.MouseButtonEvent:
				switch event.Type {
				case sdl.MOUSEBUTTONDOWN:
					if event.Button == sdl.BUTTON_LEFT {
						globalState.leftClick = true
					} else if event.Button == sdl.BUTTON_RIGHT {
						globalState.rightClick = true
					}
				}
			}
			events = append(events, event)
		}

		response := states[len(states)-1].Update(events)
		switch response.code {
		case RESPONSE_OK:
		case RESPONSE_PUSH:
			states = append(states, response.state)
			states[len(states)-1].Init(renderer)
			continue
		case RESPONSE_POP:
			states[len(states)-1].Exit()
			states = states[:len(states)-1]
			continue
		}

		response = states[len(states)-1].Render(renderer)
		switch response.code {
		case RESPONSE_OK:
		case RESPONSE_PUSH:
			states = append(states, response.state)
		case RESPONSE_POP:
			states = states[:len(states)-1]
		}
		renderer.Present()

		globalState.Frame()
	}
}
