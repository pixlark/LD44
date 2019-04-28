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
			break
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
			break
		}
	}
}

type Path struct {
	orbIndex      int
	orbIndexReset int

	orbPosition      int
	orbPositionReset int

	flagIndex int

	orbSwappedThisFrame bool

	stoppers     []Stopper
	vertSwappers []VertSwapper
}

func newPath(start, orbIndex, flagIndex int, stoppers []Stopper, vertSwappers []VertSwapper) Path {
	var p Path
	p.orbIndex = orbIndex
	p.orbIndexReset = orbIndex
	p.flagIndex = flagIndex

	p.orbPositionReset = start
	p.orbPosition = p.orbPositionReset

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

func (this *Level) stopperAt(row, col int) bool {
	path := this.paths[row]
	for _, stopper := range path.stoppers {
		if stopper.position == col {
			return true
		}
	}
	return false
}

func (this *Level) vertSwapperAt(row, col int) bool {
	path := this.paths[row]
	for _, swapper := range path.vertSwappers {
		if swapper.position == col {
			return true
		}
	}
	return false
}

func (this *Level) inRangeOfToolSpot(tool Tool, x, y int32) (int, int, bool) {
	switch tool.(type) {
	case Stopper:
		for row := range this.paths {
			for col := this.paths[row].orbPosition + 1; col < this.width-1; col++ {
				rect := this.baseRect(row, col)
				if distance(rect.X, rect.Y, x, y) < (float32(pathVertSpace) / 2.0) {
					if !this.stopperAt(row, col) {
						return row, col, true
					}
				}
			}
		}
	case VertSwapper:
		for row := 0; row < len(this.paths)-1; row++ {
			for col := this.paths[row].orbPosition + 1; col < this.width-1; col++ {
				rect := this.swapperRect(row, col)
				centerX := rect.X + rect.W/2
				centerY := rect.Y + rect.H/2
				if distance(centerX, centerY, x, y) < (float32(pathVertSpace) / 2.0) {
					if !this.vertSwapperAt(row, col) {
						return row, col, true
					}
				}
			}
		}
	}
	return 0, 0, false
}

func (this *Level) init() {
	
}

func (this *Level) reset() {
	for i := range this.paths {
		path := &this.paths[i]
		path.orbIndex = path.orbIndexReset
		path.orbPosition = path.orbPositionReset
		for i := range path.stoppers {
			path.stoppers[i].active = true
		}
	}
}

func (this *Level) step() {
	for r := range this.paths {
		this.paths[r].orbSwappedThisFrame = false
	}
pathloop:
	for r := range this.paths {
		path := &this.paths[r]

		// Stop if active stopper in the way
		for i, stopper := range path.stoppers {
			if stopper.active && stopper.position == path.orbPosition {
				path.stoppers[i].active = false
				continue pathloop
			}
		}

		// Swap if both ends of a swapper have connected
		if r != len(this.paths)-1 && !path.orbSwappedThisFrame {
			for _, swapper := range path.vertSwappers {
				topConnect := swapper.position == path.orbPosition
				otherPath := &this.paths[r+1]
				botConnect := swapper.position == otherPath.orbPosition

				if topConnect && botConnect {
					tmp := path.orbIndex
					path.orbIndex = otherPath.orbIndex
					otherPath.orbIndex = tmp

					path.orbSwappedThisFrame = true
					otherPath.orbSwappedThisFrame = true
				}
			}
		}

		// Step path
		if path.orbPosition < this.width-1 {
			path.orbPosition++
		}
	}
}

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

func (this *Level) checkEnded() (bool, bool) {
	ended := true
	for _, path := range this.paths {
		if path.orbPosition != this.width-1 {
			ended = false
		}
	}
	success := true
	if ended {
		for _, path := range this.paths {
			if path.flagIndex != path.orbIndex {
				success = false
			}
		}
	}
	return ended, success
}
