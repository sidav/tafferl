//go:build raylib

package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"image"
	"image/png"
	"os"
)

type spriteAtlas struct {
	light []rl.Texture2D
	dark  []rl.Texture2D
}

func CreateAtlasFromFile(filename string, topleftx, toplefty, originalSpriteW, originalSpriteH,
	desiredSpriteW, desiredSpriteH, totalFrames int) *spriteAtlas {
	const SPRITE_SCALE_FACTOR = 1

	newAtlas := spriteAtlas{
		// spriteSize: desiredSpriteSize * int(SPRITE_SCALE_FACTOR),
	}

	newAtlas.light = make([]rl.Texture2D, 1)
	newAtlas.dark = make([]rl.Texture2D, 1)

	file, _ := os.Open(filename)
	img, _ := png.Decode(file)
	file.Close()

	currPic := extractSubimageFromImage(img, topleftx+0*originalSpriteW, toplefty, originalSpriteW, originalSpriteH)
	rlImg := rl.NewImageFromImage(currPic)
	rl.ImageResizeNN(rlImg, int32(desiredSpriteW)*int32(SPRITE_SCALE_FACTOR), int32(desiredSpriteH)*int32(SPRITE_SCALE_FACTOR))
	newAtlas.light = append(newAtlas.light, rl.LoadTextureFromImage(rlImg))

	currPic = extractSubimageFromImage(img, topleftx+0*originalSpriteW, toplefty+originalSpriteH, originalSpriteW, originalSpriteH)
	rlImg = rl.NewImageFromImage(currPic)
	rl.ImageResizeNN(rlImg, int32(desiredSpriteW)*int32(SPRITE_SCALE_FACTOR), int32(desiredSpriteH)*int32(SPRITE_SCALE_FACTOR))
	newAtlas.dark = append(newAtlas.light, rl.LoadTextureFromImage(rlImg))

	return &newAtlas
}

func extractSubimageFromImage(img image.Image, fromx, fromy, w, h int) image.Image {
	minx, miny := img.Bounds().Min.X, img.Bounds().Min.Y
	//maxx, maxy := img.Bounds().Min.X, img.Bounds().Max.Y
	switch img.(type) {
	case *image.RGBA:
		subImg := img.(*image.RGBA).SubImage(
			image.Rect(minx+fromx, miny+fromy, minx+fromx+w, miny+fromy+h),
		)
		// reset img bounds, because RayLib goes nuts about it otherwise
		subImg.(*image.RGBA).Rect = image.Rect(0, 0, w, h)
		return subImg
	case *image.NRGBA:
		subImg := img.(*image.NRGBA).SubImage(
			image.Rect(minx+fromx, miny+fromy, minx+fromx+w, miny+fromy+h),
		)
		// reset img bounds, because RayLib goes nuts about it otherwise
		subImg.(*image.NRGBA).Rect = image.Rect(0, 0, w, h)
		return subImg
	case *image.Paletted:
		subImg := img.(*image.Paletted).SubImage(
			image.Rect(minx+fromx, miny+fromy, minx+fromx+w, miny+fromy+h),
		)
		// reset img bounds, because RayLib goes nuts about it otherwise
		subImg.(*image.Paletted).Rect = image.Rect(0, 0, w, h)
		return subImg
	default:
	}
	panic(fmt.Sprintf("\nUnknown image type %T", img))
}
