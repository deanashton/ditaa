package graphical

import (
	"image"
	"image/color"
)

// from: https://github.com/bolknote/go-gd/blob/master/gd.go
// from: http://incubator.quasimondo.com/processing/fast_blur_deluxe.php
//
// see also: idea of "horizontally and vertically applied motion blur, reapplied multiple times":
// http://www.gamasutra.com/view/feature/131511/four_tricks_for_fast_blurring_in_.php?page=2
// http://web.archive.org/web/20060718054020/http://www.acm.uiuc.edu/siggraph/workshops/wjarosz_convolution_2001.pdf (via wikipedia "box blur")
// http://incubator.quasimondo.com/processing/superfast_blur.php
// http://oppei.tumblr.com/post/47346262/super-fast-blur-v1-1-by-mario-klingemann

//FIXME: what the license???

// Stack Blur Algorithm by Mario Klingemann <mario@quasimondo.com>
// "Go" language port by Evgeny Stepanischev http://bolknote.ru
// [slightly modified for image.Image by Mateusz Czapliński]
//FIXME: make sure it works on subimages (img.Bounds().Min != (0,0))
func StackBlur(img *image.RGBA, radius int, keepalpha bool) {
	if radius < 1 {
		return
	}

	//w, h := int(img.Sx()), int(img.Sy())
	//w, h := int(img.Bounds().Max.X-img.Bounds().Min.X)
	w, h := int(img.Bounds().Max.X), int(img.Bounds().Max.Y)
	wm, hm, wh, div := w-1, h-1, w*h, radius*2+1

	len := map[bool]int{true: 3, false: 4}[keepalpha]

	rgba := make([][]byte, len)
	for i := 0; i < len; i++ {
		rgba[i] = make([]byte, wh)
	}

	vmin := make([]int, max(w, h))

	var x, y, i, yp, yi, yw, stackpointer, stackstart, rbs int
	var sir *[4]byte

	divsum := (div + 1) >> 1
	divsum *= divsum

	dv := make([]byte, 256*divsum)

	for i = 0; i < 256*divsum; i++ {
		dv[i] = byte(i / divsum)
	}

	yw, yi = 0, 0
	stack := make([][4]byte, div)
	r1 := radius + 1

	for y = 0; y < h; y++ {
		sum := make([]int, len)
		insum := make([]int, len)
		outsum := make([]int, len)

		for i = -radius; i <= radius; i++ {
			coords := yi + min(wm, max(i, 0))
			yc := coords / w
			xc := coords % w

			//p := img.ColorsForIndex(img.ColorAt(xc, yc))
			p := img.At(xc, yc)

			sir = &stack[i+radius]
			//sir[0] = (byte)(p["red"])
			//sir[1] = (byte)(p["green"])
			//sir[2] = (byte)(p["blue"])
			//sir[3] = (byte)(p["alpha"])
			{
				r, g, b, a := p.RGBA()
				sir[0] = byte(r >> 8)
				sir[1] = byte(g >> 8)
				sir[2] = byte(b >> 8)
				sir[3] = byte(a >> 8)
			}

			rbs = r1 - abs(i)
			for i := 0; i < len; i++ {
				sum[i] += int(sir[i]) * rbs
			}

			if i > 0 {
				for i := 0; i < len; i++ {
					insum[i] += int(sir[i])
				}
			} else {
				for i := 0; i < len; i++ {
					outsum[i] += int(sir[i])
				}
			}
		}

		stackpointer = radius

		for x = 0; x < w; x++ {
			for i := 0; i < len; i++ {
				rgba[i][yi] = dv[sum[i]]
				sum[i] -= outsum[i]
			}

			stackstart = stackpointer - radius + div
			sir = &stack[stackstart%div]

			for i := 0; i < len; i++ {
				outsum[i] -= int(sir[i])
			}

			if y == 0 {
				vmin[x] = min(x+radius+1, wm)
			}

			coords := yw + vmin[x]
			yc := coords / w
			xc := coords % w

			//p := img.ColorsForIndex(img.ColorAt(xc, yc))
			p := img.At(xc, yc)

			//sir[0] = (byte)(p["red"])
			//sir[1] = (byte)(p["green"])
			//sir[2] = (byte)(p["blue"])
			//sir[3] = (byte)(p["alpha"])
			{
				r, g, b, a := p.RGBA()
				sir[0] = byte(r >> 8)
				sir[1] = byte(g >> 8)
				sir[2] = byte(b >> 8)
				sir[3] = byte(a >> 8)
			}

			for i := 0; i < len; i++ {
				insum[i] += int(sir[i])
				sum[i] += insum[i]
			}

			stackpointer = (stackpointer + 1) % div
			sir = &stack[stackpointer%div]

			for i := 0; i < len; i++ {
				outsum[i] += int(sir[i])
				insum[i] -= int(sir[i])
			}

			yi++
		}

		yw += w
	}

	for x = 0; x < w; x++ {
		sum := make([]int, len)
		insum := make([]int, len)
		outsum := make([]int, len)

		yp = -radius * w

		for i = -radius; i <= radius; i++ {
			yi = max(0, yp) + x

			sir = &stack[i+radius]

			for i := 0; i < len; i++ {
				sir[i] = rgba[i][yi]
			}
			rbs = r1 - abs(i)

			for i := 0; i < len; i++ {
				sum[i] += int(rgba[i][yi]) * rbs
			}

			if i > 0 {
				for i := 0; i < len; i++ {
					insum[i] += int(sir[i])
				}
			} else {
				for i := 0; i < len; i++ {
					outsum[i] += int(sir[i])
				}
			}

			if i < hm {
				yp += w
			}
		}

		yi = x

		stackpointer = radius

		for y = 0; y < h; y++ {
			var alpha int

			if keepalpha {
				//alpha = img.ColorsForIndex(img.ColorAt(yi%w, yi/w))["alpha"]
				_, _, _, a := img.At(yi%w, yi/w).RGBA()
				alpha = int(a >> 8)
			} else {
				alpha = int(dv[sum[3]])
			}

			//newpxl := img.ColorAllocateAlpha(int(dv[sum[0]]), int(dv[sum[1]]), int(dv[sum[2]]), alpha)
			newpxl := color.RGBA{byte(dv[sum[0]]), byte(dv[sum[1]]), byte(dv[sum[2]]), byte(alpha)}
			//if newpxl == -1 {
			//	newpxl = img.ColorClosestAlpha(int(dv[sum[0]]), int(dv[sum[1]]), int(dv[sum[2]]), alpha)
			//}

			//img.SetPixel(yi%w, yi/w, newpxl)
			img.SetRGBA(yi%w, yi/w, newpxl)

			for i := 0; i < len; i++ {
				sum[i] -= outsum[i]
			}

			stackstart = stackpointer - radius + div
			sir = &stack[stackstart%div]

			for i := 0; i < len; i++ {
				outsum[i] -= int(sir[i])
			}

			if x == 0 {
				vmin[y] = min(y+r1, hm) * w
			}

			p := x + vmin[y]

			for i := 0; i < len; i++ {
				sir[i] = rgba[i][p]
				insum[i] += int(sir[i])
				sum[i] += insum[i]
			}

			stackpointer = (stackpointer + 1) % div
			sir = &stack[stackpointer]

			for i := 0; i < len; i++ {
				outsum[i] += int(sir[i])
				insum[i] -= int(sir[i])
			}

			yi += w
		}
	}
}

func min(n1, n2 int) int {
	if n1 < n2 {
		return n1
	}
	return n2
}

func max(n1, n2 int) int {
	if n1 > n2 {
		return n1
	}
	return n2
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
