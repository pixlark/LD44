package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

var _ = fmt.Println // debugging

const (
	RESPONSE_OK   = iota
	RESPONSE_PUSH = iota
	RESPONSE_POP  = iota
)

type Response struct {
	code int
	state State
}

type State interface {
	Init()
	Update() Response
	Render(renderer *sdl.Renderer) Response
	Exit()
}

func main() {
	sdl.Init(sdl.INIT_EVERYTHING);
	
	window, _ := sdl.CreateWindow("title!", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, 0)
	renderer, _ := sdl.CreateRenderer(window, -1, 0)

	states := make([]State, 4)
	states = append(states, &MainState{})
	states[len(states) - 1].Init()

	globalState.Init()

	mainloop:
	for len(states) > 0 {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				break mainloop
			}
		}
		
		response := states[len(states) - 1].Update()
		switch response.code {
		case RESPONSE_OK:
		case RESPONSE_PUSH:
			states = append(states, response.state)
			states[len(states) - 1].Init()
			continue
		case RESPONSE_POP:
			states[len(states) - 1].Exit()
			states = states[:len(states) - 1]
			continue
		}

		response = states[len(states) - 1].Render(renderer)
		switch response.code {
		case RESPONSE_OK:
		case RESPONSE_PUSH:
			states = append(states, response.state)
		case RESPONSE_POP:
			states = states[:len(states) - 1]
		}
		renderer.Present()
		
		globalState.Frame()
	}
}
