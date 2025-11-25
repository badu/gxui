package purego

import (
	"github.com/badu/gxui"
)

func translateKeyboardKey(in Key) gxui.KeyboardKey {
	switch in {
	case KeySpace:
		return gxui.KeySpace
	case KeyApostrophe:
		return gxui.KeyApostrophe
	case KeyComma:
		return gxui.KeyComma
	case KeyMinus:
		return gxui.KeyMinus
	case KeyPeriod:
		return gxui.KeyPeriod
	case KeySlash:
		return gxui.KeySlash
	case Key0:
		return gxui.Key0
	case Key1:
		return gxui.Key1
	case Key2:
		return gxui.Key2
	case Key3:
		return gxui.Key3
	case Key4:
		return gxui.Key4
	case Key5:
		return gxui.Key5
	case Key6:
		return gxui.Key6
	case Key7:
		return gxui.Key7
	case Key8:
		return gxui.Key8
	case Key9:
		return gxui.Key9
	case KeySemicolon:
		return gxui.KeySemicolon
	case KeyEqual:
		return gxui.KeyEqual
	case KeyA:
		return gxui.KeyA
	case KeyB:
		return gxui.KeyB
	case KeyC:
		return gxui.KeyC
	case KeyD:
		return gxui.KeyD
	case KeyE:
		return gxui.KeyE
	case KeyF:
		return gxui.KeyF
	case KeyG:
		return gxui.KeyG
	case KeyH:
		return gxui.KeyH
	case KeyI:
		return gxui.KeyI
	case KeyJ:
		return gxui.KeyJ
	case KeyK:
		return gxui.KeyK
	case KeyL:
		return gxui.KeyL
	case KeyM:
		return gxui.KeyM
	case KeyN:
		return gxui.KeyN
	case KeyO:
		return gxui.KeyO
	case KeyP:
		return gxui.KeyP
	case KeyQ:
		return gxui.KeyQ
	case KeyR:
		return gxui.KeyR
	case KeyS:
		return gxui.KeyS
	case KeyT:
		return gxui.KeyT
	case KeyU:
		return gxui.KeyU
	case KeyV:
		return gxui.KeyV
	case KeyW:
		return gxui.KeyW
	case KeyX:
		return gxui.KeyX
	case KeyY:
		return gxui.KeyY
	case KeyZ:
		return gxui.KeyZ
	case KeyLeftBracket:
		return gxui.KeyLeftBracket
	case KeyBackslash:
		return gxui.KeyBackslash
	case KeyRightBracket:
		return gxui.KeyRightBracket
	case KeyGraveAccent:
		return gxui.KeyGraveAccent
	case KeyWorld1:
		return gxui.KeyWorld1
	case KeyWorld2:
		return gxui.KeyWorld2
	case KeyEscape:
		return gxui.KeyEscape
	case KeyEnter:
		return gxui.KeyEnter
	case KeyTab:
		return gxui.KeyTab
	case KeyBackspace:
		return gxui.KeyBackspace
	case KeyInsert:
		return gxui.KeyInsert
	case KeyDelete:
		return gxui.KeyDelete
	case KeyRight:
		return gxui.KeyRight
	case KeyLeft:
		return gxui.KeyLeft
	case KeyDown:
		return gxui.KeyDown
	case KeyUp:
		return gxui.KeyUp
	case KeyPageUp:
		return gxui.KeyPageUp
	case KeyPageDown:
		return gxui.KeyPageDown
	case KeyHome:
		return gxui.KeyHome
	case KeyEnd:
		return gxui.KeyEnd
	case KeyCapsLock:
		return gxui.KeyCapsLock
	case KeyScrollLock:
		return gxui.KeyScrollLock
	case KeyNumLock:
		return gxui.KeyNumLock
	case KeyPrintScreen:
		return gxui.KeyPrintScreen
	case KeyPause:
		return gxui.KeyPause
	case KeyF1:
		return gxui.KeyF1
	case KeyF2:
		return gxui.KeyF2
	case KeyF3:
		return gxui.KeyF3
	case KeyF4:
		return gxui.KeyF4
	case KeyF5:
		return gxui.KeyF5
	case KeyF6:
		return gxui.KeyF6
	case KeyF7:
		return gxui.KeyF7
	case KeyF8:
		return gxui.KeyF8
	case KeyF9:
		return gxui.KeyF9
	case KeyF10:
		return gxui.KeyF10
	case KeyF11:
		return gxui.KeyF11
	case KeyF12:
		return gxui.KeyF12
	case KeyKP0:
		return gxui.KeyKp0
	case KeyKP1:
		return gxui.KeyKp1
	case KeyKP2:
		return gxui.KeyKp2
	case KeyKP3:
		return gxui.KeyKp3
	case KeyKP4:
		return gxui.KeyKp4
	case KeyKP5:
		return gxui.KeyKp5
	case KeyKP6:
		return gxui.KeyKp6
	case KeyKP7:
		return gxui.KeyKp7
	case KeyKP8:
		return gxui.KeyKp8
	case KeyKP9:
		return gxui.KeyKp9
	case KeyKPDecimal:
		return gxui.KeyKpDecimal
	case KeyKPDivide:
		return gxui.KeyKpDivide
	case KeyKPMultiply:
		return gxui.KeyKpMultiply
	case KeyKPSubtract:
		return gxui.KeyKpSubtract
	case KeyKPAdd:
		return gxui.KeyKpAdd
	case KeyKPEnter:
		return gxui.KeyKpEnter
	case KeyKPEqual:
		return gxui.KeyKpEqual
	case KeyLeftShift:
		return gxui.KeyLeftShift
	case KeyLeftControl:
		return gxui.KeyLeftControl
	case KeyLeftAlt:
		return gxui.KeyLeftAlt
	case KeyLeftSuper:
		return gxui.KeyLeftSuper
	case KeyRightShift:
		return gxui.KeyRightShift
	case KeyRightControl:
		return gxui.KeyRightControl
	case KeyRightAlt:
		return gxui.KeyRightAlt
	case KeyRightSuper:
		return gxui.KeyRightSuper
	case KeyMenu:
		return gxui.KeyMenu
	default:
		return gxui.KeyUnknown
	}
}

func translateKeyboardModifier(from ModifierKey) gxui.KeyboardModifier {
	out := gxui.ModNone
	if from&ModShift != 0 {
		out |= gxui.ModShift
	}
	if from&ModControl != 0 {
		out |= gxui.ModControl
	}
	if from&ModAlt != 0 {
		out |= gxui.ModAlt
	}
	if from&ModSuper != 0 {
		out |= gxui.ModSuper
	}
	return out
}
