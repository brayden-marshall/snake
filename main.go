package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math/rand"
	"time"
)

const (
	windowWidth    = 800
	windowHeight   = 800
	gridWidth      = 40.0
	gridHeight     = 40.0
	widthInterval  = float64(windowWidth / gridWidth)
	heightInterval = float64(windowHeight / gridHeight)
)

type Direction int

const (
	DirectionUp Direction = iota
	DirectionLeft
	DirectionDown
	DirectionRight
	DirectionMax
)

func newApple(random *rand.Rand) pixel.Rect {
	minX := 10 * widthInterval
	minY := 10 * heightInterval
	return pixel.R(minX, minY, minX+widthInterval, minY+heightInterval)
}

func drawApple(apple pixel.Rect, imd *imdraw.IMDraw) {
	imd.Color = colornames.Red
	imd.Push(apple.Min)
	imd.Push(apple.Max)
	imd.Rectangle(0)
}

func moveApple(apple *pixel.Rect, random *rand.Rand, snake Snake) {
	// two sets containing information about all invalid locations to move the apple
	invalidX := make(map[int]bool)
	invalidY := make(map[int]bool)
	for _, cell := range snake.Body {
		invalidX[int(cell.Min.X)] = true
		invalidY[int(cell.Min.Y)] = true
	}

	// finding all valid positions that we can move the apple
	var validPositions []pixel.Vec
	for i := 0.0; i < gridWidth; i++ {
		for j := 0.0; j < gridHeight; j++ {
			x := i * widthInterval
			y := j * heightInterval
			if !invalidX[int(x)] || !invalidY[int(y)] {
				validPositions = append(validPositions, pixel.V(x, y))
			}
		}
	}

	// picking a random spot to move the apple
	newPosition := validPositions[random.Intn(len(validPositions))]
	newX := newPosition.X
	newY := newPosition.Y
	*apple = pixel.R(newX, newY, newX+widthInterval, newY+heightInterval)
}

func run() {
	// creating the window
	cfg := pixelgl.WindowConfig{
		Title:  "Snake",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// imd is our drawing target
	imd := imdraw.New(nil)

	// creating the snake with length 5
	snake := NewSnake(5)

	// setting up random source for apple movement
	randSource := rand.NewSource(time.Now().UnixNano())
	randomGen := rand.New(randSource)

	// creating the apple
	apple := newApple(randomGen)

	var start time.Time
	var dt time.Duration
	var score int = len(snake.Body)

	for !win.Closed() && snake.Alive {
		start = time.Now()
		// clearing and drawing the frame
		win.Clear(colornames.Black)
		imd.Clear()
		snake.Draw(imd)
		drawApple(apple, imd)
		imd.Draw(win)
		win.Update()

		// listening for directional input
		if win.Pressed(pixelgl.KeyUp) {
			snake.Direction = DirectionUp
		} else if win.Pressed(pixelgl.KeyLeft) {
			snake.Direction = DirectionLeft
		} else if win.Pressed(pixelgl.KeyDown) {
			snake.Direction = DirectionDown
		} else if win.Pressed(pixelgl.KeyRight) {
			snake.Direction = DirectionRight
		}

		// moving the snake only once every 50 milliseconds
		if dt > 50*time.Millisecond {
			snake.Move()
			dt = 0
		}

		// if snake has eaten apple
		if snake.Body[0].Min == apple.Min {
			snake.Grow(10)
			score++
			moveApple(&apple, randomGen, snake)
		}

		dt += time.Since(start)
	}
}

func main() {
	pixelgl.Run(run)
}
