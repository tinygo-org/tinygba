package tinygba

import (
	"device/gba"

	"image/color"
	"runtime/volatile"
	"unsafe"
)

func NewDisplayMode4() *DisplayMode4 {
	return &DisplayMode4{
		activeFrame: gba.DISPCNT_FRAMESELECT_FRAME0,
		frame0:      (*[160][120]volatile.Register16)(unsafe.Pointer(uintptr(gba.MEM_VRAM_FRONT))),
		frame1:      (*[160][120]volatile.Register16)(unsafe.Pointer(uintptr(gba.MEM_VRAM_BACK))),
	}
}

type DisplayMode4 struct {
	activeFrame uint8
	frame0      *[160][120]volatile.Register16
	frame1      *[160][120]volatile.Register16
}

func (d *DisplayMode4) Configure() {
	// Use video mode 4 (in BG4, 2 frames of 8bpp bitmap in VRAM)
	gba.DISP.DISPCNT.Set(gba.DISPCNT_BGMODE_4<<gba.DISPCNT_BGMODE_Pos |
		gba.DISPCNT_SCREENDISPLAY_BG2_ENABLE<<gba.DISPCNT_SCREENDISPLAY_BG2_Pos)
}

func (d *DisplayMode4) Size() (x, y int16) {
	return 240, 160
}

// SetPixel draws to the inactive frame in Mode 4.
// Mode4 we can only write word-length data with correct alignment.
// So we need to read the data for the word into which we want to write,
// mask it for the data we want to write, and only then write the data
// to the target location.
func (d *DisplayMode4) SetPixel(x, y int16, c color.RGBA) {
	offset := x >> 1
	pixel := uint16(((c.R<<3)>>8)<<5 | ((c.G<<3)>>8)<<2 | ((c.B << 2) >> 8))

	if d.activeFrame == gba.DISPCNT_FRAMESELECT_FRAME0 {
		val := d.frame1[y][offset].Get()
		if x&1 > 0 {
			// odd
			d.frame1[y][offset].Set((val & 0x00ff) | (pixel << 8))
		} else {
			// even
			d.frame1[y][offset].Set((val & 0xff00) | pixel)
		}
	} else {
		val := d.frame0[y][offset].Get()
		if x&1 > 0 {
			// odd
			d.frame0[y][offset].Set((val & 0x00ff) | (pixel << 8))
		} else {
			// even
			d.frame0[y][offset].Set((val & 0xff00) | pixel)
		}
	}
}

// Display switches the active frame in Mode 4.
func (d *DisplayMode4) Display() error {
	// switch the active frame
	if d.activeFrame == gba.DISPCNT_FRAMESELECT_FRAME0 {
		gba.DISP.DISPCNT.SetBits(gba.DISPCNT_FRAMESELECT_FRAME1 << gba.DISPCNT_FRAMESELECT_Pos)

		d.activeFrame = gba.DISPCNT_FRAMESELECT_FRAME1
	} else {
		gba.DISP.DISPCNT.ClearBits(gba.DISPCNT_FRAMESELECT_FRAME1 << gba.DISPCNT_FRAMESELECT_Pos)

		d.activeFrame = gba.DISPCNT_FRAMESELECT_FRAME0
	}

	return nil
}
