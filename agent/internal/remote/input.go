package remote

import (
	"encoding/json"
	"log"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procSendInput    = user32.NewProc("SendInput")
	procGetDC        = user32.NewProc("GetDC")
	procReleaseDC    = user32.NewProc("ReleaseDC")
	gdi32            = syscall.NewLazyDLL("gdi32.dll")
	procGetDeviceCaps = gdi32.NewProc("GetDeviceCaps")
)

const (
	INPUT_MOUSE    = 0
	INPUT_KEYBOARD = 1

	MOUSEEVENTF_MOVE        = 0x0001
	MOUSEEVENTF_ABSOLUTE    = 0x8000
	MOUSEEVENTF_LEFTDOWN    = 0x0002
	MOUSEEVENTF_LEFTUP      = 0x0004
	MOUSEEVENTF_RIGHTDOWN   = 0x0008
	MOUSEEVENTF_RIGHTUP     = 0x0010
	MOUSEEVENTF_MIDDLEDOWN  = 0x0020
	MOUSEEVENTF_MIDDLEUP    = 0x0040
	MOUSEEVENTF_WHEEL       = 0x0800
	MOUSEEVENTF_VIRTUALDESK = 0x4000

	KEYEVENTF_KEYDOWN = 0x0000
	KEYEVENTF_KEYUP   = 0x0002

	VK_BACK     = 0x08
	VK_TAB      = 0x09
	VK_RETURN   = 0x0D
	VK_SHIFT    = 0x10
	VK_CONTROL  = 0x11
	VK_MENU     = 0x12
	VK_PAUSE    = 0x13
	VK_CAPITAL  = 0x14
	VK_ESCAPE   = 0x1B
	VK_SPACE    = 0x20
	VK_PRIOR    = 0x21
	VK_NEXT     = 0x22
	VK_END      = 0x23
	VK_HOME     = 0x24
	VK_LEFT     = 0x25
	VK_UP       = 0x26
	VK_RIGHT    = 0x27
	VK_DOWN     = 0x28
	VK_DELETE   = 0x2E
	VK_LSHIFT   = 0xA0
	VK_RSHIFT   = 0xA1
	VK_LCONTROL = 0xA2
	VK_RCONTROL = 0xA3
	VK_LMENU    = 0xA4
	VK_RMENU    = 0xA5
	VK_LWIN     = 0x5B
	VK_RWIN     = 0x5C
	VK_F1       = 0x70
	VK_F2       = 0x71
	VK_F3       = 0x72
	VK_F4       = 0x73
	VK_F5       = 0x74
	VK_F6       = 0x75
	VK_F7       = 0x76
	VK_F8       = 0x77
	VK_F9       = 0x78
	VK_F10      = 0x79
	VK_F11      = 0x7A
	VK_F12      = 0x7B
)

type mouseInput struct {
	Dx          int32
	Dy          int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	_           [4]byte
	DwExtraInfo uintptr
}

type keyboardInput struct {
	WVk         uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

type inputUnion struct {
	Mi mouseInput
}

type inputStruct struct {
	Type uint32
	_    [4]byte
	Union inputUnion
}

type controlPacket struct {
	Type     string            `json:"type"`
	Button   string            `json:"button,omitempty"`
	X        float64           `json:"x,omitempty"`
	Y        float64           `json:"y,omitempty"`
	Delta    int32             `json:"delta,omitempty"`
	Key      string            `json:"key,omitempty"`
	Code     string            `json:"code,omitempty"`
	Text     string            `json:"text,omitempty"`
	Modifiers []string         `json:"modifiers,omitempty"`
}

var screenWidth, screenHeight int32

func init() {
	screenWidth, screenHeight = getScreenSize()
}

func getScreenSize() (int32, int32) {
	dc, _, _ := procGetDC.Call(0)
	if dc == 0 {
		return 1920, 1080
	}
	defer procReleaseDC.Call(0, dc)
	w, _, _ := procGetDeviceCaps.Call(dc, 118)
	h, _, _ := procGetDeviceCaps.Call(dc, 117)
	return int32(w), int32(h)
}

func HandleMessage(data []byte) {
	var pkt controlPacket
	if err := json.Unmarshal(data, &pkt); err != nil {
		log.Println("[Remote] Invalid packet:", err)
		return
	}

	if pkt.Type == "" {
		log.Println("[Remote] Packet missing type field")
		return
	}

	switch pkt.Type {
	case "mouse_move":
		handleMouseMove(pkt)
	case "mouse_down":
		handleMouseDown(pkt)
	case "mouse_up":
		handleMouseUp(pkt)
	case "double_click":
		handleDoubleClick(pkt)
	case "mouse_wheel":
		handleMouseWheel(pkt)
	case "key_down":
		handleKeyDown(pkt)
	case "key_up":
		handleKeyUp(pkt)
	case "clipboard":
		handleClipboard(pkt)
	default:
		log.Printf("[Remote] Unknown packet type: %s", pkt.Type)
	}
}

func mouseButtonToFlags(button string, down bool) uint32 {
	switch button {
	case "left":
		if down {
			return MOUSEEVENTF_LEFTDOWN
		}
		return MOUSEEVENTF_LEFTUP
	case "right":
		if down {
			return MOUSEEVENTF_RIGHTDOWN
		}
		return MOUSEEVENTF_RIGHTUP
	case "middle":
		if down {
			return MOUSEEVENTF_MIDDLEDOWN
		}
		return MOUSEEVENTF_MIDDLEUP
	}
	return 0
}

func mouseEvent(flags uint32, dx, dy int32, data uint32) {
	var mi mouseInput
	mi.DwFlags = flags | MOUSEEVENTF_ABSOLUTE | MOUSEEVENTF_VIRTUALDESK
	if flags&MOUSEEVENTF_WHEEL != 0 {
		mi.MouseData = data
	}
	if flags&(MOUSEEVENTF_MOVE|MOUSEEVENTF_ABSOLUTE) == MOUSEEVENTF_MOVE|MOUSEEVENTF_ABSOLUTE {
		if screenWidth > 0 {
			mi.Dx = int32(int64(dx) * 65535 / int64(screenWidth-1))
		}
		if screenHeight > 0 {
			mi.Dy = int32(int64(dy) * 65535 / int64(screenHeight-1))
		}
	}

	var input inputStruct
	input.Type = INPUT_MOUSE
	input.Union.Mi = mi

	sendOneInput(input)
}

func keyboardEvent(vk uint16, flags uint32) {
	var ki keyboardInput
	ki.WVk = vk
	ki.DwFlags = flags

	var raw [40]byte
	raw[0] = INPUT_KEYBOARD
	raw[4] = 0
	raw[5] = 0
	raw[6] = 0
	raw[7] = 0
	*(*uint16)(unsafe.Pointer(&raw[8])) = ki.WVk
	*(*uint16)(unsafe.Pointer(&raw[10])) = ki.WScan
	*(*uint32)(unsafe.Pointer(&raw[12])) = ki.DwFlags
	*(*uint32)(unsafe.Pointer(&raw[16])) = ki.Time
	*(*uintptr)(unsafe.Pointer(&raw[24])) = ki.DwExtraInfo

	procSendInput.Call(1, uintptr(unsafe.Pointer(&raw[0])), 40)
}

func sendOneInput(input inputStruct) {
	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}

var jsCodeToVK = map[string]uint16{
	"Backspace":   VK_BACK,
	"Tab":          VK_TAB,
	"Enter":        VK_RETURN,
	"Escape":       VK_ESCAPE,
	"Space":        VK_SPACE,
	"Delete":       VK_DELETE,
	"Home":         VK_HOME,
	"End":          VK_END,
	"PageUp":       VK_PRIOR,
	"PageDown":     VK_NEXT,
	"ArrowUp":      VK_UP,
	"ArrowDown":    VK_DOWN,
	"ArrowLeft":    VK_LEFT,
	"ArrowRight":   VK_RIGHT,
	"ShiftLeft":    VK_LSHIFT,
	"ShiftRight":   VK_RSHIFT,
	"ControlLeft":  VK_LCONTROL,
	"ControlRight": VK_RCONTROL,
	"AltLeft":      VK_LMENU,
	"AltRight":     VK_RMENU,
	"MetaLeft":     VK_LWIN,
	"MetaRight":    VK_RWIN,
	"F1":           VK_F1,
	"F2":           VK_F2,
	"F3":           VK_F3,
	"F4":           VK_F4,
	"F5":           VK_F5,
	"F6":           VK_F6,
	"F7":           VK_F7,
	"F8":           VK_F8,
	"F9":           VK_F9,
	"F10":          VK_F10,
	"F11":          VK_F11,
	"F12":          VK_F12,
}

func codeToVK(code string) uint16 {
	if vk, ok := jsCodeToVK[code]; ok {
		return vk
	}
	if len(code) == 4 && code[:3] == "Key" && code[3] >= 'A' && code[3] <= 'Z' {
		return uint16(code[3])
	}
	if len(code) == 6 && code[:5] == "Digit" && code[5] >= '0' && code[5] <= '9' {
		return 0x30 + uint16(code[5]-'0')
	}
	return 0
}

func handleMouseMove(pkt controlPacket) {
	mouseEvent(MOUSEEVENTF_MOVE, int32(pkt.X), int32(pkt.Y), 0)
}

func handleMouseDown(pkt controlPacket) {
	flags := mouseButtonToFlags(pkt.Button, true)
	if flags != 0 {
		mouseEvent(flags, int32(pkt.X), int32(pkt.Y), 0)
	}
}

func handleMouseUp(pkt controlPacket) {
	flags := mouseButtonToFlags(pkt.Button, false)
	if flags != 0 {
		mouseEvent(flags, 0, 0, 0)
	}
}

func handleDoubleClick(pkt controlPacket) {
	flags := mouseButtonToFlags(pkt.Button, true)
	if flags == 0 {
		return
	}
	x := int32(pkt.X)
	y := int32(pkt.Y)
	mouseEvent(flags, x, y, 0)
	time.Sleep(30 * time.Millisecond)
	mouseEvent(flags ^ 1, 0, 0, 0) // up = flip lowest bit
	time.Sleep(30 * time.Millisecond)
	mouseEvent(flags, x, y, 0)
	time.Sleep(30 * time.Millisecond)
	mouseEvent(flags ^ 1, 0, 0, 0)
}

func handleMouseWheel(pkt controlPacket) {
	mouseEvent(MOUSEEVENTF_WHEEL, 0, 0, uint32(pkt.Delta))
}

func handleKeyDown(pkt controlPacket) {
	vk := codeToVK(pkt.Code)
	if vk == 0 {
		return
	}
	for _, m := range pkt.Modifiers {
		mkVk := codeToVK(m)
		if mkVk != 0 {
			keyboardEvent(mkVk, KEYEVENTF_KEYDOWN)
		}
	}
	keyboardEvent(vk, KEYEVENTF_KEYDOWN)
}

func handleKeyUp(pkt controlPacket) {
	vk := codeToVK(pkt.Code)
	if vk == 0 {
		return
	}
	keyboardEvent(vk, KEYEVENTF_KEYUP)
	for i := len(pkt.Modifiers) - 1; i >= 0; i-- {
		mkVk := codeToVK(pkt.Modifiers[i])
		if mkVk != 0 {
			keyboardEvent(mkVk, KEYEVENTF_KEYUP)
		}
	}
}

func handleClipboard(pkt controlPacket) {
	if pkt.Text == "" {
		return
	}
	if err := setClipboard(pkt.Text); err != nil {
		log.Println("Remote clipboard error:", err)
	}
}
