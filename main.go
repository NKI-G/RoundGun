package main

// 이코드의 저자는 뉴키니와 챗쥐피티입니다.
import (
	"fmt"
	"math"
	"os"

	library "roundgun/lib"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Player struct {
	Sprite *Sprite
	Hp int
	IsHit bool
	HitCount int
}

type Built struct {
	Sprite *Sprite
	Dir float64
	Count int
	Damage int
}

type Monstor struct {
	Sprite *Sprite
	Hp int
	IsHit bool
	HitCount int
	Level int
}

func monstorUpdate(monstor *Monstor, player Sprite) {
	// Update monster's color or HP based on level
	if (monstor.IsHit == false){
		if monstor.Level == 1 {
			monstor.Sprite.Color = [4]int{155, 100, 100, 255}
		} else {
			monstor.Hp = 0
		}
	}
	// Calculate speed based on monster's level
	var speed int
	if monstor.Level < 18 {
		speed = monstor.Level + 2
	} else {
		speed = 20
	}

	// Calculate direction towards the player
	dir := math.Atan2(float64(player.DfPos[1])-float64(monstor.Sprite.DfPos[1]), float64(player.DfPos[0])-float64(monstor.Sprite.DfPos[0]))

	// Update monster's position based on direction and speed
	monstor.Sprite.DfPos[0] += int(float64(speed) * math.Cos(dir))
	monstor.Sprite.DfPos[1] += int(float64(speed) * math.Sin(dir))

	if (monstor.IsHit && monstor.HitCount <= 30){
		monstor.Sprite.Color = [4]int{255, 0, 0, 255}
		monstor.HitCount++
	}else{
		monstor.HitCount = 0
		monstor.IsHit = false
		println(monstor.IsHit)
	}
}

func playerUpdate(player *Player) {
	if (player.IsHit && player.HitCount <= 30){
		player.Sprite.Color = [4]int{255, 0, 0, 255}
		player.HitCount++
	}else{
		player.HitCount = 0
		player.IsHit = false
		println(player.IsHit)
	}
}

type Sprite struct {
	Color [4]int
	Size  int
	DfPos [2]int // 기본 좌표. 하지만 그냥 좌표가 됨
}

func newSprite(color [4]int, size int, dfPos [2]int) *Sprite {
	return &Sprite{Color: color, Size: size, DfPos: dfPos}
}

func spriteMove(s *Sprite, pos [2]int) {
	s.DfPos = pos
}

func spriteDraw(s *Sprite, renderer *sdl.Renderer, cameraPos [2]int) {
	library.DrawCicle(renderer, int32(s.DfPos[0])-int32(cameraPos[0]), int32(s.DfPos[1])-int32(cameraPos[1]), int32(s.Size), uint8(s.Color[0]), uint8(s.Color[1]), uint8(s.Color[2]), uint8(s.Color[3]))
}

func main() {
	// SDL 초기화
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDL: %v\n", err)
		return
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize TTF: %v\n", err)
		return
	}
	defer ttf.Quit()

	windowsize := [2]int{800, 500}

	// 윈도우 생성
	window, err := sdl.CreateWindow("RoundGun", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(windowsize[0]), int32(windowsize[1]), sdl.WINDOW_SHOWN)
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
	
	font, err := ttf.OpenFont("assets/SpoqaHanSansNeo-Bold.ttf", 24)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %v\n", err)
		return
	}
	defer font.Close()

	// 배경색 설정
	bgColor := sdl.Color{R: 0, G: 0, B: 0, A: 255}
	renderer.SetDrawColor(bgColor.R, bgColor.G, bgColor.B, bgColor.A)
	renderer.Clear()

	// 변수들
	clickPos := [2]int{0, 0}
	isClick := false
	clickRadius := 40
	clickAlpha := 255
	stepCount := 4
	step := 0
	cameraPos := [2]int{0, 0}
	player := Player{Sprite: newSprite([4]int{255, 255, 255, 255}, 40, [2]int{400, 250}), Hp: 20} // 플레이어 초기 위치 중앙
	mousePos := [2]int{0, 0}
	bullets := []Built{}
	var deltaCountBullets float32 = 0
	monstors := []*Monstor{}
	var deltaCountMonstors float32 = 0
	var deltaDamageCount float32 = 0



	lastTime := sdl.GetTicks64()

	// 메인 루프
	running := true
	for running {
		// 현재 프레임 시간
		currentTime := sdl.GetTicks64()
		// DeltaTime 계산 (초 단위로 변환)
		deltaTime := float32(currentTime-lastTime) / 1000.0
		lastTime = currentTime

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("Quit")
				running = false
			case *sdl.MouseButtonEvent:
				if e.State == sdl.PRESSED {
					isClick = true
					clickRadius = 30 // 클릭 시 반지름 초기화
					clickAlpha = 255 // 클릭 시 투명도 초기화
					// 클릭한 화면상의 좌표를 월드 좌표로 변환
					clickPos[0] = int(e.X) + cameraPos[0]
					clickPos[1] = int(e.Y) + cameraPos[1]
					step = 0 // 이동 단계를 초기화
				}
			case *sdl.MouseMotionEvent:
				mousePos[0] = int(e.X)
				mousePos[1] = int(e.Y)
			case *sdl.KeyboardEvent:
				if e.State == sdl.PRESSED && e.Keysym.Sym == sdl.K_UP {
					cameraPos[1] += 50
				}
				if e.State == sdl.PRESSED && e.Keysym.Sym == sdl.K_DOWN {
					cameraPos[1] -= 50
				}
				if e.State == sdl.PRESSED && e.Keysym.Sym == sdl.K_LEFT {
					cameraPos[0] -= 50
				}
				if e.State == sdl.PRESSED && e.Keysym.Sym == sdl.K_RIGHT {
					cameraPos[0] += 50
				}
			default:
				_ = e
			}
		}

		// 화면 업데이트
		renderer.SetDrawColor(bgColor.R, bgColor.G, bgColor.B, bgColor.A)
		renderer.Clear()

		// 플레이어 그리기
		spriteDraw(player.Sprite, renderer, cameraPos)
		playerUpdate(&player)

		// 클릭 위치에 원 그리기
		if isClick {
			// 4단계로 나누어 이동
			if step < stepCount {
				fraction := (float64(step+1) / float64(stepCount))
				player.Sprite.DfPos[0] = int(float64(player.Sprite.DfPos[0]) + fraction*float64(clickPos[0]-player.Sprite.DfPos[0]))
				player.Sprite.DfPos[1] = int(float64(player.Sprite.DfPos[1]) + fraction*float64(clickPos[1]-player.Sprite.DfPos[1]))
				step++
			} else {
				isClick = false // 이동이 완료되면 클릭 상태를 해제
			}

			// 클릭 원 그리기
			alpha := uint8(clickAlpha)
			radius := clickRadius
			gfx.AACircleRGBA(renderer, int32(clickPos[0])-int32(cameraPos[0]), int32(clickPos[1])-int32(cameraPos[1]), int32(radius), 255, 255, 255, alpha)

			clickRadius -= 10 // 원의 크기 감소
			clickAlpha -= 80  // 투명도 감소
			if clickRadius <= 0 {
				clickRadius = 0
			}
		}

		// 마우스의 월드 좌표 구하기
		mouseWorldPos := [2]int{mousePos[0] + cameraPos[0], mousePos[1] + cameraPos[1]}

		// 마우스와 플레이어 사이의 각도 계산
		mouseDir := math.Atan2(float64(mouseWorldPos[1])-float64(player.Sprite.DfPos[1]), float64(mouseWorldPos[0])-float64(player.Sprite.DfPos[0]))

		// 몬스터 생성
		deltaCountMonstors += deltaTime
		if deltaCountMonstors >= 10 {
			X, Y := library.MonstorPosCreator(cameraPos, windowsize)
			monstors = append(monstors, &Monstor{Sprite: newSprite([4]int{100, 100, 100, 255}, 30, [2]int{X, Y}), Hp: 20, IsHit: false, Level: 1, HitCount: 0})
			deltaCountMonstors = 0
		}

		if (player.Hp <= 0){
			running = false
		}

		// 몬스터 처리
		for i := 0; i < len(monstors); i++ {
			if monstors[i].Hp <= 0 {
				monstors = append(monstors[:i], monstors[i+1:]...)
				i-- // 요소를 제거했으므로 인덱스를 감소시킴
				continue
			}
			monstorUpdate(monstors[i], *player.Sprite)
			spriteDraw(monstors[i].Sprite, renderer, cameraPos)
			if (library.IsCollision((*library.Sprite)(monstors[i].Sprite), (*library.Sprite)(player.Sprite))){
				if (deltaDamageCount >= 3){
					player.Hp -= 1
					deltaDamageCount = 0
				}
			}
		}

		deltaDamageCount++

		// 총알 생성
		deltaCountBullets += deltaTime
		bulletSpeed := 5.0 // 총알 속도
		if deltaCountBullets >= 1 {
			bullets = append(bullets, Built{Sprite: newSprite([4]int{100, 100, 100, 255}, 10, player.Sprite.DfPos), Dir: mouseDir, Damage: 3})
			deltaCountBullets = 0
		}

		// 총알 처리
		for i := 0; i < len(bullets); i++ {
			spriteDraw(bullets[i].Sprite, renderer, cameraPos)

			// 총알 이동
			bullets[i].Sprite.DfPos[0] += int(bulletSpeed * math.Cos(bullets[i].Dir))
			bullets[i].Sprite.DfPos[1] += int(bulletSpeed * math.Sin(bullets[i].Dir))
			bullets[i].Count++

			// 총알 제거 (거리가 너무 멀어지면)
			if bullets[i].Count >= 200 {
				bullets = append(bullets[:i], bullets[i+1:]...)
				i-- // 요소를 제거했으므로 인덱스를 감소시킴
				continue
			}

			// 몬스터와 충돌 감지
			for j := 0; j < len(monstors); j++ {
				if library.IsCollision((*library.Sprite)(bullets[i].Sprite), (*library.Sprite)(monstors[j].Sprite)) {
					monstors[j].Hp -= bullets[i].Damage
					monstors[j].IsHit = true
					bullets = append(bullets[:i], bullets[i+1:]...)
					i-- // 요소를 제거했으므로 인덱스를 감소시킴
					break
				}
			}
		}

		// 카메라가 플레이어를 부드럽게 따라가게 함
		cameraSpeed := 0.1 // 카메라가 플레이어를 따라가는 속도 (0.0 ~ 1.0 사이의 값)
		cameraPos[0] = int(float64(cameraPos[0]) + cameraSpeed*float64(player.Sprite.DfPos[0]-cameraPos[0]-windowsize[0]/2))
		cameraPos[1] = int(float64(cameraPos[1]) + cameraSpeed*float64(player.Sprite.DfPos[1]-cameraPos[1]-windowsize[1]/2))

		// 나머지 ui
		text := fmt.Sprintf("Frame:")
		_, err = library.RenderText(font, renderer, text, sdl.Color{255,255,255,255}, 10, 10)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to render text: %v\n", err)
			return
		}

		text = fmt.Sprintf("Frame:")
		_, err = library.RenderText(font, renderer, text, sdl.Color{255,255,255,255}, 10, 34)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to render text: %v\n", err)
			return
		}

		// 화면에 반영
		renderer.Present()

		sdl.Delay(33) // 대략 30 FPS
	}
	
	if (player.Hp <= 0){
		for i := 0; i < 33*10; i++ {
			renderer.SetDrawColor(0, 0, 0, 255)
			renderer.Clear()
	
			text := fmt.Sprintf("Game Over")
			_, err := library.RenderText(font, renderer, text, sdl.Color{255,255,255,255}, 50, 50)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to render text: %v\n", err)
				return
			}
	
			text = fmt.Sprintf("Frame: %d", i/33)
			_, err = library.RenderText(font, renderer, text, sdl.Color{255,255,255,255}, 50, 74)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to render text: %v\n", err)
				return
			}
	
			renderer.Present()
		}
	}
}
