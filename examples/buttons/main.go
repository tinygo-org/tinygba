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

	// Colors
	black = color.RGBA{}
	white = color.RGBA{255, 255, 255, 255}
	green = color.RGBA{0, 255, 0, 255}
	red   = color.RGBA{255, 0, 0, 255}

	// Google colors
	gBlue   = color.RGBA{66, 163, 244, 255}
	gRed    = color.RGBA{219, 68, 55, 255}
	gYellow = color.RGBA{244, 160, 0, 255}
	gGreen  = color.RGBA{15, 157, 88, 255}
)

func main() {
	// Set up the display
	display.Configure()

	clearScreen(black)

	for {
		tinygba.WaitForVBlank()

		update()
	}
}

func update() {
	key := tinygba.ReadButtons()

	switch {
	case tinygba.ButtonStart.IsPushed(key):
		clearScreen(black)

	case tinygba.ButtonSelect.IsPushed(key):
		clearScreen(white)

	case tinygba.ButtonRight.IsPushed(key):
		clearScreen(green)

	case tinygba.ButtonLeft.IsPushed(key):
		clearScreen(red)

	case tinygba.ButtonDown.IsPushed(key):
		clearScreen(gBlue)

	case tinygba.ButtonUp.IsPushed(key):
		clearScreen(gRed)

	case tinygba.ButtonA.IsPushed(key):
		clearScreen(gYellow)

	case tinygba.ButtonB.IsPushed(key):
		clearScreen(gGreen)

	}
}

func clearScreen(c color.RGBA) {
	tinygba.FillScreen(c)
}
