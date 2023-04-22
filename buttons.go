package tinygba

import (
	"device/gba"
)

// Buttons
var (
	ButtonA      = Button{gba.KEYINPUT_BUTTON_A_Pos}
	ButtonB      = Button{gba.KEYINPUT_BUTTON_B_Pos}
	ButtonSelect = Button{gba.KEYINPUT_BUTTON_SELECT_Pos}
	ButtonStart  = Button{gba.KEYINPUT_BUTTON_START_Pos}
	ButtonRight  = Button{gba.KEYINPUT_BUTTON_RIGHT_Pos}
	ButtonLeft   = Button{gba.KEYINPUT_BUTTON_LEFT_Pos}
	ButtonUp     = Button{gba.KEYINPUT_BUTTON_UP_Pos}
	ButtonDown   = Button{gba.KEYINPUT_BUTTON_DOWN_Pos}
	ButtonR      = Button{gba.KEYINPUT_BUTTON_R_Pos}
	ButtonL      = Button{gba.KEYINPUT_BUTTON_L_Pos}
)

const (
	ButtonPushed   uint16 = 0
	ButtonReleased uint16 = 1
)

type Button struct {
	pos int16
}

// ReadButtons gets the current state of the buttons.
// It is best to only read this one time per video frame.
func ReadButtons() uint16 {
	return gba.KEY.INPUT.Get()
}

// IsPushed checks to see if the button is pushed.
func (b Button) IsPushed(key uint16) bool {
	return key&(1<<b.pos) == ButtonPushed
}

// IsReleased checks to see if the button is released.
func (b Button) IsReleased(key uint16) bool {
	return key&(1<<b.pos) == ButtonReleased
}
