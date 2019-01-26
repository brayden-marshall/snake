package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type Snake struct {
	Body      []pixel.Rect
	Direction Direction
	Alive     bool
}

func NewSnake(size int) Snake {
	snake := Snake{
		Direction: DirectionLeft,
		Alive:     true,
	}
	minX := (gridWidth / 2) * widthInterval
	minY := (gridHeight / 2) * heightInterval
	maxX := minX + widthInterval
	maxY := minY + widthInterval
	for i := 0; i < size; i++ {
		snake.Body = append(snake.Body, pixel.R(minX, minY, maxX, maxY))
		minX += widthInterval
		maxX += widthInterval
	}
	return snake
}

// returns true if the snake is in bounds of the window, else false
func (s *Snake) InBounds() bool {
	return s.Body[0].Min.X >= 0 && s.Body[0].Min.Y >= 0 &&
		s.Body[0].Max.X <= windowWidth && s.Body[0].Max.Y <= windowHeight
}

func (s *Snake) SelfCollide() bool {
	head := s.Body[0]
	// iterates through all body parts to see if head is touching any of them
	for _, cell := range s.Body[1:] {
		if head.Min == cell.Min {
			return true
		}
	}
	return false
}

func (s *Snake) Move() {
	// creating a copy of s.body
	var oldBody = make([]pixel.Rect, len(s.Body))
	copy(oldBody, s.Body)
	// shifting all body pieces one to the right
	for i := 0; i < len(s.Body)-1; i++ {
		s.Body[i+1] = oldBody[i]
	}

	var offset pixel.Vec
	switch s.Direction {
	case DirectionUp:
		offset = pixel.V(0, heightInterval)
	case DirectionLeft:
		offset = pixel.V(-widthInterval, 0)
	case DirectionDown:
		offset = pixel.V(0, -heightInterval)
	case DirectionRight:
		offset = pixel.V(widthInterval, 0)
	}

	// moving the snake's head based on an offset(direction)
	s.Body[0] = s.Body[0].Moved(offset)

	if !s.InBounds() {
		s.Die()
	}

	if s.SelfCollide() {
		s.Die()
	}
}

func (s *Snake) Grow(n int) {
	// adding another element to the end of the snake's body
	s.Body = append(s.Body, s.Body[len(s.Body)-1])
	if n == 1 {
		return
	}
	s.Grow(n - 1)
}

func (s *Snake) Die() {
	s.Alive = false
}

func (s *Snake) Draw(imd *imdraw.IMDraw) {
	for i := range s.Body {
		imd.Color = colornames.Grey
		imd.Push(s.Body[i].Min)
		imd.Push(s.Body[i].Max)
		imd.Rectangle(0)

		// drawing thin black outline to distinguish seperate squares
		imd.Color = colornames.Black
		imd.Push(s.Body[i].Min)
		imd.Push(s.Body[i].Max)
		imd.Rectangle(1)
	}
}
