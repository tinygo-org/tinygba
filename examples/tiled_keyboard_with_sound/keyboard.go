package main

// Tile-based keyboard rendering using the tiles package (GBA Mode 0).
//
// Layout (one tile = 8×8 px; each key = 2×2 tiles = 16×16 px, no gap):
//
//	Row 0 (10 keys "qwertyuiop"): tile x 5–24, tile y 13–14
//	Row 1 ( 9 keys "asdfghjkl"): tile x 6–23, tile y 15–16
//	Row 2 ( 7 keys "zxcvbnm"):    tile x 8–21, tile y 17–18
//
// VRAM layout (character block 0, 4bpp):
//
//	Tile  0         — transparent (blank)
//	Tile  1         — solid gray  (normal key background)
//	Tile  2         — solid purple (selected key background)
//	Tiles 10–113   — letter glyph quarters, 4 per letter (a–z)
//	                  letter i → tiles 10+i*4 .. 10+i*4+3
//	                  (top-left, top-right, bottom-left, bottom-right)
//
// Screen blocks (each 2 KB):
//
//	Block 8 → BG1 tilemap — colored key backgrounds
//	Block 9 → BG0 tilemap — letter overlays (transparent bg)
//
// Centering: each 5×7 glyph is placed at key-relative pixel (5, 4), so it
// straddles all four 8×8 tiles of the 2×2 key block:
//
//	Key cols 5–7  → left tiles  (tile-cols 5,6,7)
//	Key cols 8–9  → right tiles (tile-cols 0,1)
//	Key rows 4–7  → top tiles   (tile-rows 4,5,6,7)
//	Key rows 8–10 → bottom tiles (tile-rows 0,1,2)

import (
	"tinygo.org/x/tinygba"
)

// Keyboard layout: three QWERTY rows.
var layout = [][]string{
	{"q", "w", "e", "r", "t", "y", "u", "i", "o", "p"},
	{"a", "s", "d", "f", "g", "h", "j", "k", "l"},
	{"z", "x", "c", "v", "b", "n", "m"},
}

// Cursor state.
var (
	selectedRow        = 0
	selectedCol        = 0
	prevButtons uint16 = 0xFFFF // all bits high = no buttons pressed (GBA active-low)
)

// VRAM / tilemap constants.
const (
	kbCharBlock    uint8  = 0  // character block holding tile pixel data
	kbBgBlock      uint8  = 8  // screen block for BG1 (key backgrounds)
	kbFgBlock      uint8  = 9  // screen block for BG0 (letter glyphs)
	tileKeyNormal  uint16 = 1  // solid gray
	tileKeySelect  uint16 = 2  // solid highlight
	tileLetterBase        = 10 // first tile index for letter quarters
)

// Palette color indices (all in palette 0).
const (
	palKeyGray = 1 // normal key background
	palKeyHL   = 2 // selected key background
	palLetter  = 3 // letter glyph pixels
)

// Keyboard tile Y-rows: row r starts at tile row (13 + r*2).
const kbTileY0 = 13

// keyTileX returns the leftmost tile X for the key at (row, col).
func keyTileX(row, col int) uint8 {
	switch row {
	case 0:
		return uint8(5 + col*2) // 10 keys → tiles 5–24
	case 1:
		return uint8(6 + col*2) // 9 keys → tiles 6–23
	default:
		return uint8(8 + col*2) // 7 keys → tiles 8–21
	}
}

// keyTileY returns the top tile Y for keyboard row r.
func keyTileY(row int) uint8 { return uint8(kbTileY0 + row*2) }

// prevTileRow / prevTileCol track which key was last highlighted.
var prevTileRow, prevTileCol = -1, -1

// justPressed returns true only on the frame a button transitions from released to pressed.
func justPressed(button tinygba.Button, curr, prev uint16) bool {
	return button.IsPushed(curr) && !button.IsPushed(prev)
}

// InitTiles sets up the display (Mode 0), palette, tile pixel data, and the
// initial tilemap. Call once at startup in place of machine.Display.Configure().
func InitTiles() {
	// Enable Mode 0 with BG0 (letters, front) and BG1 (backgrounds, back).
	tinygba.ConfigureLayers(tinygba.Layer0, tinygba.Layer1)
	tinygba.SetupLayer(tinygba.Layer1, kbCharBlock, kbBgBlock, tinygba.Colors16, tinygba.Size32x32, 1)
	tinygba.SetupLayer(tinygba.Layer0, kbCharBlock, kbFgBlock, tinygba.Colors16, tinygba.Size32x32, 0)
	tinygba.SetScroll(tinygba.Layer0, 0, 0)
	tinygba.SetScroll(tinygba.Layer1, 0, 0)

	// Palette 0:
	//   index 0 — transparent / backdrop (left as black)
	//   index 1 — normal key gray
	//   index 2 — selected key purple
	//   index 3 — letter white
	tinygba.SetPaletteColor(0, palKeyGray, tinygba.RGB(12, 12, 12))
	tinygba.SetPaletteColor(0, palKeyHL, tinygba.RGB(14, 8, 20))
	tinygba.SetPaletteColor(0, palLetter, tinygba.RGB(31, 31, 31))

	// Tile 1: solid gray (all pixels = palette index 1).
	var solid [32]byte
	for i := range solid {
		solid[i] = palKeyGray | (palKeyGray << 4)
	}
	tinygba.DefineTile4bpp(kbCharBlock, tileKeyNormal, solid)

	// Tile 2: solid highlight (all pixels = palette index 2).
	for i := range solid {
		solid[i] = palKeyHL | (palKeyHL << 4)
	}
	tinygba.DefineTile4bpp(kbCharBlock, tileKeySelect, solid)

	// Tiles 10–113: four quarter-tiles per letter (a–z).
	for i, rows := range letterGlyphs {
		tl, tr, bl, br := glyphQuarterTiles(rows)
		base := uint16(tileLetterBase + i*4)
		tinygba.DefineTile4bpp(kbCharBlock, base+0, tl)
		tinygba.DefineTile4bpp(kbCharBlock, base+1, tr)
		tinygba.DefineTile4bpp(kbCharBlock, base+2, bl)
		tinygba.DefineTile4bpp(kbCharBlock, base+3, br)
	}

	// Clear both screen blocks to tile 0 (transparent).
	tinygba.FillTiled(kbBgBlock, 0, 0, 30, 20, 0, 0)
	tinygba.FillTiled(kbFgBlock, 0, 0, 30, 20, 0, 0)

	// Write background tiles for every key (all normal to start).
	for row, rowKeys := range layout {
		for col := range rowKeys {
			setKeyBg(row, col, tileKeyNormal)
		}
	}
	// Write the four letter quarter-tiles for each key into BG0.
	setLetterTiles()

	// Highlight the initial selection.
	setKeyBg(selectedRow, selectedCol, tileKeySelect)
	prevTileRow, prevTileCol = selectedRow, selectedCol
}

// DrawTile updates only the tiles that changed due to a new key selection.
// Call every frame after UpdateTiledMode has run.
func DrawTile() {
	if selectedRow == prevTileRow && selectedCol == prevTileCol {
		return
	}
	setKeyBg(prevTileRow, prevTileCol, tileKeyNormal)
	setKeyBg(selectedRow, selectedCol, tileKeySelect)
	prevTileRow, prevTileCol = selectedRow, selectedCol
}

// Update handles input and navigates the keyboard in tile mode.
// Call every frame before DrawTile.
func Update() {
	curr := tinygba.ReadButtons()

	var playBeep bool

	switch {
	case justPressed(tinygba.ButtonRight, curr, prevButtons):
		selectedCol++
		if selectedCol >= len(layout[selectedRow]) {
			selectedCol = 0
		}
		playBeep = true
	case justPressed(tinygba.ButtonLeft, curr, prevButtons):
		selectedCol--
		if selectedCol < 0 {
			selectedCol = len(layout[selectedRow]) - 1
		}
		playBeep = true
	case justPressed(tinygba.ButtonDown, curr, prevButtons):
		selectedRow++
		if selectedRow >= len(layout) {
			selectedRow = 0
		}
		if selectedCol >= len(layout[selectedRow]) {
			selectedCol = len(layout[selectedRow]) - 1
		}
		playBeep = true
	case justPressed(tinygba.ButtonUp, curr, prevButtons):
		selectedRow--
		if selectedRow < 0 {
			selectedRow = len(layout) - 1
		}
		if selectedCol >= len(layout[selectedRow]) {
			selectedCol = len(layout[selectedRow]) - 1
		}
		playBeep = true
	}

	if playBeep {
		tinygba.PlayNote(1800, 0, 6, 57)
	}

	prevButtons = curr
}

// setKeyBg writes the 2×2 BG1 background tiles for one key.
func setKeyBg(row, col int, tile uint16) {
	x := keyTileX(row, col)
	y := keyTileY(row)
	tinygba.SetTile(kbBgBlock, x, y, tile, 0, false, false)
	tinygba.SetTile(kbBgBlock, x+1, y, tile, 0, false, false)
	tinygba.SetTile(kbBgBlock, x, y+1, tile, 0, false, false)
	tinygba.SetTile(kbBgBlock, x+1, y+1, tile, 0, false, false)
}

// setLetterTiles writes the four quarter-tiles for each key's letter into BG0.
func setLetterTiles() {
	for row, rowKeys := range layout {
		for col, key := range rowKeys {
			letterIdx := int(key[0] - 'a') // 'a'→0, 'z'→25
			x := keyTileX(row, col)
			y := keyTileY(row)
			base := uint16(tileLetterBase + letterIdx*4)
			tinygba.SetTile(kbFgBlock, x, y, base+0, 0, false, false)
			tinygba.SetTile(kbFgBlock, x+1, y, base+1, 0, false, false)
			tinygba.SetTile(kbFgBlock, x, y+1, base+2, 0, false, false)
			tinygba.SetTile(kbFgBlock, x+1, y+1, base+3, 0, false, false)
		}
	}
}

// glyphQuarterTiles splits a centered 5×7 glyph across the four 8×8 tiles of
// a 2×2 key block, returning (top-left, top-right, bottom-left, bottom-right).
func glyphQuarterTiles(rows [7]uint8) (tl, tr, bl, br [32]byte) {
	for gr := 0; gr < 7; gr++ {
		bits := rows[gr]
		var g [5]uint8
		for i := range g {
			if (bits>>(4-uint(i)))&1 != 0 {
				g[i] = palLetter
			}
		}
		var tileRow int
		if gr < 4 {
			tileRow = 4 + gr
		} else {
			tileRow = gr - 4
		}
		base := tileRow * 4
		if gr < 4 {
			tl[base+2] = 0 | (g[0] << 4)
			tl[base+3] = g[1] | (g[2] << 4)
			tr[base+0] = g[3] | (g[4] << 4)
		} else {
			bl[base+2] = 0 | (g[0] << 4)
			bl[base+3] = g[1] | (g[2] << 4)
			br[base+0] = g[3] | (g[4] << 4)
		}
	}
	return
}

// letterGlyphs holds 5×7 pixel bitmaps for 'a'–'z' (indices 0–25).
// Each [7]uint8 is seven rows; each uint8's lower 5 bits are the row pixels
// (bit 4 = leftmost column, bit 0 = rightmost column).
var letterGlyphs = [26][7]uint8{
	{14, 17, 17, 31, 17, 17, 17}, // a  .###. #...# #...# ##### #...# #...# #...#
	{30, 17, 17, 30, 17, 17, 30}, // b  ####. #...# #...# ####. #...# #...# ####.
	{14, 17, 16, 16, 16, 17, 14}, // c  .###. #...# #.... #.... #.... #...# .###.
	{30, 17, 17, 17, 17, 17, 30}, // d  ####. #...# #...# #...# #...# #...# ####.
	{31, 16, 16, 30, 16, 16, 31}, // e  ##### #.... #.... ####. #.... #.... #####
	{31, 16, 16, 30, 16, 16, 16}, // f  ##### #.... #.... ####. #.... #.... #....
	{14, 17, 16, 23, 17, 17, 14}, // g  .###. #...# #.... #.### #...# #...# .###.
	{17, 17, 17, 31, 17, 17, 17}, // h  #...# #...# #...# ##### #...# #...# #...#
	{31, 4, 4, 4, 4, 4, 31},      // i  ##### ..#.. ..#.. ..#.. ..#.. ..#.. #####
	{15, 1, 1, 1, 17, 17, 14},    // j  .#### ....# ....# ....# #...# #...# .###.
	{17, 18, 20, 24, 20, 18, 17}, // k  #...# #..#. #.#.. ##... #.#.. #..#. #...#
	{16, 16, 16, 16, 16, 16, 31}, // l  #.... #.... #.... #.... #.... #.... #####
	{17, 27, 21, 17, 17, 17, 17}, // m  #...# ##.## #.#.# #...# #...# #...# #...#
	{17, 25, 21, 19, 17, 17, 17}, // n  #...# ##..# #.#.# #..## #...# #...# #...#
	{14, 17, 17, 17, 17, 17, 14}, // o  .###. #...# #...# #...# #...# #...# .###.
	{30, 17, 17, 30, 16, 16, 16}, // p  ####. #...# #...# ####. #.... #.... #....
	{14, 17, 17, 17, 21, 18, 13}, // q  .###. #...# #...# #...# #.#.# #..#. .##.#
	{30, 17, 17, 30, 20, 18, 17}, // r  ####. #...# #...# ####. #.#.. #..#. #...#
	{15, 16, 16, 14, 1, 1, 30},   // s  .#### #.... #.... .###. ....# ....# ####.
	{31, 4, 4, 4, 4, 4, 4},       // t  ##### ..#.. ..#.. ..#.. ..#.. ..#.. ..#..
	{17, 17, 17, 17, 17, 17, 14}, // u  #...# #...# #...# #...# #...# #...# .###.
	{17, 17, 17, 17, 10, 10, 4},  // v  #...# #...# #...# #...# .#.#. .#.#. ..#..
	{17, 17, 17, 17, 21, 27, 17}, // w  #...# #...# #...# #...# #.#.# ##.## #...#
	{17, 10, 10, 4, 10, 10, 17},  // x  #...# .#.#. .#.#. ..#.. .#.#. .#.#. #...#
	{17, 17, 10, 4, 4, 4, 4},     // y  #...# #...# .#.#. ..#.. ..#.. ..#.. ..#..
	{31, 1, 2, 4, 8, 16, 31},     // z  ##### ....# ...#. ..#.. .#... #.... #####
}
