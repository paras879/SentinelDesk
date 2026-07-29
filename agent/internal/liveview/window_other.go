//go:build !windows

package liveview

import (
	"image"
)

func (c *Capturer) HideDashboard() {
	// Not implemented for non-Windows platforms
}

func (c *Capturer) ShowDashboard() {
	// Not implemented for non-Windows platforms
}

func excludeDashboardWindow(img *image.RGBA, screenBounds image.Rectangle) {
	// Not implemented for non-Windows platforms
}
