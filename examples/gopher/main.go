// this is basically https://gitlab.com/TitiMoby/tinygo-examples/-/blob/main/gba/03-Gopher/main.go
// with the newer tinygba API and some animation improvements.
package main

import (
	"machine"

	"image/color"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freesans"
	"tinygo.org/x/tinyfont/gophers"
	"tinygo.org/x/tinygba"
)

var (
	// Display from machine
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

	// Coordinates
	x int16 = 100 //TODO: horizontally center
	y int16 = 100 //TODO: vertically center

	border int16 = 40
)

func main() {
	// Set up the display
	display.Configure()

	// Display Gopher text message and draw our Gophers
	drawGophers()

	// Infinite loop to avoid exiting the application
	for {
		tinygba.WaitForVBlank()

		update()
	}
}

func drawGophers() {
	// Display a textual message "Gopher" with Google colors
	tinyfont.DrawChar(&display, &freesans.Bold24pt7b, 36, 60, 'G', gBlue)
	tinyfont.DrawChar(&display, &freesans.Bold24pt7b, 71, 60, 'o', gRed)
	tinyfont.DrawChar(&display, &freesans.Bold24pt7b, 98, 60, 'p', gYellow)
	tinyfont.DrawChar(&display, &freesans.Bold24pt7b, 126, 60, 'h', gGreen)
	tinyfont.DrawChar(&display, &freesans.Bold24pt7b, 154, 60, 'e', gBlue)
	tinyfont.DrawChar(&display, &freesans.Bold24pt7b, 180, 60, 'r', gRed)

	// Display a "press START button" message - center
	tinyfont.WriteLine(&display, &tinyfont.TomThumb, 85, 90, "Press START button", white)

	// Display two gophers
	tinyfont.DrawChar(&display, &gophers.Regular58pt, 5, 140, 'B', green)
	tinyfont.DrawChar(&display, &gophers.Regular58pt, 195, 140, 'X', red)
}

func update() {
	key := tinygba.ReadButtons()

	switch {
	case tinygba.ButtonStart.IsPushed(key):
		clearScreen()

		// Display gopher
		tinyfont.DrawChar(&display, &gophers.Regular58pt, x, y, 'B', green)
	case tinygba.ButtonSelect.IsPushed(key):
		clearScreen()
		drawGophers()

	case tinygba.ButtonRight.IsPushed(key):
		tinyfont.DrawChar(&display, &gophers.Regular58pt, x, y, 'B', black)

		if x <= screenWidth-border*2 {
			x = x + 10
		}

		// display gopher at right
		tinyfont.DrawChar(&display, &gophers.Regular58pt, x, y, 'B', green)

	case tinygba.ButtonLeft.IsPushed(key):
		tinyfont.DrawChar(&display, &gophers.Regular58pt, x, y, 'B', black)

		if x >= border {
			x = x - 10
		}

		// display gopher at right
		tinyfont.DrawChar(&display, &gophers.Regular58pt, x, y, 'B', green)

	case tinygba.ButtonDown.IsPushed(key):
		tinyfont.DrawChar(&display, &gophers.Regular58pt, x, y, 'B', black)

		if y <= screenHeight-border {
			y = y + 10
		}

		tinyfont.DrawChar(&display, &gophers.Regular58pt, x, y, 'B', green)

	case tinygba.ButtonUp.IsPushed(key):
		tinyfont.DrawChar(&display, &gophers.Regular58pt, x, y, 'B', black)

		if y >= border*2 {
			y = y - 10
		}

		tinyfont.DrawChar(&display, &gophers.Regular58pt, x, y, 'B', green)

	case tinygba.ButtonA.IsPushed(key):
		tinyfont.DrawChar(&display, &gophers.Regular58pt, x, y, 'B', black)

		// Display the gopher up
		y = y - 20
		tinyfont.DrawChar(&display, &gophers.Regular58pt, x, y, 'B', green)

		tinyfont.DrawChar(&display, &gophers.Regular58pt, x, y, 'B', black)
		// Display the gopher down
		y = y + 20
		tinyfont.DrawChar(&display, &gophers.Regular58pt, x, y, 'B', green)
	}
}

func clearScreen() {
	tinygba.FillScreen(black)
}
