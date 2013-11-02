package main

import (
	"code.google.com/p/freetype-go/freetype/raster"
	"math"
	//"fmt"
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

func (s *Shape) makeDocumentPath() raster.Path {
	bb := Bounds(s.Points)
	p1 := Point{X: bb.Min.X, Y: bb.Min.Y}
	p2 := Point{X: bb.Max.X, Y: bb.Min.Y}
	p3 := Point{X: bb.Max.X, Y: bb.Max.Y}
	p4 := Point{X: bb.Min.X, Y: bb.Max.Y}
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
		case TYPE_STORAGE, TYPE_IO, TYPE_DECISION, TYPE_MANUAL_OPERATION, TYPE_TRAPEZOID, TYPE_ELLIPSE:
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
