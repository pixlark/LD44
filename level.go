package main

import "github.com/veandco/go-sdl2/sdl"

type Path struct {
	orbIndex    int
	orbReset    int
	flagIndex   int

	orbPosition int
}

func NewPath(orbIndex, flagIndex int) Path {
	var p Path
	p.orbIndex = orbIndex
	p.flagIndex = flagIndex
	p.orbReset = 0
	p.orbPosition = p.orbReset
	return p
}

func pathRect(i int) sdl.Rect {
	return sdl.Rect{
		pathLeft, pathTop + int32(i)*pathVertSpace,
		pathRight - pathLeft, pathThickness,
	}
}

type Level struct {
	paths  []Path
	width  int
}

func (this *Level) Init() {
	this.paths = []Path{
		NewPath(1, 1),
		NewPath(2, 2),
		NewPath(3, 3),
	}
	this.width = 6
}

func (this *Level) Reset() {
	for i := range this.paths {
		path := &this.paths[i]
		path.orbPosition = path.orbReset
	}
}

func (this *Level) Step() {
	for i := range this.paths {
		path := &this.paths[i]
		if path.orbPosition < this.width - 1 {
			path.orbPosition++
		}
	}
}
