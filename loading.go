package main

import rl "github.com/gen2brain/raylib-go/raylib"

var (
	defaultFont      rl.Font
	tilesAtlaces     = map[string]*spriteAtlas{}
	pawnsAtlaces     = map[string]*spriteAtlas{}
	furnitureAtlaces = map[string]*spriteAtlas{}

	uiAtlaces = map[string]*spriteAtlas{}
)

func loadResources() {
	// defaultFont = rl.LoadFont("resources/flexi.ttf")
	tilesAtlaces = make(map[string]*spriteAtlas)
	pawnsAtlaces = make(map[string]*spriteAtlas)
	furnitureAtlaces = make(map[string]*spriteAtlas)
	uiAtlaces = make(map[string]*spriteAtlas)

	pawnsAtlaces["player"] = CreateAtlasFromFile("assets/player.png", 0, 0, 16, 16, 16, 16, 1)
}
