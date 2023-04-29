package main

import (
	"image/color"
	"math/rand"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freesans"
	"tinygo.org/x/tinygba"
)

const (
	BCK = iota
	SNAKE
	APPLE
	TEXT
)

const (
	GameSplash = iota
	GameStart
	GamePlay
	GameOver
)

const (
	SnakeUp = iota
	SnakeDown
	SnakeLeft
	SnakeRight
)

const (
	WIDTHBLOCKS  = 24
	HEIGHTBLOCKS = 16
)

var colors = []color.RGBA{
	color.RGBA{0, 0, 0, 255},
	color.RGBA{255, 255, 255, 255},
	color.RGBA{255, 0, 0, 255},
}

var (
	// Those variable are there for a more easy reading of the apple shape.
	re = colors[APPLE] // red
	bk = colors[BCK]   // background
	gr = colors[SNAKE] // green

	// The array is split for a visual purpose too.
	appleBuf = []color.RGBA{
		bk, bk, bk, bk, bk, gr, gr, gr, bk, bk,
		bk, bk, bk, bk, gr, gr, gr, bk, bk, bk,
		bk, bk, bk, re, gr, gr, re, bk, bk, bk,
		bk, bk, re, re, re, re, re, re, bk, bk,
		bk, re, re, re, re, re, re, re, re, bk,
		bk, re, re, re, re, re, re, re, re, bk,
		bk, re, re, re, re, re, re, re, re, bk,
		bk, bk, re, re, re, re, re, re, bk, bk,
		bk, bk, bk, re, re, re, re, bk, bk, bk,
		bk, bk, bk, bk, bk, bk, bk, bk, bk, bk,
	}
)

type Snake struct {
	body      [104][2]int16
	length    int16
	direction int16
}

type Game struct {
	colors         []color.RGBA
	snake          Snake
	appleX, appleY int16
	Status         uint8
	score          int
}

var splashed = false
var scoreStr = []byte("SCORE: 123")

func NewGame() *Game {
	return &Game{
		colors: []color.RGBA{
			color.RGBA{0, 0, 0, 255},
			color.RGBA{0, 200, 0, 255},
			color.RGBA{250, 0, 0, 255},
			color.RGBA{160, 160, 160, 255},
		},
		snake: Snake{
			body: [104][2]int16{
				{0, 3},
				{0, 2},
				{0, 1},
			},
			length:    3,
			direction: SnakeLeft,
		},
		appleX: 5,
		appleY: 5,
		Status: GameSplash,
	}
}

func (g *Game) Splash() {
	if !splashed {
		g.splash()
		splashed = true
	}
}

func (g *Game) Start() {
	clearScreen()

	g.initSnake()
	g.drawSnake()
	g.createApple()

	g.Status = GamePlay
}

func (g *Game) Play(direction int) {
	switch direction {
	case SnakeLeft:
		if g.snake.direction != SnakeLeft {
			g.snake.direction = SnakeLeft
		}
	case SnakeRight:
		if g.snake.direction != SnakeRight {
			g.snake.direction = SnakeRight
		}
	case SnakeUp:
		if g.snake.direction != SnakeUp {
			g.snake.direction = SnakeUp
		}
	case SnakeDown:
		if g.snake.direction != SnakeDown {
			g.snake.direction = SnakeDown
		}
	}

	g.moveSnake()
}

func (g *Game) Over() {
	splashed = false

	g.Status = GameSplash
}

func (g *Game) splash() {
	clearScreen()

	tinyfont.WriteLine(&display, &freesans.Bold12pt7b, 10, 50, "SNAKE", g.colors[TEXT])
	tinyfont.WriteLine(&display, &freesans.Regular12pt7b, 18, 100, "Press START", g.colors[TEXT])

	if g.score > 0 {
		scoreStr[7] = 48 + uint8((g.score)/100)
		scoreStr[8] = 48 + uint8(((g.score)/10)%10)
		scoreStr[9] = 48 + uint8((g.score)%10)

		tinyfont.WriteLine(&display, &tinyfont.TomThumb, 50, 120, string(scoreStr), g.colors[TEXT])
	}
}

func (g *Game) initSnake() {
	g.snake.body[0][0] = 0
	g.snake.body[0][1] = 3
	g.snake.body[1][0] = 0
	g.snake.body[1][1] = 2
	g.snake.body[2][0] = 0
	g.snake.body[2][1] = 1

	g.snake.length = 3
	g.snake.direction = SnakeRight
}

func (g *Game) collisionWithSnake(x, y int16) bool {
	for i := int16(0); i < g.snake.length; i++ {
		if x == g.snake.body[i][0] && y == g.snake.body[i][1] {
			return true
		}
	}
	return false
}

func (g *Game) createApple() {
	g.appleX = int16(rand.Int31n(16))
	g.appleY = int16(rand.Int31n(13))
	for g.collisionWithSnake(g.appleX, g.appleY) {
		g.appleX = int16(rand.Int31n(16))
		g.appleY = int16(rand.Int31n(13))
	}
	g.drawApple(g.appleX, g.appleY)
}

func (g *Game) moveSnake() {
	x := g.snake.body[0][0]
	y := g.snake.body[0][1]

	switch g.snake.direction {
	case SnakeLeft:
		x--
		break
	case SnakeUp:
		y--
		break
	case SnakeDown:
		y++
		break
	case SnakeRight:
		x++
		break
	}
	if x >= WIDTHBLOCKS {
		x = 0
	}
	if x < 0 {
		x = WIDTHBLOCKS - 1
	}
	if y >= HEIGHTBLOCKS {
		y = 0
	}
	if y < 0 {
		y = HEIGHTBLOCKS - 1
	}

	if g.collisionWithSnake(x, y) {
		g.score = int(g.snake.length - 3)
		g.Status = GameOver
	}

	// draw head
	g.drawSnakePartial(x, y, g.colors[SNAKE])
	if x == g.appleX && y == g.appleY {
		g.snake.length++
		g.createApple()
	} else {
		// remove tail
		g.drawSnakePartial(g.snake.body[g.snake.length-1][0], g.snake.body[g.snake.length-1][1], g.colors[BCK])
	}
	for i := g.snake.length - 1; i > 0; i-- {
		g.snake.body[i][0] = g.snake.body[i-1][0]
		g.snake.body[i][1] = g.snake.body[i-1][1]
	}
	g.snake.body[0][0] = x
	g.snake.body[0][1] = y
}

func (g *Game) drawApple(x, y int16) {
	tinygba.FillRectangleWithBuffer(10*x, 10*y, 10, 10, appleBuf)
}

func (g *Game) drawSnake() {
	for i := int16(0); i < g.snake.length; i++ {
		g.drawSnakePartial(g.snake.body[i][0], g.snake.body[i][1], g.colors[SNAKE])
	}
}

func (g *Game) drawSnakePartial(x, y int16, c color.RGBA) {
	modY := int16(9)
	if y == 12 {
		modY = 8
	}
	tinygba.FillRectangle(10*x, 10*y, 9, modY, c)
}
