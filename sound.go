// Package sound provides access to the Game Boy Advance's sound hardware.
//
// Uses Sound Channel 2 (square wave with envelope) for tone generation.
// Register layouts based on GBATEK: https://problemkaputt.de/gbatek.htm#gbasoundcontroller
//
// Example usage:
//
//	sound.Enable(7, 7)
//	sound.PlayNote(1750, 2, 15, 30) // A4 at 440Hz, 50% duty, max volume
package tinygba

import "device/gba"

// EnableSound enables master sound and configures channel 2 output on both speakers.
// Volume should be between 0 and 7.
func EnableSound(leftVolume, rightVolume uint16) {
	// SOUNDCNT_X bit 7: master sound enable
	gba.SOUND.CNT_X.Set(1 << 7)

	// SOUNDCNT_L:
	//   Bits 0-2:   right master volume
	//   Bits 4-6:   left master volume
	//   Bits 8-11:  channels 1-4 enable right
	//   Bits 12-15: channels 1-4 enable left
	// Enable Sound 2 on both speakers (bit 9 = ch2 right, bit 13 = ch2 left)
	gba.SOUND.CNT_L.Set(rightVolume | (leftVolume << 4) | (1 << 9) | (1 << 13))

	// SOUNDCNT_H bits 0-1: PSG channel volume ratio (2 = 100%)
	gba.SOUND.CNT_H.Set(2)
}

// DisableSound disables master sound.
func DisableSound() {
	gba.SOUND.CNT_X.Set(0)
}

// PlayNote plays a tone on Sound Channel 2.
//
// frequency: register value = 2048 - (131072 / hz). E.g. A4 (440Hz) = 1750.
// duty: wave pattern (0=12.5%, 1=25%, 2=50%, 3=75%).
// volume: initial envelope volume (0-15).
// duration: sound length in (64-n)/256s units (1-63), or 0 for continuous.
func PlayNote(frequency, duty, volume uint16, duration uint8) {
	// SOUND2CNT_L (0x4000068) - Duty/Length/Envelope (all in one register):
	//   Bits 0-5:   sound length
	//   Bits 6-7:   wave duty
	//   Bits 8-10:  envelope step time (0 = hold volume)
	//   Bit  11:    envelope direction (0=decrease, 1=increase)
	//   Bits 12-15: initial volume
	gba.SOUND2.CNT_L.Set(uint16(duration) | (duty << 6) | (volume << 12))

	// SOUND2CNT_H (0x400006C) - Frequency/Control:
	//   Bits 0-10: frequency
	//   Bit  14:   length flag (1=stop when length expires)
	//   Bit  15:   initial/restart sound
	if duration > 0 {
		gba.SOUND2.CNT_H.Set(frequency | (1 << 14) | (1 << 15))
	} else {
		gba.SOUND2.CNT_H.Set(frequency | (1 << 15))
	}
}
