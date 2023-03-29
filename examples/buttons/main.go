package main

import (
	"device/gba"
	"machine"

	"image/color"
	"runtime/interrupt"

	"tinygo.org/x/tinydraw"
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

	gba.DISP.DISPSTAT.SetBits(gba.DISPSTAT_VBLANK_IRQ_ENABLE << gba.DISPSTAT_VBLANK_IRQ_Pos)
	interrupt.New(machine.IRQ_VBLANK, update).Enable()

	clearScreen(black)

	for {
	}
}

func update(interrupt.Interrupt) {
	switch {
	case tinygba.ButtonStart.IsPushed():
		clearScreen(black)

	case tinygba.ButtonSelect.IsPushed():
		clearScreen(white)

	case tinygba.ButtonRight.IsPushed():
		clearScreen(green)

	case tinygba.ButtonLeft.IsPushed():
		clearScreen(red)

	case tinygba.ButtonDown.IsPushed():
		clearScreen(gBlue)

	case tinygba.ButtonUp.IsPushed():
		clearScreen(gRed)

	case tinygba.ButtonA.IsPushed():
		clearScreen(gYellow)

	case tinygba.ButtonB.IsPushed():
		clearScreen(gGreen)

	}
}

func clearScreen(c color.RGBA) {
	tinydraw.FilledRectangle(
		display,
		int16(0), int16(0),
		screenWidth, screenHeight,
		c,
	)
}
