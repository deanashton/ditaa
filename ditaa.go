package main

import (
	"fmt"
	"os"
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
	grid := NewTextGrid()
	err = grid.LoadFrom(r)
	if err != nil {
		return err
	}
	diagram := NewDiagram(grid)
	image := NewBitmapRenderer().RenderToImage(diagram)
	// TODO: write .png to outfile
	_ = image
	return nil
}
