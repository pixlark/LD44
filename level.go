package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

var _ = fmt.Println

type Tool interface {
	addToLevel(level *Level, row, col int)
	removeFromLevel(level *Level, row int)
}

type Stopper struct {
	position int
	active   bool
}

func newStopper(index int) Stopper {
	return Stopper{index, true}
}

func (this Stopper) addToLevel(level *Level, row, col int) {
	path := &level.paths[row]
	this.position = col
	path.stoppers = append(path.stoppers, this)
}

func (this Stopper) removeFromLevel(level *Level, row int) {
	path := &level.paths[row]
	for i := range path.stoppers {
		if path.stoppers[i] == this {
			path.stoppers[i] = path.stoppers[len(path.stoppers)-1]
			path.stoppers = path.stoppers[:len(path.stoppers)-1]
		}
	}
}

// Swappers are represented by their top anchor
type VertSwapper struct {
	position int
}

func newVertSwapper(position int) VertSwapper {
	return VertSwapper{position}
}

func (this VertSwapper) addToLevel(level *Level, row, col int) {
	path := &level.paths[row]
	this.position = col
	path.vertSwappers = append(path.vertSwappers, this)
}

func (this VertSwapper) removeFromLevel(level *Level, row int) {
	path := &level.paths[row]
	for i := range path.vertSwappers {
		if path.vertSwappers[i] == this {
			path.vertSwappers[i] = path.vertSwappers[len(path.vertSwappers)-1]
			path.vertSwappers = path.vertSwappers[:len(path.vertSwappers)-1]			
		}
	}
}

type Path struct {
	orbIndex  int
	orbReset  int
	flagIndex int

	stoppers     []Stopper
	vertSwappers []VertSwapper

	orbPosition int
}

func newPath(orbIndex, flagIndex int, stoppers []Stopper, vertSwappers []VertSwapper) Path {
	var p Path
	p.orbIndex = orbIndex
	p.flagIndex = flagIndex
	p.orbReset = 0
	p.orbPosition = p.orbReset
	p.stoppers = stoppers
	p.vertSwappers = vertSwappers
	return p
}

type Level struct {
	paths []Path
	width int
}

func (this *Level) baseRect(path, pos int) sdl.Rect {
	rect := sdl.Rect{
		pathLeft, pathTop + int32(path)*pathVertSpace,
		0, 0,
	}
	rect.X += int32(pos) * ((pathRight - pathLeft) / int32(this.width-1))
	return rect
}

func (this *Level) pathRect(path int) sdl.Rect {
	rect := this.baseRect(path, 0)
	rect.W = pathRight - pathLeft
	rect.H = pathThickness
	return rect
}

func (this *Level) stopperRect(path, pos int) sdl.Rect {
	rect := this.baseRect(path, pos)
	rect.X -= stopperSize / 2
	rect.Y -= stopperSize / 2
	rect.W = stopperSize
	rect.H = stopperSize
	return rect
}

func (this *Level) swapperRect(path, pos int) sdl.Rect {
	rect := this.baseRect(path, pos)
	rect.X -= swapperWidth / 2
	rect.Y += swapperPad
	rect.W = swapperWidth
	rect.H = swapperHeight
	return rect
}

func (this *Level) inRangeOfToolSpot(tool Tool, x, y int32) (int, int, bool) {
	switch tool.(type) {
	case Stopper:
		for row := range this.paths {
			for col := 1; col < this.width-1; col++ {
				rect := this.baseRect(row, col)
				if distance(rect.X, rect.Y, x, y) < (float32(pathVertSpace) / 2.0) {
					return row, col, true
				}
			}
		}
	case VertSwapper:
		for row := 0; row < len(this.paths)-1; row++ {
			for col := 1; col < this.width-1; col++ {
				rect := this.swapperRect(row, col)
				centerX := rect.X + rect.W / 2
				centerY := rect.Y + rect.H / 2
				if distance(centerX, centerY, x, y) < (float32(pathVertSpace) / 2.0) {
					return row, col, true
				}
			}
		}
	}
	return 0, 0, false
}

func (this *Level) init() {
	this.paths = []Path{
		newPath(1, 1, []Stopper{newStopper(2)}, []VertSwapper{}),
		newPath(2, 2, []Stopper{newStopper(1)}, []VertSwapper{newVertSwapper(3)}),
		newPath(3, 3, []Stopper{}, []VertSwapper{}),
	}
	this.width = 6
}

func (this *Level) reset() {
	for i := range this.paths {
		path := &this.paths[i]
		path.orbPosition = path.orbReset
		for i := range path.stoppers {
			path.stoppers[i].active = true
		}
	}
}

func (this *Level) step() {
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

const (
	TOOL_NONE    = iota
	TOOL_STOPPER = iota
)

func (this *Level) canDragTool() (Tool, int, int) {
	mx, my, _ := sdl.GetMouseState()
	for row := range this.paths {
		path := &this.paths[row]
		// Check for stoppers
		for _, stopper := range path.stoppers {
			rect := this.stopperRect(row, stopper.position)
			if globalState.leftClick && inRect(mx, my, rect) {
				return stopper, row, stopper.position
			}
		}
		// Check for vertical swappers
		for _, swapper := range path.vertSwappers {
			rect := this.swapperRect(row, swapper.position)
			if globalState.leftClick && inRect(mx, my, rect) {
				return swapper, row, swapper.position
			}
		}
	}
	return nil, -1, -1
}
