package main

import (	
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var _ = fmt.Println

const (
	pathLeftPad   int32 = 40
	pathTopPad    int32 = 100
	pathTop       int32 = pathTopPad
	pathThickness int32 = 5
	pathVertSpace int32 = screenH / 6
	pathLeft      int32 = pathLeftPad
	pathRight     int32 = screenW - pathLeftPad

	orbSize int32 = 40
)

func loadTexture(renderer *sdl.Renderer, path string) *sdl.Texture {
	texture, err := img.LoadTexture(renderer, "orb.png")
	if err != nil {
		fatal("Could not open texture")
	}
	return texture
}

func loadFont(path string, size int) *ttf.Font {
	font, err := ttf.OpenFont(path, size)
	if err != nil {
		fatal("Could not open font")
	}
	return font
}

type GameState struct {
	level  Level
	assets map[string]*sdl.Texture
	font   *ttf.Font
	going  bool
}

func (this *GameState) Init(renderer *sdl.Renderer) {
	this.assets = make(map[string]*sdl.Texture)
	this.assets["orb"] = loadTexture(renderer, "orb.png")

	this.font = loadFont("DejaVuSans.ttf", 15)
	this.going = false
	
	this.level.Init()
}

func (this *GameState) Update(events []sdl.Event) Response {
	for _, event := range events {
		switch event := event.(type) {
		case *sdl.KeyboardEvent:
			switch event.Type {
			case sdl.KEYDOWN:
				break
			}
		}
	}

	return Response{RESPONSE_OK, nil}
}

func fontRender(renderer *sdl.Renderer, font *ttf.Font, text string) *sdl.Texture {
	surface, _ := font.RenderUTF8Solid(text, sdl.Color{0, 0, 0, 0xff})
	defer surface.Free()
	texture,_ := renderer.CreateTextureFromSurface(surface)
	return texture
}

func (this *GameState) Render(renderer *sdl.Renderer) Response {
	renderer.SetDrawColor(0, 0, 0, 0xff)
	renderer.Clear()

	// Draw path lines
	renderer.SetDrawColor(0xff, 0xff, 0xff, 0xff)
	for i, _ := range this.level.paths {
		rect := pathRect(i)
		renderer.FillRect(&rect)
	}
	
	// Draw orbs
	for i, path := range this.level.paths {
		// Orb texture
		rect := pathRect(i)
		// Move to position along path
		rect.X += int32(path.orbPosition) * (rect.W / int32(this.level.width - 1))
		// Offset to center
		rect.X -= orbSize / 2
		rect.Y -= orbSize / 2
		rect.W = orbSize
		rect.H = orbSize
		renderer.Copy(this.assets["orb"], nil, &rect)

		// Orb index
		// TODO(pixlark): Make this look less shite
		fontTexture := fontRender(renderer, this.font, fmt.Sprintf("%d", path.orbIndex))
		defer fontTexture.Destroy()
		renderer.Copy(fontTexture, nil, &rect)
	}

	// Go/Reset button
	if Button(renderer, this.font, sdl.Rect{0, 0, 100, 50}, "Play") {
		fmt.Println("HASDLFKJWEF")
	}
	
	return Response{RESPONSE_OK, nil}
}

func (this *GameState) Exit() {
	
}
