package library

import (
	"math"
	"math/rand"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Sprite struct {
	Color [4]int
	Size  int
	DfPos [2]int // 기본 좌표. 하지만 그냥 좌표가 됨
}

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

func MonstorPosCreator(cameraPos [2]int, windowsize [2]int) (int, int) {
	// 화면 경계 계산
	leftBoundary := cameraPos[0] - 300
	rightBoundary := cameraPos[0] + windowsize[0] + 300
	topBoundary := cameraPos[1] - 300
	bottomBoundary := cameraPos[1] + windowsize[1] + 300

	var spawnX, spawnY int

	for {
		// 무작위로 몬스터의 스폰 위치 선택
		spawnX = rand.Intn(rightBoundary-leftBoundary) + leftBoundary
		spawnY = rand.Intn(bottomBoundary-topBoundary) + topBoundary

		// 스폰된 위치가 화면 밖에 있는지 확인
		if spawnX < cameraPos[0] || spawnX > cameraPos[0]+windowsize[0] || spawnY < cameraPos[1] || spawnY > cameraPos[1]+windowsize[1] {
			break // 조건에 맞는 위치가 나왔으므로 반복문 탈출
		}
	}

	return spawnX, spawnY
}

func IsCollision(s1, s2 *Sprite) bool {
	// 두 스프라이트의 중심 좌표와 반지름을 계산합니다.
	distanceX := s1.DfPos[0] - s2.DfPos[0]
	distanceY := s1.DfPos[1] - s2.DfPos[1]
	distance := math.Sqrt(float64(distanceX*distanceX + distanceY*distanceY))

	// 두 스프라이트의 반지름을 합친 값과 거리 비교
	return distance <= float64(s1.Size+s2.Size)
}

func RenderText(font *ttf.Font, renderer *sdl.Renderer, text string, color sdl.Color, x, y int32) (*sdl.Texture, error) {
	surface, err := font.RenderUTF8Blended(text, color)
	if err != nil {
		return nil, err
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, err
	}

	rect := sdl.Rect{X: x, Y: y, W: surface.W, H: surface.H}
	err = renderer.Copy(texture, nil, &rect)
	if err != nil {
		return nil, err
	}

	return texture, nil
}