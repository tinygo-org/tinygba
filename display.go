package tinygba

import (
	"machine"

	"image/color"
)

var screenWidth, screenHeight = machine.Display.Size()

// FillScreen fills the screen with a given color.
func FillScreen(c color.RGBA) {
	FillRectangle(0, 0, screenWidth, screenHeight, c)
}

// FillRectangle fills a rectangle at a given coordinates with a color.
func FillRectangle(x, y, width, height int16, c color.RGBA) error {
	for i := int16(0); i < height; i++ {
		for j := int16(0); j < width; j++ {
			machine.Display.SetPixel(x+j, y+i, c)
		}
	}

	return nil
}

// FillRectangleWithBuffer fills a rectangular area of the display with buffer.
func FillRectangleWithBuffer(x, y, width, height int16, buffer []color.RGBA) error {
	offset := 0
	for i := int16(0); i < height; i++ {
		for j := int16(0); j < width; j++ {
			machine.Display.SetPixel(x+j, y+i, buffer[offset])
			offset += 1
		}
	}

	return nil
}
