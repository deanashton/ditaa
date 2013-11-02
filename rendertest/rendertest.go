package main

import (
	"encoding/xml"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

const (
	sources = "../orig-java/tests/xmls"
	results = "imgs"
)

type Ref struct {
	From string `xml:"reference,attr"`
}

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

type Color struct {
	R int `xml:"r,attr"`
	G int `xml:"g,attr"`
	B int `xml:"b,attr"`
	A int `xml:"a,attr"`
	Ref
}

type Point struct {
	X      float64 `xml:"x,attr"`
	Y      float64 `xml:"y,attr"`
	Locked bool    `xml:"locked,attr"`
	Type   int     `xml:"type,attr"`
}

type Shape struct {
	Type        int     `xml:"type"`
	FillColor   *Color  `xml:"fillColor"`
	StrokeColor Color   `xml:"strokeColor"`
	Closed      bool    `xml:"isClosed"`
	Dashed      bool    `xml:"isStrokeDashed"`
	Points      []Point `xml:"points>point"`
}

type Diagram struct {
	XMLName xml.Name `xml:"diagram"`
	Grid    Grid     `xml:"grid"`
	Shapes  []Shape  `xml:"shapes>shape"`
	//TODO: []Text
}

type Options struct{}

func LoadDiagram(path string) (*Diagram, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("loading diagram from '%s': %s", path, err)
	}
	defer r.Close()

	diagram := Diagram{}
	err = xml.NewDecoder(r).Decode(&diagram)
	if err != nil {
		return nil, fmt.Errorf("decoding diagram from '%s': %s", path, err)
	}
	//panic(fmt.Sprintf("%s: %#v", path, diagram))
	return &diagram, nil
}

func RenderDiagram(img *image.NRGBA, diagram *Diagram, opt Options) error {
	return nil
}

func runRender(src, dst string) error {
	diagram, err := LoadDiagram(src)
	if err != nil {
		return err
	}
	img := image.NewNRGBA(image.Rect(0, 0, diagram.Grid.W, diagram.Grid.H))
	err = RenderDiagram(img, diagram, Options{})
	if err != nil {
		return err
	}
	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	err = png.Encode(w, img)
	return err
}

func run() error {
	fnames := []string{}

	os.Mkdir(results, 0666)

	//todo: load files from ../orig-java/tests/xmls/*.xml, then try to render them into some output dir, and link them all on one html page
	err := filepath.Walk(sources, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".xml" {
			return nil
		}
		fnames = append(fnames, info.Name())
		return runRender(path, filepath.Join(results, info.Name()+".png"))
	})

	if err != nil {
		return err
	}

	return err
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
