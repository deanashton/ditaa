package main

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/akavel/ditaa/graphical"
)

const (
	DEFAULT_TAB_SIZE = 8
	CELL_WIDTH       = 10
	CELL_HEIGHT      = 14
)

func main() {
	if len(os.Args[1:]) != 2 {
		fmt.Fprintf(os.Stderr, "USAGE: %s INFILE OUTFILE.png\n", os.Args[0])
		os.Exit(1)
	}

	err := run(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(2)
	}
}

func run(infile, outfile string) error {
	r, err := os.Open(infile)
	if err != nil {
		return err
	}
	grid := NewTextGrid(0, 0)
	err = grid.LoadFrom(r)
	if err != nil {
		return err
	}
	diagram := NewDiagram(grid)

	img := image.NewRGBA(image.Rect(0, 0, diagram.G.Grid.W, diagram.G.Grid.H))
	err = graphical.RenderDiagram(img, &diagram.G, graphical.Options{DropShadows: true}, "orig-java/src/org/stathissideris/ascii2image/graphics/font.ttf")
	if err != nil {
		return err
	}
	w, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer w.Close()

	wbuf := bufio.NewWriter(w)
	err = png.Encode(wbuf, img)
	if err != nil {
		return err
	}
	err = wbuf.Flush()
	return err
}
