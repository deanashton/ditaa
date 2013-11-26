package graphical

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	"code.google.com/p/graphics-go/graphics"
	"code.google.com/p/graphics-go/graphics/interp"
	"encoding/xml"
	"image"
	"image/color"
	"io/ioutil"
	"sort"
)

const (
	fontpath = "../orig-java/src/org/stathissideris/ascii2image/graphics/font.ttf"
)

type Label struct {
	Text         string  `xml:"text"`
	FontSize     float64 `xml:"font>size"`
	X            int     `xml:"xPos"`
	Y            int     `xml:"yPos"`
	Color        Color   `xml:"color"`
	OnLine       bool    `xml:"isTextOnLine"`
	Outline      bool    `xml:"hasOutline"`
	OutlineColor Color   `xml:"outlineColor"`
}

type Diagram struct {
	XMLName xml.Name `xml:"diagram"`
	Grid    Grid     `xml:"grid"`
	Shapes  []Shape  `xml:"shapes>shape"`
	Labels  []Label  `xml:"texts>text"`
}

type Options struct {
	DropShadows bool
}

func renderShadows(img *image.RGBA, shapes []Shape, g Grid, opt Options) {
	for _, shape := range shapes {
		if len(shape.Points) == 0 || !shape.DropsShadow() || shape.Type == TYPE_CUSTOM {
			continue
		}
		path := shape.MakeIntoRenderPath(g /*, opt*/)
		if path == nil {
			continue
		}
		Fill(img, path, color.RGBA{150, 150, 150, 255})
	}
	offset := g.CellW
	if g.CellH < offset {
		offset = g.CellH
	}
	offsetf := float64(offset) / 3.3333
	img2 := image.NewRGBA(img.Bounds())
	graphics.I.Translate(float64(offsetf), float64(offsetf)).Transform(img2, img, interp.Bilinear)
	*img = *img2
}

func blurShadows(img *image.RGBA) {
	radius := 4
	StackBlur(img, radius, true)

	// remove blur artifacts from the top-left border of image
	bb := img.Rect
	radius += 2
	for y := bb.Min.Y; y <= bb.Min.Y+radius; y++ {
		for x := bb.Min.X; x <= bb.Max.X; x++ {
			img.SetRGBA(x, y, WHITE)
		}
	}
	for y := bb.Min.Y + radius + 1; y <= bb.Max.Y; y++ {
		for x := bb.Min.X; x <= bb.Min.X+radius; x++ {
			img.SetRGBA(x, y, WHITE)
		}
	}
}

type LargeFirst []Shape

func (t LargeFirst) Len() int           { return len(t) }
func (t LargeFirst) Less(i, j int) bool { return t[i].CalcArea() > t[j].CalcArea() }
func (t LargeFirst) Swap(i, j int) {
	tmp := t[i]
	t[i] = t[j]
	t[j] = tmp
}

func RenderDiagram(img *image.RGBA, diagram *Diagram, opt Options) error {
	fontfile, err := ioutil.ReadFile(fontpath)
	if err != nil {
		return err
	}
	font, err := truetype.Parse(fontfile)
	if err != nil {
		return err
	}
	_ = font

	for y := 0; y < diagram.Grid.H; y++ {
		for x := 0; x < diagram.Grid.W; x++ {
			img.SetRGBA(x, y, WHITE)
		}
	}

	//TODO: antialiasing options

	// drop shadows
	if opt.DropShadows {
		renderShadows(img, diagram.Shapes, diagram.Grid, opt)

		//TODO: blur shadows
		if true {
			blurShadows(img)
		}
	}

	//render storage shapes
	//special case since they are '3d' and should be
	//rendered bottom to top
	//TODO: known bug: if a storage object is within a bigger normal box, it will be overwritten in the main drawing loop
	//(BUT this is not possible since tags are applied to all shapes overlaping shapes)
	storageShapes := []Shape{}
	for _, shape := range diagram.Shapes {
		if shape.Type == TYPE_STORAGE {
			//TODO: freetype-go doesn't implement stroking cubic paths -- need to fix or walk around
			//storageShapes = append(storageShapes, shape)
		}
	}
	//TODO: sort storage shapes
	for _, shape := range storageShapes {
		path := shape.MakeIntoRenderPath(diagram.Grid /*, opt*/)
		//TODO: handle dashed
		color := WHITE
		if shape.FillColor != nil {
			color = shape.FillColor.RGBA()
		}
		Fill(img, path, color)
		Stroke(img, path, shape.StrokeColor.RGBA())
	}

	sort.Sort(LargeFirst(diagram.Shapes))

	// render rest of shapes + collect point markers
	pointMarkers := []Shape{}
	for _, shape := range diagram.Shapes {
		switch shape.Type {
		case TYPE_POINT_MARKER:
			pointMarkers = append(pointMarkers, shape)
			continue
		case TYPE_STORAGE:
			continue
		case TYPE_CUSTOM:
			//TODO: render custom shape
			continue
		}
		if len(shape.Points) == 0 {
			continue
		}

		path := shape.MakeIntoRenderPath(diagram.Grid /*, opt*/)

		// fill
		if path != nil && shape.Closed && !shape.Dashed {
			color := WHITE
			if shape.FillColor != nil {
				color = shape.FillColor.RGBA()
			}
			Fill(img, path, color)
		}

		// draw
		if shape.Type != TYPE_ARROWHEAD {
			//TODO: support dashed lines
			Stroke(img, path, shape.StrokeColor.RGBA())
		}
	}

	// render point markers
	for _, shape := range pointMarkers {
		outer, inner := shape.MakeMarkerPaths(diagram.Grid)
		Fill(img, outer, shape.StrokeColor.RGBA())
		Fill(img, inner, WHITE)
	}

	// handle text
	for _, label := range diagram.Labels {
		ctx := freetype.NewContext()
		ctx.SetFont(font)
		ctx.SetFontSize(label.FontSize)
		ctx.SetSrc(image.NewUniform(label.Color.RGBA()))
		ctx.SetDst(img)
		ctx.SetClip(img.Bounds())
		//TODO: handle outline
		ctx.DrawString(label.Text, P(Point{X: float64(label.X), Y: float64(label.Y)}))
	}
	return nil
}
