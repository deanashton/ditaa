package main

import (
	"code.google.com/p/freetype-go/freetype/raster"
	"fmt"
	"math"
)

type Grid struct {
	W     int `xml:"width"`
	H     int `xml:"height"`
	CellW int `xml:"cellWidth"`
	CellH int `xml:"cellHeight"`
}

type ShapeType int

const (
	TYPE_SIMPLE ShapeType = iota
	TYPE_ARROWHEAD
	TYPE_POINT_MARKER
	TYPE_DOCUMENT
	TYPE_STORAGE
	TYPE_IO
	TYPE_DECISION
	TYPE_MANUAL_OPERATION // upside-down trapezoid
	TYPE_TRAPEZOID        // rightside-up trapezoid
	TYPE_ELLIPSE
	TYPE_CUSTOM ShapeType = 9999
)

type Shape struct {
	Type        ShapeType `xml:"type"`
	FillColor   *Color    `xml:"fillColor"`
	StrokeColor Color     `xml:"strokeColor"`
	Closed      bool      `xml:"isClosed"`
	Dashed      bool      `xml:"isStrokeDashed"`
	Points      []Point   `xml:"points>point"`
}

func (s *Shape) MakeMarkerPaths(g Grid) (outer, inner raster.Path) {
	if len(s.Points) != 1 {
		return nil, nil
	}
	center := s.Points[0]
	diameter := 0.7 * math.Min(float64(g.CellW), float64(g.CellH))
	return Circle(float64(center.X), float64(center.Y), (diameter+STROKE_WIDTH)*0.5),
		Circle(float64(center.X), float64(center.Y), (diameter-STROKE_WIDTH)*0.5)
}

type Rect struct{ Min, Max Point }

func Bounds(pp []Point) Rect {
	if len(pp) == 0 {
		return Rect{}
	}
	r := Rect{pp[0], pp[0]}
	for _, p := range pp {
		if p.X < r.Min.X {
			r.Min.X = p.X
		}
		if p.X > r.Max.X {
			r.Max.X = p.X
		}
		if p.Y < r.Min.Y {
			r.Min.Y = p.Y
		}
		if p.Y > r.Max.Y {
			r.Max.Y = p.Y
		}
	}
	return r
}

func specPoints(bb Rect) (p1, p2, p3, p4 Point) {
	p1 = Point{X: bb.Min.X, Y: bb.Min.Y}
	p2 = Point{X: bb.Max.X, Y: bb.Min.Y}
	p3 = Point{X: bb.Max.X, Y: bb.Max.Y}
	p4 = Point{X: bb.Min.X, Y: bb.Max.Y}
	return
}

func (s *Shape) makeDocumentPath() raster.Path {
	bb := Bounds(s.Points)
	p1, p2, p3, p4 := specPoints(bb)
	pmid := Point{X: 0.5 * (bb.Min.X + bb.Max.X), Y: bb.Max.Y}

	path := raster.Path{}
	path.Start(P(p1))
	path.Add1(P(p2))
	path.Add1(P(p3))

	controlDX := (bb.Max.X - bb.Min.X) / 6
	controlDY := (bb.Max.Y - bb.Min.Y) / 8
	path.Add2(P(Point{X: pmid.X + controlDX, Y: pmid.Y - controlDY}), P(pmid))
	path.Add2(P(Point{X: pmid.X - controlDX, Y: pmid.Y + controlDY}), P(p4))
	path.Add1(P(p1))
	return path
}

func (s *Shape) makeIOPath(g Grid, opt Options) raster.Path {
	if len(s.Points) != 4 {
		return nil
	}
	bb := Bounds(s.Points)
	p1, p2, p3, p4 := specPoints(bb)
	//TODO: handle opt.FixedSlope
	offset := float64(g.CellW) * 0.5

	path := raster.Path{}
	path.Start(P(Point{X: p1.X + offset, Y: p1.Y}))
	path.Add1(P(Point{X: p2.X + offset, Y: p2.Y}))
	path.Add1(P(Point{X: p3.X - offset, Y: p3.Y}))
	path.Add1(P(Point{X: p4.X - offset, Y: p4.Y}))
	path.Add1(P(Point{X: p1.X + offset, Y: p1.Y})) // close path
	return path
}

func (s *Shape) makeTrapezoidPath(g Grid, opt Options, inverted bool) raster.Path {
	if len(s.Points) != 4 {
		return nil
	}
	bb := Bounds(s.Points)
	//TODO: handle opt.FixedSlope
	offset := float64(g.CellW) * 0.5
	if inverted {
		offset = -offset
	}
	ul := Point{X: bb.Min.X + offset, Y: bb.Min.Y}
	ur := Point{X: bb.Max.X - offset, Y: bb.Min.Y}
	br := Point{X: bb.Max.X + offset, Y: bb.Max.Y}
	bl := Point{X: bb.Min.X - offset, Y: bb.Max.Y}
	//pmid := Point{X:0.5*(bb.Min.X+bb.Max.X), Y:bb.Max.Y}

	path := raster.Path{}
	path.Start(P(ul))
	path.Add1(P(ur))
	path.Add1(P(br))
	path.Add1(P(bl))
	path.Add1(P(ul)) // close path
	return path
}

func (s *Shape) MakeIntoRenderPath(g Grid, opt Options) raster.Path {
	if s.Type == TYPE_POINT_MARKER {
		panic("please handle markers separately")
		return nil
		//return s.makeMarkerPath(g)
	}
	if len(s.Points) == 4 {
		switch s.Type {
		case TYPE_DOCUMENT:
			return s.makeDocumentPath()
		case TYPE_IO:
			return s.makeIOPath(g, opt)
		case TYPE_MANUAL_OPERATION:
			return s.makeTrapezoidPath(g, opt, true)
		case TYPE_TRAPEZOID:
			return s.makeTrapezoidPath(g, opt, false)
		case TYPE_STORAGE, TYPE_DECISION, TYPE_ELLIPSE:
			_ = fmt.Sprintf
			//panic(fmt.Sprintf("niy for type %d", s.Type))
			//TODO: fixme
			return nil
		}
	}
	if len(s.Points) < 2 {
		return nil
	}
	path := raster.Path{}
	point, prev, next := s.Points[0], s.Points[len(s.Points)-1], s.Points[1]
	_, _ = prev, next
	//var entry, exit *Point
	switch point.Type {
	case POINT_NORMAL:
		path.Start(P(point))
	case POINT_ROUND:
		//TODO: fixme
		path.Start(P(point))
		//panic("niy")
	}
	for i := 1; i < len(s.Points); i++ {
		prev = point
		point = s.Points[i]
		if i < len(s.Points)-1 {
			next = s.Points[i+1]
		} else {
			next = s.Points[0]
		}
		switch point.Type {
		case POINT_NORMAL:
			path.Add1(P(point))
		case POINT_ROUND:
			//TODO: fixme
			path.Add1(P(point))
			//panic("niy")
		}
	}
	if s.Closed && len(s.Points) > 2 {
		path.Add1(P(s.Points[0])) //FIXME: other for POINT_ROUND?
	}
	return path
}
