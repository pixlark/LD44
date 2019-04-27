package main

import "github.com/veandco/go-sdl2/sdl"

type Path struct {
	orbIndex    int
	orbPosition int
	flagIndex   int
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
		Path{1, 0, 1},
		Path{2, 1, 2},
		Path{3, 5, 3},
	}
	
	this.width = 6
}
