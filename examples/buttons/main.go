package main

import (
	"machine"

	"image/color"

	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/proggy"
	"tinygo.org/x/tinygba"
)

var (
	// Colors
	gbaPurple = color.RGBA{106, 55, 188, 255}
	white     = color.RGBA{255, 255, 255, 255}
	green     = color.RGBA{0, 255, 0, 255}
)

var (
	upX int16 = 45
	upY int16 = 70

	downX int16 = 40
	downY int16 = 90

	leftX int16 = 20
	leftY int16 = 80

	rightX int16 = 60
	rightY int16 = 80

	aX int16 = 200
	aY int16 = 70

	bX int16 = 180
	bY int16 = 80

	lX int16 = 30
	lY int16 = 30

	rX int16 = 200
	rY int16 = 30

	startX int16 = 60
	startY int16 = 120

	selectX int16 = 60
	selectY int16 = 130
)

func main() {
	machine.Display.Configure()

	tinygba.FillScreen(gbaPurple)

	tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, upX, upY, "Up", white)
	tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, downX, downY, "Down", white)
	tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, leftX, leftY, "Left", white)
	tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, rightX, rightY, "Right", white)
	tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, aX, aY, "A", white)
	tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, bX, bY, "B", white)
	tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, lX, lY, "L", white)
	tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, rX, rY, "R", white)
	tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, startX, startY, "Start", white)
	tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, selectX, selectY, "Select", white)

	for {
		tinygba.WaitForVBlank()
		update()
	}
}

func update() {
	key := tinygba.ReadButtons()
	if tinygba.ButtonUp.IsPushed(key) {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, upX, upY, "Up", green)
	} else {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, upX, upY, "Up", white)
	}

	if tinygba.ButtonDown.IsPushed(key) {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, downX, downY, "Down", green)
	} else {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, downX, downY, "Down", white)
	}

	if tinygba.ButtonLeft.IsPushed(key) {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, leftX, leftY, "Left", green)
	} else {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, leftX, leftY, "Left", white)
	}

	if tinygba.ButtonRight.IsPushed(key) {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, rightX, rightY, "Right", green)
	} else {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, rightX, rightY, "Right", white)
	}

	if tinygba.ButtonA.IsPushed(key) {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, aX, aY, "A", green)
	} else {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, aX, aY, "A", white)
	}

	if tinygba.ButtonB.IsPushed(key) {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, bX, bY, "B", green)
	} else {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, bX, bY, "B", white)
	}

	if tinygba.ButtonL.IsPushed(key) {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, lX, lY, "L", green)
	} else {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, lX, lY, "L", white)
	}

	if tinygba.ButtonR.IsPushed(key) {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, rX, rY, "R", green)
	} else {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, rX, rY, "R", white)
	}

	if tinygba.ButtonStart.IsPushed(key) {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, startX, startY, "Start", green)
	} else {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, startX, startY, "Start", white)
	}

	if tinygba.ButtonSelect.IsPushed(key) {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, selectX, selectY, "Select", green)
	} else {
		tinyfont.WriteLine(&machine.Display, &proggy.TinySZ8pt7b, selectX, selectY, "Select", white)
	}
}
