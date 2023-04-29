# tinygba

Tools and helpers for developing GBA programs using TinyGo.

Still highly experimental and subject to sudden changes.

## How to use

```go
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
)

func main() {
	display.Configure()

	tinygba.FillScreen(black)

	for {
		tinygba.WaitForVBlank()

		update()
	}
}

func update() {
	key := tinygba.ReadButtons()

	switch {
	case tinygba.ButtonStart.IsPushed(key):
		tinygba.FillScreen(black)

	case tinygba.ButtonSelect.IsPushed(key):
		tinygba.FillScreen(white)

	case tinygba.ButtonRight.IsPushed(key):
		tinygba.FillScreen(green)

	case tinygba.ButtonLeft.IsPushed(key):
		tinygba.FillScreen(red)
	}
}
```

## Roadmap

- pallettes
- sprites
- ?
