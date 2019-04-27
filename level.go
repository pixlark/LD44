package main

import "github.com/veandco/go-sdl2/sdl"

type Stopper struct {
	position int
	active   bool
}

func NewStopper(index int) Stopper {
	return Stopper{index, true}
}

type Path struct {
	orbIndex  int
	orbReset  int
	flagIndex int

	stoppers []Stopper

	orbPosition int
}

func NewPath(orbIndex, flagIndex int, stoppers []Stopper) Path {
	var p Path
	p.orbIndex = orbIndex
	p.flagIndex = flagIndex
	p.orbReset = 0
	p.orbPosition = p.orbReset
	p.stoppers = stoppers
	return p
}

type Level struct {
	paths []Path
	width int
}

func (this *Level) pathRect(path, pos int) sdl.Rect {
	rect := sdl.Rect{
		pathLeft, pathTop + int32(path)*pathVertSpace,
		pathRight - pathLeft, pathThickness,
	}
	rect.X += int32(pos) * (rect.W / int32(this.width-1))
	return rect
}

func (this *Level) Init() {
	this.paths = []Path{
		NewPath(1, 1, []Stopper{NewStopper(2)}),
		NewPath(2, 2, []Stopper{NewStopper(1)}),
		NewPath(3, 3, []Stopper{}),
	}
	this.width = 6
}

func (this *Level) Reset() {
	for i := range this.paths {
		path := &this.paths[i]
		path.orbPosition = path.orbReset
		for i := range path.stoppers {
			path.stoppers[i].active = true
		}
	}
}

func (this *Level) Step() {
	for i := range this.paths {
		path := &this.paths[i]

		// Stop if active stopper in the way
		stopped := false
		for i, stopper := range path.stoppers {
			if stopper.active && stopper.position == path.orbPosition {
				stopped = true
				path.stoppers[i].active = false
				break
			}
		}
		if stopped {
			continue
		}

		// Step path
		if path.orbPosition < this.width-1 {
			path.orbPosition++
		}
	}
}
