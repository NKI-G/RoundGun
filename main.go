package main

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	// SDL 초기화
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %v\n", err)
		return
	}
	defer sdl.Quit()

	// 윈도우 생성
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create window: %v\n", err)
		return
	}
	defer window.Destroy()

	// 렌더러 생성
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %v\n", err)
		return
	}
	defer renderer.Destroy()

	// 배경색 설정
	bgColor := sdl.Color{R: 255, G: 255, B: 255, A: 255}
	renderer.SetDrawColor(bgColor.R, bgColor.G, bgColor.B, bgColor.A)
	renderer.Clear() // 배경색으로 화면을 채웁니다

	// 이미지 로드
	player, err := img.Load("./test.png")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load PNG: %v\n", err)
		return
	}
	defer player.Free()

	// 텍스처로 변환
	texture, err := renderer.CreateTextureFromSurface(player)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create texture: %v\n", err)
		return
	}
	defer texture.Destroy()

	// 메인 루프
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Quit")
				running = false
			default:
				// 다른 이벤트 처리
				_ = e
			}
		}

		// 화면 업데이트
		renderer.SetDrawColor(bgColor.R, bgColor.G, bgColor.B, bgColor.A)
		renderer.Clear()

		// SDL 렌더러를 통해 텍스처 복사
		renderer.Copy(texture, nil, &sdl.Rect{X: 100, Y: 100, W: 200, H: 200})

		// 안티앨리어싱된 원의 경계 그리기
		gfx.AACircleRGBA(renderer, 500, 200, 79, 0, 0, 255, 255)

		// 채워진 원 그리기
		gfx.FilledCircleRGBA(renderer, 500-1, 200-2, 77, 0, 0, 255, 255)
		gfx.FilledCircleRGBA(renderer, 500+1, 200+2, 77, 0, 0, 255, 255)
		gfx.FilledCircleRGBA(renderer, 500+1, 200-2, 77, 0, 0, 255, 255)
		gfx.FilledCircleRGBA(renderer, 500-1, 200+2, 77, 0, 0, 255, 255)
		gfx.FilledCircleRGBA(renderer, 500, 200, 78, 0, 0, 255, 255)

		// 화면에 반영
		renderer.Present()

		sdl.Delay(33) // 대략 30 FPS
	}
}
