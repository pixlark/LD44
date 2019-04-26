package main

import (
	_ "fmt"
	"math"
	"github.com/veandco/go-sdl2/sdl"
)

type RGBA struct {
	R float32
	G float32
	B float32
	A float32
}

type HSV struct {
	H float32
	S float32
	V float32
}

func (this HSV) Rgba() RGBA {
	var hh, p, q, t, ff float32
	var i   int
	var out RGBA
	out.A = 1.0

	if this.S <= 0.0 {
		out.R = this.V
		out.G = this.V
		out.B = this.V
		return out
	}
	
	hh = this.H
	if (hh >= 360.0) {
		hh = 0.0
	}
	hh /= 60.0
	i = int(hh)
	ff = hh - float32(i)
	p = this.V * (1.0 - this.S)
	q = this.V * (1.0 - (this.S * ff))
	t = this.V * (1.0 - (this.S * (1.0 - ff)))

	switch i {
	case 0:
		out.R = this.V
		out.G = t
		out.B = p
	case 1:
		out.R = q
		out.G = this.V
		out.B = p
	case 2:
		out.R = p
		out.G = this.V
		out.B = t
	case 3:
		out.R = p
		out.G = q
		out.B = this.V
	case 4:
		out.R = t
		out.G = p
		out.B = this.V
	default:
		out.R = this.V
		out.G = p
		out.B = q
	}
	return out
}

/*
 * GLOBAL STATE
 */

type GlobalState struct {
	lastCounter uint64
	deltaTime float32
}

func (this *GlobalState) Init() {
	this.lastCounter = sdl.GetPerformanceCounter()
	this.deltaTime = 0.0001
}

var globalState GlobalState

/*
 * MAIN STATE
 */

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

/*
 * STATE
 */

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

/*
 * MAIN
 */

func main() {
	sdl.Init(sdl.INIT_EVERYTHING);
	
	window, _ := sdl.CreateWindow("title!", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, 0)
	renderer, _ := sdl.CreateRenderer(window, -1, 0)

	states := make([]State, 4)
	states = append(states, &MainState{})
	states[len(states) - 1].Init()

	globalState.Init()
	
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
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

		{
			counter := sdl.GetPerformanceCounter()
			globalState.deltaTime =
				float32(counter - globalState.lastCounter) /
				float32(sdl.GetPerformanceFrequency())
			globalState.lastCounter = counter
		}
		
		renderer.Present()
	}
}
