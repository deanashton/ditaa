package graphical

import (
	"image"
	"image/color"
	"math"

	"code.google.com/p/jamslam-freetype-go/freetype/raster"
)

const (
	STROKE_WIDTH float64 = 1
	MAGIC_K      float64 = 0.5522847498
)

type Color struct {
	R uint8 `xml:"r,attr"`
	G uint8 `xml:"g,attr"`
	B uint8 `xml:"b,attr"`
	A uint8 `xml:"a,attr"`
}

func (c Color) RGBA() color.RGBA {
	return color.RGBA{c.R, c.G, c.B, c.A}
}

var WHITE = Color{255, 255, 255, 255}

type PointType int

const (
	POINT_NORMAL PointType = iota
	POINT_ROUND
)

type Point struct {
	X      float64   `xml:"x,attr"`
	Y      float64   `xml:"y,attr"`
	Locked bool      `xml:"locked,attr"`
	Type   PointType `xml:"type,attr"`
}

func (p1 Point) NorthOf(p2 Point) bool { return p1.Y < p2.Y }
func (p1 Point) SouthOf(p2 Point) bool { return p1.Y > p2.Y }
func (p1 Point) WestOf(p2 Point) bool  { return p1.X < p2.X }
func (p1 Point) EastOf(p2 Point) bool  { return p1.X > p2.X }

func P(p Point) raster.Point {
	//TODO: handle fractional part too, but probably not needed
	return raster.Point{
		raster.Fix32(int(p.X)) << 8,
		raster.Fix32(int(p.Y)) << 8,
	}
}

func ftofix(f float64) raster.Fix32 {
	//TODO: verify this is OK
	a := math.Trunc(f)
	b := math.Ldexp(math.Abs(f-a), 8)
	return raster.Fix32(a)<<8 + raster.Fix32(b)
}

func Stroke(img *image.RGBA, path raster.Path, color color.RGBA) {
	//TODO: support dashed lines
	g := raster.NewRasterizer(img.Rect.Max.X+1, img.Rect.Max.Y+1) //TODO: +1 or not?
	raster.Stroke(g, path, ftofix(STROKE_WIDTH), nil, nil)
	painter := raster.NewRGBAPainter(img)
	painter.SetColor(color)
	g.Rasterize(painter)
}

func Fill(img *image.RGBA, path raster.Path, color color.RGBA) {
	g := raster.NewRasterizer(img.Rect.Max.X+1, img.Rect.Max.Y+1) //TODO: +1 or not?
	g.AddPath(path)
	painter := raster.NewRGBAPainter(img)
	painter.SetColor(color)
	g.Rasterize(painter)
}

func Circle(x, y, r float64) raster.Path {
	//panic(fmt.Sprint(x, y, r))
	P := func(x, y float64) raster.Point {
		return raster.Point{ftofix(x), ftofix(y)}
	}
	p1 := P(x+r, y)
	p2 := P(x, y+r)
	p3 := P(x-r, y)
	p4 := P(x, y-r)
	kr := MAGIC_K * r
	path := raster.Path{}
	// see: http://hansmuller-flex.blogspot.com/2011/04/approximating-circular-arc-with-cubic.html
	//  or: http://www.whizkidtech.redprince.net/bezier/circle/
	// etc. -- google "drawing circle with cubic curves"
	path.Start(p1)
	path.Add3(P(x+r, y+kr), P(x+kr, y+r), p2)
	path.Add3(P(x-kr, y+r), P(x-r, y+kr), p3)
	path.Add3(P(x-r, y-kr), P(x-kr, y-r), p4)
	path.Add3(P(x+kr, y-r), P(x+r, y-kr), p1)
	return path
}
