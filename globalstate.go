package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type GlobalState struct {
	lastCounter uint64
	deltaTime   float32
	leftClick   bool
	rightClick  bool
}

func (this *GlobalState) Init() {
	this.lastCounter = sdl.GetPerformanceCounter()
	this.deltaTime = 0.0001
}

func (this *GlobalState) Frame() {
	counter := sdl.GetPerformanceCounter()
	this.deltaTime =
		float32(counter-this.lastCounter) /
			float32(sdl.GetPerformanceFrequency())
	this.lastCounter = counter

	this.leftClick = false
	this.rightClick = false
}

var globalState GlobalState
