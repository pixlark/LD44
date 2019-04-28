package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var _ = fmt.Println

const (
	pathLeftPad   int32 = 40
	pathTopPad    int32 = 100
	pathTop       int32 = pathTopPad
	pathThickness int32 = 5
	pathVertSpace int32 = screenH / 5
	pathLeft      int32 = pathLeftPad
	pathRight     int32 = screenW - pathLeftPad
	pathEndWidth  int32 = 10

	orbSize int32 = 40

	stopperSize int32 = 25

	swapperWidth  int32 = 10
	swapperPad    int32 = 10
	swapperHeight int32 = pathVertSpace - (swapperPad * 2)

	secondsPerStep float32 = 1.0
)

const (
	INTERACT_SETUP = iota
	INTERACT_GOING = iota
	INTERACT_PAUSE = iota
)

type GameState struct {
	levelPath   string
	level       Level
	assets      map[string]*sdl.Texture
	font        *ttf.Font
	interaction int
	stepTimer   float32

	transientResetRow    int
	transientResetColumn int
	transientTool        Tool
}

func gameStateWithLevelPath(path string) GameState {
	var state GameState
	state.levelPath = path
	return state
}

func (this *GameState) init(renderer *sdl.Renderer) {
	this.assets = make(map[string]*sdl.Texture)
	this.assets["orb"] = loadTexture(renderer, "orb.png")
	this.assets["flag"] = loadTexture(renderer, "flag.png")
	this.assets["path-left"] = loadTexture(renderer, "path-left.png")
	this.assets["path-right"] = loadTexture(renderer, "path-right.png")
	this.assets["path"] = loadTexture(renderer, "path.png")

	this.font = loadFont("DejaVuSans.ttf", 32)
	this.interaction = INTERACT_SETUP

	//fmt.Printf("Loading %s\n", this.levelPath)
	this.level = loadLevel(this.levelPath)
}

func (this *GameState) update(events []sdl.Event) Response {
	for _, event := range events {
		var _ = event
	}

	// Update step timer if we're going
	if this.interaction == INTERACT_GOING {
		this.stepTimer -= globalState.deltaTime
		if this.stepTimer <= 0.0 {
			// Check end-game state
			ended, success := this.level.checkEnded()
			if ended && success {
				return Response{RESPONSE_POP, nil}
			}
			this.level.step()
			this.stepTimer = secondsPerStep
		}
	}

	mx, my, _ := sdl.GetMouseState()
	if this.interaction == INTERACT_SETUP && this.transientTool == nil {
		// Check for tool dragging if we're not going and not already dragging
		dragged, row, column := this.level.canDragTool()
		if dragged != nil {
			dragged.removeFromLevel(&this.level, row)
			this.transientTool = dragged
			this.transientResetRow = row
			this.transientResetColumn = column
		}
	} else if this.transientTool != nil {
		// Deal with resetting the transient tool
		reset := false
		if this.interaction != INTERACT_SETUP {
			reset = true
		} else if globalState.leftClick {
			row, col, ok := this.level.inRangeOfToolSpot(this.transientTool, mx, my)
			if ok {
				this.transientTool.addToLevel(&this.level, row, col)
				this.transientTool = nil
			} else {
				reset = true
			}
		}
		if reset {
			this.transientTool.addToLevel(&this.level, this.transientResetRow, this.transientResetColumn)
			this.transientTool = nil
		}
	}

	return Response{RESPONSE_OK, nil}
}

func (this *GameState) render(renderer *sdl.Renderer) Response {
	renderer.SetDrawColor(0, 0, 0, 0xff)
	renderer.Clear()

	// Draw path lines
	//renderer.SetDrawColor(0xaa, 0xaa, 0xaa, 0xff)
	for i := range this.level.paths {
		rect := this.level.pathRect(i)
		{
			r := rect
			r.W = pathEndWidth
			renderer.Copy(this.assets["path-left"], nil, &r)
		}
		{
			r := rect
			r.X += pathEndWidth
			r.W -= pathEndWidth * 2
			renderer.Copy(this.assets["path"], nil, &r)
		}
		{
			r := rect
			r.X = r.X + (r.W - pathEndWidth)
			r.W = pathEndWidth
			renderer.Copy(this.assets["path-right"], nil, &r)
		}
		//renderer.FillRect(&rect)
	}

	// Draw stoppers
	for r, path := range this.level.paths {
		for _, stopper := range path.stoppers {
			if stopper.active {
				renderer.SetDrawColor(0xff, 0xff, 0xff, 0xff)
			} else {
				renderer.SetDrawColor(0xcc, 0xcc, 0xcc, 0xff)
			}
			rect := this.level.stopperRect(r, stopper.position)
			renderer.FillRect(&rect)
		}
	}

	// Draw non-swapper spots
	renderer.SetDrawColor(0xaa, 0x00, 0x00, 0xff)
	for r, path := range this.level.paths {
		if r == len(this.level.paths)-1 {
			if len(path.vertSwappers) > 0 {
				fatal("Somehow a nonSwapSpot got put on the bottom-most row!")
			}
			continue
		}
		for _, spot := range path.nonSwapSpots {
			rect := this.level.swapperRect(r, spot)
			rect.W += 10
			rect.X -= 5
			renderer.FillRect(&rect)
		}
	}

	// Draw swappers
	renderer.SetDrawColor(0xee, 0xee, 0xee, 0xff)
	for r, path := range this.level.paths {
		if r == len(this.level.paths)-1 {
			if len(path.vertSwappers) > 0 {
				fatal("Somehow a vertical swapper got put on the bottom-most row!")
			}
			continue
		}
		for _, swapper := range path.vertSwappers {
			rect := this.level.swapperRect(r, swapper.position)
			renderer.FillRect(&rect)
		}
	}

	// Draw flags
	for i, path := range this.level.paths {
		rect := this.level.baseRect(i, this.level.width-1)
		rect.X -= orbSize / 2
		rect.Y -= orbSize / 2
		rect.W = orbSize
		rect.H = orbSize
		renderer.Copy(this.assets["flag"], nil, &rect)

		fontTexture := fontRender(renderer, this.font, fmt.Sprintf("%d", path.flagIndex), sdl.Color{0xff, 0xff, 0xff, 0xff})
		defer fontTexture.Destroy()
		_, _, fW, fH, _ := fontTexture.Query()
		fontRect := centerRectInRect(sdl.Rect{0, 0, fW, fH}, rect)
		renderer.Copy(fontTexture, nil, &fontRect)
	}

	// Draw orbs
	for i, path := range this.level.paths {
		// Orb texture
		rect := this.level.baseRect(i, path.orbPosition)
		// Move to position along path
		// Offset to center
		rect.X -= orbSize / 2
		rect.Y -= orbSize / 2
		rect.W = orbSize
		rect.H = orbSize
		renderer.Copy(this.assets["orb"], nil, &rect)

		// Orb index
		// TODO(pixlark): Make this look less shite
		fontTexture := fontRender(renderer, this.font, fmt.Sprintf("%d", path.orbIndex), sdl.Color{0, 0, 0, 0xff})
		defer fontTexture.Destroy()

		_, _, fW, fH, _ := fontTexture.Query()
		fontRect := centerRectInRect(sdl.Rect{0, 0, fW, fH}, rect)
		renderer.Copy(fontTexture, nil, &fontRect)
	}

	// Pause button
	if this.interaction == INTERACT_GOING {
		if button(renderer, this.font, sdl.Rect{120, 0, 100, 50}, "Pause", white) {
			this.interaction = INTERACT_PAUSE
		}
	} else if this.interaction == INTERACT_PAUSE {
		if button(renderer, this.font, sdl.Rect{120, 0, 100, 50}, "Resume", white) {
			this.interaction = INTERACT_GOING
		}
	}

	// Go/Reset button
	var buttonText string
	if this.interaction != INTERACT_SETUP {
		buttonText = "Reset"
	} else {
		buttonText = "Go"
	}
	if button(renderer, this.font, sdl.Rect{0, 0, 100, 50}, buttonText, white) {
		switch this.interaction {
		case INTERACT_SETUP:
			this.interaction = INTERACT_GOING
			this.stepTimer = secondsPerStep
		default: // kind of a hack; why no fallthrough option, golang?
			this.interaction = INTERACT_SETUP
			this.level.reset()
		}
	}

	// Transient tool
	mx, my, _ := sdl.GetMouseState()
	if this.transientTool != nil {
		switch this.transientTool.(type) {
		case Stopper:
			renderer.SetDrawColor(0xff, 0xff, 0xff, 0xff)
			rect := sdl.Rect{
				mx - (stopperSize / 2), my - (stopperSize / 2),
				stopperSize, stopperSize,
			}
			renderer.FillRect(&rect)
		case VertSwapper:
			renderer.SetDrawColor(0xff, 0xff, 0xff, 0xff)
			rect := sdl.Rect{
				mx - (swapperWidth / 2), my - (swapperHeight / 2),
				swapperWidth, swapperHeight,
			}
			renderer.FillRect(&rect)
		}
	}

	return Response{RESPONSE_OK, nil}
}

func (this *GameState) exit() {

}
