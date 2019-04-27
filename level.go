package main

import "github.com/veandco/go-sdl2/sdl"

type Path struct {
	orbIndex  int
	orbReset  int
	flagIndex int

	stoppers       []bool
	activeStoppers []bool

	orbPosition int
}

func NewPath(orbIndex, flagIndex int, stoppers []bool) Path {
	var p Path
	p.orbIndex = orbIndex
	p.flagIndex = flagIndex
	p.orbReset = 0
	p.orbPosition = p.orbReset
	p.stoppers = stoppers
	p.activeStoppers = make([]bool, len(stoppers))
	copy(p.activeStoppers, stoppers)
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
	rect.X += int32(pos) * (rect.W / int32(this.width - 1))
	return rect
}

func (this *Level) Init() {
	this.paths = []Path{
		NewPath(1, 1, []bool{false, false, true, false, false, false}),
		NewPath(2, 2, []bool{false, true, false, false, false, false}),
		NewPath(3, 3, []bool{false, false, false, false, false, false}),
	}
	this.width = 6
}

func (this *Level) Reset() {
	for i := range this.paths {
		path := &this.paths[i]
		path.orbPosition = path.orbReset
		copy(path.activeStoppers, path.stoppers)
	}
}

func (this *Level) Step() {
	for i := range this.paths {
		path := &this.paths[i]

		// Stop if active stopper in the way
		if path.activeStoppers[path.orbPosition] {
			path.activeStoppers[path.orbPosition] = false
			continue
		}
		
		// Step path
		if path.orbPosition < this.width-1 {
			path.orbPosition++
		}
	}
}
