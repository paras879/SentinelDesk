package liveview

import (
	"bytes"
	"crypto/md5"
	"image"
	"image/jpeg"
	"log"
	"sync"
	"sync/atomic"

	"github.com/kbinani/screenshot"
	"golang.org/x/image/draw"
)

var jpegBufferPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

type Capturer struct {
	mu           sync.Mutex
	lastHash     [16]byte
	width        int
	height       int
	resizedBuf   *image.RGBA
	captureCount atomic.Int64
}

func NewCapturer(width, height int) *Capturer {
	return &Capturer{
		width:      width,
		height:     height,
		resizedBuf: image.NewRGBA(image.Rect(0, 0, width, height)),
	}
}

func (c *Capturer) Resize(width, height int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.width = width
	c.height = height
	c.resizedBuf = image.NewRGBA(image.Rect(0, 0, width, height))
}

func (c *Capturer) Capture(quality int) ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return nil, err
	}

	excludeDashboardWindow(img, bounds)

	resized := c.resizedBuf
	if resized.Bounds().Dx() != c.width || resized.Bounds().Dy() != c.height {
		resized = image.NewRGBA(image.Rect(0, 0, c.width, c.height))
		c.resizedBuf = resized
	}
	draw.BiLinear.Scale(resized, resized.Bounds(), img, img.Bounds(), draw.Over, nil)

	if !c.hasChanged(resized) {
		return nil, nil
	}

	buf := jpegBufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	if err := jpeg.Encode(buf, resized, &jpeg.Options{Quality: quality}); err != nil {
		jpegBufferPool.Put(buf)
		return nil, err
	}

	out := make([]byte, buf.Len())
	copy(out, buf.Bytes())
	buf.Reset()
	jpegBufferPool.Put(buf)

	c.captureCount.Add(1)
	if c.captureCount.Load()%10 == 1 {
		log.Printf("[Capture] frame captured quality=%d size=%d count=%d", quality, len(out), c.captureCount.Load())
	}
	return out, nil
}

func (c *Capturer) hasChanged(img *image.RGBA) bool {
	sampled := downsampleHash(img, 32)
	hash := md5.Sum(sampled)
	if hash == c.lastHash {
		return false
	}
	hash = md5.Sum(img.Pix)
	if hash == c.lastHash {
		return false
	}
	c.lastHash = hash
	return true
}

func downsampleHash(img *image.RGBA, grid int) []byte {
	b := img.Bounds()
	xStep := b.Dx() / grid
	yStep := b.Dy() / grid
	if xStep < 1 {
		xStep = 1
	}
	if yStep < 1 {
		yStep = 1
	}
	out := make([]byte, 0, grid*grid*3)
	for y := b.Min.Y; y < b.Max.Y; y += yStep {
		for x := b.Min.X; x < b.Max.X; x += xStep {
			off := img.PixOffset(x, y)
			out = append(out, img.Pix[off], img.Pix[off+1], img.Pix[off+2])
		}
	}
	return out
}


