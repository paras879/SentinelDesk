package liveview

import (
	"bytes"
	"crypto/md5"
	"image"
	"image/jpeg"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/kbinani/screenshot"
	"golang.org/x/image/draw"
)

type Capturer struct {
	mu       sync.Mutex
	lastHash [16]byte
	width    int
	height   int
}

func NewCapturer(width, height int) *Capturer {
	return &Capturer{
		width:  width,
		height: height,
	}
}

func (c *Capturer) Capture() ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}

	excludeDashboardWindow(img, bounds)

	resized := image.NewRGBA(image.Rect(0, 0, c.width, c.height))
	draw.BiLinear.Scale(resized, resized.Bounds(), img, img.Bounds(), draw.Over, nil)

	hash := md5.Sum(resized.Pix)
	if hash == c.lastHash {
		return nil, nil
	}
	c.lastHash = hash

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, resized, &jpeg.Options{Quality: 70}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procEnumWindows      = user32.NewProc("EnumWindows")
	procGetWindowTextW   = user32.NewProc("GetWindowTextW")
	procGetWindowRect    = user32.NewProc("GetWindowRect")
	procIsWindowVisible  = user32.NewProc("IsWindowVisible")
	procShowWindow       = user32.NewProc("ShowWindow")
)

const (
	SW_MINIMIZE = 6
	SW_RESTORE  = 9
)

var (
	dashboardHWND      uintptr
	dashboardMinimized bool
)

type rect struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}

func (c *Capturer) HideDashboard() {
	hwnd := getDashboardHWND()
	if hwnd == 0 {
		return
	}
	procShowWindow.Call(hwnd, SW_MINIMIZE)
	dashboardMinimized = true
}

func (c *Capturer) ShowDashboard() {
	hwnd := getDashboardHWND()
	if hwnd == 0 {
		return
	}
	procShowWindow.Call(hwnd, SW_RESTORE)
	dashboardMinimized = false
}

func excludeDashboardWindow(img *image.RGBA, screenBounds image.Rectangle) {
	if dashboardMinimized {
		return
	}

	r := findDashboardRect()
	if r == nil {
		return
	}

	clip := image.Rect(
		int(r.left), int(r.top),
		int(r.right), int(r.bottom),
	).Intersect(screenBounds)

	if clip.Empty() {
		return
	}

	for y := clip.Min.Y; y < clip.Max.Y; y++ {
		offset := img.PixOffset(clip.Min.X, y)
		for x := clip.Min.X; x < clip.Max.X; x++ {
			img.Pix[offset+0] = 0
			img.Pix[offset+1] = 0
			img.Pix[offset+2] = 0
			img.Pix[offset+3] = 255
			offset += 4
		}
	}
}

func getDashboardHWND() uintptr {
	if dashboardHWND == 0 {
		dashboardHWND = findWindow("SentinelDesk")
	}
	return dashboardHWND
}

func findDashboardRect() *rect {
	hwnd := getDashboardHWND()
	if hwnd == 0 {
		return nil
	}

	visible, _, _ := procIsWindowVisible.Call(hwnd)
	if visible == 0 {
		return nil
	}

	var r rect
	ret, _, _ := procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&r)))
	if ret == 0 {
		return nil
	}

	return &r
}

func findWindow(titleSubstr string) uintptr {
	var hwnd uintptr
	cb := syscall.NewCallback(func(h uintptr, lparam uintptr) uintptr {
		var buf [512]uint16
		ret, _, _ := procGetWindowTextW.Call(h, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
		if ret == 0 {
			return 1
		}
		title := strings.ToLower(syscall.UTF16ToString(buf[:]))
		if strings.Contains(title, strings.ToLower(titleSubstr)) {
			hwnd = h
			return 0
		}
		return 1
	})
	procEnumWindows.Call(cb, 0)
	return hwnd
}
