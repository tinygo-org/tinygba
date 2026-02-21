package main

import (
	"tinygo.org/x/tinygba"
)

func main() {
	InitTiles()
	tinygba.EnableSound(7, 7)
	for {
		tinygba.WaitForVBlank()
		Update()
		DrawTile()
	}
}
