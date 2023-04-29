package main

import (
	"machine"

	"image/color"

	"tinygo.org/x/tinygba"
)

var (
	display = machine.Display

	// Screen resolution
	screenWidth, screenHeight = display.Size()

	game = NewGame()
)

func main() {
	// Set up the display
	display.Configure()

	clearScreen()

	for {
		tinygba.WaitForVBlank()

		update()
	}
}

func update() {
	key := tinygba.ReadButtons()

	switch game.Status {
	case GameSplash:
		game.Splash()

		if tinygba.ButtonStart.IsPushed(key) {
			game.Start()
		}

	case GamePlay:
		switch {
		case tinygba.ButtonStart.IsPushed(key):
			game.Over()

		case tinygba.ButtonRight.IsPushed(key):
			game.Play(SnakeRight)

		case tinygba.ButtonLeft.IsPushed(key):
			game.Play(SnakeLeft)

		case tinygba.ButtonDown.IsPushed(key):
			game.Play(SnakeDown)

		case tinygba.ButtonUp.IsPushed(key):
			game.Play(SnakeUp)

		default:
			game.Play(-1)
		}

	case GameOver:
		game.Splash()

		if tinygba.ButtonStart.IsPushed(key) {
			game.Start()
		}
	}
}

func clearScreen() {
	tinygba.FillScreen(color.RGBA{})
}
