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
	gridWidth      = 50.0
	gridHeight     = 50.0
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
	// FIXME: apple can currently spawn on top of snake
	minX := float64(random.Intn(gridWidth-1)) * widthInterval
	minY := float64(random.Intn(gridHeight-1)) * heightInterval
	return pixel.R(minX, minY, minX+widthInterval, minY+heightInterval)
}

func drawApple(apple pixel.Rect, imd *imdraw.IMDraw) {
	imd.Color = colornames.Red
	imd.Push(apple.Min)
	imd.Push(apple.Max)
	imd.Rectangle(0)
}

func moveApple(apple *pixel.Rect, random *rand.Rand) {
	// FIXME: apple can currently spawn on top of snake
	newX := float64(random.Intn(gridWidth-1)) * widthInterval
	newY := float64(random.Intn(gridHeight-1)) * heightInterval
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
			snake.Grow(2)
			score++
			moveApple(&apple, randomGen)
		}

		dt += time.Since(start)
	}
}

func main() {
	pixelgl.Run(run)
}
