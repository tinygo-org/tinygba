package tinygba

import (
	"device/gba"
)

// WaitForVBlank waits until the VBlank flag is set.
func WaitForVBlank() {
	// TODO: sleep until the next VBlank instead of busy waiting.
	for gba.DISP.DISPSTAT.Get()&(1<<gba.DISPSTAT_VBLANK_Pos) == 0 {
	}
}
