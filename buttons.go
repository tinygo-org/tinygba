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
	ButtonPushed   = 0
	ButtonReleased = 1
)

type Button struct {
	pos int16
}

func (b Button) Get() int {
	return int(gba.KEY.INPUT.Get() & (1 << b.pos))
}

func (b Button) IsPushed() bool {
	return b.Get() == ButtonPushed
}

func (b Button) IsReleased() bool {
	return b.Get() == ButtonReleased
}
