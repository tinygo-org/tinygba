// Package tiles provides access to the Game Boy Advance's Mode 0 tile-based
// background rendering hardware.
//
// Mode 0 supports four tiled backgrounds (BG0–BG3). The CPU writes 2-byte
// tilemap entries only when cell content changes, letting hardware handle
// all pixel rendering each frame.
//
// Register layouts based on GBATEK: https://problemkaputt.de/gbatek.htm#gbavideobgmodescontrol
//
// Example usage:
//
//	tiles.Configure(tiles.Layer0, tiles.Layer1)
//	tiles.SetupLayer(tiles.Layer1, 0, 8, tiles.Colors16, tiles.Size32x32, 1)
//	tiles.SetupLayer(tiles.Layer0, 0, 9, tiles.Colors16, tiles.Size32x32, 0)
//	tiles.SetPaletteColor(0, 1, tiles.RGB(0, 31, 0)) // palette 0, index 1 = green
//	tiles.DefineTile4bpp(0, 1, solidTilePixels)
//	tiles.Fill(8, 0, 0, 30, 20, 1, 0) // fill screen block 8 with tile 1
package tinygba

import (
	"device/gba"
	"runtime/volatile"
	"unsafe"
)

// Layer represents a GBA background layer (BG0–BG3).
type Layer uint8

const (
	Layer0 Layer = 0 // highest display priority
	Layer1 Layer = 1
	Layer2 Layer = 2
	Layer3 Layer = 3
)

// ColorMode controls tile color depth.
type ColorMode uint8

const (
	Colors16  ColorMode = 0 // 4bpp: 32 bytes/tile, 16 palettes × 16 colors
	Colors256 ColorMode = 1 // 8bpp: 64 bytes/tile, 1 palette × 256 colors
)

// MapSize controls tilemap dimensions.
type MapSize uint8

const (
	Size32x32 MapSize = 0 // 256×256 px (default; covers GBA screen)
	Size64x32 MapSize = 1
	Size32x64 MapSize = 2
	Size64x64 MapSize = 3
)

// Memory base addresses for palette RAM and VRAM.
const (
	memPAL  uintptr = 0x05000000
	memVRAM uintptr = 0x06000000
)

// mem16 returns a pointer to a volatile 16-bit register at the given address.
// All VRAM and palette RAM writes go through this to prevent optimizer caching.
func mem16(addr uintptr) *volatile.Register16 {
	return (*volatile.Register16)(unsafe.Pointer(addr))
}

// ConfigureLayers sets DISPCNT to Mode 0 and enables the given layers.
// Must be called instead of (not alongside) machine.Display.ConfigureLayers().
func ConfigureLayers(layers ...Layer) {
	// Mode 0: bits 0-2 of DISPCNT are 000 (no need to set them)
	var bits uint16
	for _, l := range layers {
		bits |= 1 << (8 + uint(l))
	}
	gba.DISP.DISPCNT.Set(bits)
}

// SetupLayer configures a background layer by writing its BGxCNT register.
//
//   - charBlock:   0–3  — which 16KB VRAM block holds this layer's tile graphics
//   - screenBlock: 0–31 — which 2KB VRAM block holds this layer's tilemap
//   - mode:        Colors16 or Colors256
//   - size:        map dimensions
//   - priority:    0 (front) to 3 (back)
func SetupLayer(layer Layer, charBlock, screenBlock uint8, mode ColorMode, size MapSize, priority uint8) {
	value := uint16(priority) |
		(uint16(charBlock) << gba.BGCNT_CHAR_BASE_Pos) |
		(uint16(mode) << gba.BGCNT_COLORS_Pos) |
		(uint16(screenBlock) << gba.BGCNT_BASE_Pos) |
		(uint16(size) << gba.BGCNT_SIZE_Pos)
	switch layer {
	case Layer0:
		gba.BGCNT0.CNT.Set(value)
	case Layer1:
		gba.BGCNT1.CNT.Set(value)
	case Layer2:
		gba.BGCNT2.CNT.Set(value)
	case Layer3:
		gba.BGCNT3.CNT.Set(value)
	}
}

// SetScroll sets the scroll offset for a layer (writes BGxHOFS / BGxVOFS).
// Values are masked to 9 bits per the GBA spec.
func SetScroll(layer Layer, x, y uint16) {
	x &= 0x1FF
	y &= 0x1FF
	switch layer {
	case Layer0:
		gba.BG0.HOFS.Set(x)
		gba.BG0.VOFS.Set(y)
	case Layer1:
		gba.BG1.HOFS.Set(x)
		gba.BG1.VOFS.Set(y)
	case Layer2:
		gba.BG2.HOFS.Set(x)
		gba.BG2.VOFS.Set(y)
	case Layer3:
		gba.BG3.HOFS.Set(x)
		gba.BG3.VOFS.Set(y)
	}
}

// RGB constructs a GBA 15-bit color from 0–31 components (R, G, B).
// Format: bits 0-4 = R, bits 5-9 = G, bits 10-14 = B.
func RGB(r, g, b uint8) uint16 {
	return uint16(r) | uint16(g)<<5 | uint16(b)<<10
}

// SetPaletteColor writes a 15-bit BGR color into background palette RAM.
// In Colors16 mode: palette 0–15, index 0–15. index 0 = transparent.
// In Colors256 mode: palette 0, index 0–255. index 0 = transparent.
func SetPaletteColor(palette, index uint8, color uint16) {
	addr := memPAL + uintptr(palette)*32 + uintptr(index)*2
	mem16(addr).Set(color)
}

// DefineTile4bpp writes a 32-byte 4bpp tile into VRAM.
//
//   - charBlock:  0–3
//   - tileIndex:  0–511 (16KB / 32 bytes = 512 tiles per block)
//   - pixels:     32 bytes — each byte encodes 2 pixels; lo nibble = left pixel,
//     hi nibble = right pixel. Pixel value 0 = transparent.
func DefineTile4bpp(charBlock uint8, tileIndex uint16, pixels [32]byte) {
	base := memVRAM + uintptr(charBlock)*16384 + uintptr(tileIndex)*32
	for i := 0; i < 16; i++ {
		val := uint16(pixels[i*2]) | uint16(pixels[i*2+1])<<8
		mem16(base + uintptr(i)*2).Set(val)
	}
}

// DefineTile8bpp writes a 64-byte 8bpp tile into VRAM.
//
//   - charBlock:  0–3
//   - tileIndex:  0–255 (16KB / 64 bytes = 256 tiles per block)
//   - pixels:     64 bytes — each byte is one palette index. 0 = transparent.
func DefineTile8bpp(charBlock uint8, tileIndex uint16, pixels [64]byte) {
	base := memVRAM + uintptr(charBlock)*16384 + uintptr(tileIndex)*64
	for i := 0; i < 32; i++ {
		val := uint16(pixels[i*2]) | uint16(pixels[i*2+1])<<8
		mem16(base + uintptr(i)*2).Set(val)
	}
}

// SetTile writes one 2-byte tilemap entry at tile coordinates (x, y).
//
//   - screenBlock: 0–31
//   - x, y:        0–31 (tile coordinates within a 32×32 map)
//   - tileIndex:   0–1023
//   - palette:     0–15 (Colors16 only; ignored for Colors256)
//   - flipH, flipV: mirror the tile horizontally/vertically
func SetTile(screenBlock uint8, x, y uint8, tileIndex uint16, palette uint8, flipH, flipV bool) {
	addr := memVRAM + uintptr(screenBlock)*2048 + (uintptr(y)*32+uintptr(x))*2
	entry := tileIndex
	if flipH {
		entry |= 1 << 10
	}
	if flipV {
		entry |= 1 << 11
	}
	entry |= uint16(palette) << 12
	mem16(addr).Set(entry)
}

// FillTiled sets a w×h region of a screen block to the same tile entry.
// Useful for clearing areas or painting solid-color backgrounds.
func FillTiled(screenBlock uint8, x, y, w, h uint8, tileIndex uint16, palette uint8) {
	for row := uint8(0); row < h; row++ {
		for col := uint8(0); col < w; col++ {
			SetTile(screenBlock, x+col, y+row, tileIndex, palette, false, false)
		}
	}
}
