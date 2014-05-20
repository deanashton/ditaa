package main

import (
	"github.com/akavel/ditaa/graphical"
)

func NewSmallLine(grid *TextGrid, c Cell, cellw, cellh int) *graphical.Shape {
	switch {
	case grid.IsHorizontalLine(c):
		return graphical.Shape{Points: []graphical.Point{
			{X: c.X * cellw, Y: c.Y*cellh + cellh/2},
			{X: c.X*cellw + cellw - 1, Y: c.Y*cellh + cellh/2},
		}}
	case grid.IsVerticalLine(c):
		return graphical.Shape{Points: []graphical.Point{
			{X: c.X*cellw + cellw/2, Y: c.Y * cellh},
			{X: c.X*cellw + cellw/2, Y: c.Y*cellh + cellh - 1},
		}}
	}
	return nil
}

func ConnectEndsToAnchors(s *Shape, grid *TextGrid, gg *graphical.Grid) {
	if s.IsClosed() {
		return
	}
	n := len(s.Points)
	for _, line := range []struct{ end, next *graphical.Point }{
		{&s.Points[0], &s.Points[1]},
		{&s.Points[n-1], &s.Points[n-2]},
	} {
		var x, y int
		switch {
		case line.next.NorthOf(line.end):
			x, y = line.end.X, line.end.Y+gg.CellH
		case line.next.SouthOf(line.end):
			x, y = line.end.X, line.end.Y-gg.CellH
		case line.next.WestOf(line.end):
			x, y = line.end.X+gg.CellW, line.end.Y
		case line.next.EastOf(line.end):
			x, y = line.end.X-gg.CellW, line.end.Y
		}
		anchor := gg.CellFor(graphical.Point{X: x, Y: y})
		if grid.IsArrowhead(anchor) || grid.IsCorner(anchor) || grid.IsIntersection(anchor) {
			line.end.X, line.end.Y = gg.CellMidX(anchor), gg.CellMidY(anchor)
			line.end.Locked = true
		}
	}
}
