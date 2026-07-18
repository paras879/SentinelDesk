package remote

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32Clip           = windows.NewLazySystemDLL("user32.dll")
	procOpenClipboard    = user32Clip.NewProc("OpenClipboard")
	procEmptyClipboard   = user32Clip.NewProc("EmptyClipboard")
	procSetClipboardData = user32Clip.NewProc("SetClipboardData")
	procCloseClipboard   = user32Clip.NewProc("CloseClipboard")
	kernel32Clip         = windows.NewLazySystemDLL("kernel32.dll")
	procGlobalAlloc      = kernel32Clip.NewProc("GlobalAlloc")
	procRtlMoveMemory    = kernel32Clip.NewProc("RtlMoveMemory")
)

const (
	CF_UNICODETEXT = 13
	GMEM_ZEROINIT  = 0x0040
)

func setClipboard(text string) error {
	procOpenClipboard.Call(0)
	procEmptyClipboard.Call()

	utf16, _ := windows.UTF16FromString(text)
	size := len(utf16) * 2

	hMem, _, _ := procGlobalAlloc.Call(GMEM_ZEROINIT, uintptr(size))
	if hMem == 0 {
		procCloseClipboard.Call()
		return nil
	}

	procRtlMoveMemory.Call(hMem, uintptr(unsafe.Pointer(&utf16[0])), uintptr(size))

	procSetClipboardData.Call(CF_UNICODETEXT, hMem)
	procCloseClipboard.Call()
	return nil
}
