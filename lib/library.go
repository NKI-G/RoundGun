package library

import (
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

func DrawCicle(renderer *sdl.Renderer, x, y, size int32, r, g, b, a uint8){
	// 안티앨리어싱된 원의 경계 그리기
	gfx.AACircleRGBA(renderer, x, y, size, r, r, b, a)

	// 채워진 원 그리기
	gfx.FilledCircleRGBA(renderer, x-1, y-2, size-2, r, r, b, a)
	gfx.FilledCircleRGBA(renderer, x+1, y+2, size-2, r, r, b, a)
	gfx.FilledCircleRGBA(renderer, x+1, y-2, size-2, r, r, b, a)
	gfx.FilledCircleRGBA(renderer, x-1, y+2, size-2, r, r, b, a)
	gfx.FilledCircleRGBA(renderer, x, y, size-1, r, r, b, a)
}