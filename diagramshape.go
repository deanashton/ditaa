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

func createOpenFromBoundaryCells(grid *TextGrid, cells *CellSet, cellw, cellh int, allCornersRound bool) []graphical.Shape {
	if cells.Type(grid) != SET_OPEN {
		panic("CellSet is closed and cannot be handled by this method")
	}
	if len(cells.Set) == 0 {
		return []graphical.Shape{}
	}

	shapes := []graphical.Shape{}
	workGrid := NewTextGrid(grid.Width(), grid.Height())
	CopySelectedCells(workGrid, cells, grid)

	visited := NewCellSet()
	for c := range cells.Set {
		if workGrid.IsLinesEnd(c) {
			nextCells := workGrid.FollowCell(c, nil)
			shapes = append(shapes, growEdgesFromCell(workGrid, cellw, cellh, allCornersRound, nextCells.SomeCell(), c, visited)...)
			break
		}
	}

	//dashed shapes should "infect" the rest of the shapes
	dashedShapeExists := false
	for _, s := range shapes {
		if s.Dashed {
			dashedShapeExists = true
			break
		}
	}
	if dashedShapeExists {
		for i := range shapes {
			shapes[i].Dashed = true
		}
	}

	return shapes
}

func makePointForCell(c Cell, grid *TextGrid, cellw, cellh int, allCornersRound bool) graphical.Point {
	var typ graphical.PointType
	switch {
	case grid.IsCorner(c) && allCornersRound:
		typ = POINT_ROUND
	case grid.IsNormalCorner(c):
		typ = POINT_NORMAL
	case grid.IsRoundCorner(c):
		typ = POINT_ROUND
	case grid.IsLinesEnd(c) || grid.IsIntersection(c):
		typ = POINT_NORMAL
	default:
		panic(fmt.Sprint("Cannot make point for cell", c))
	}
	return graphical.Point{
		X:    c.X*cellw + cellw/2,
		Y:    c.Y*cellh + cellh/2,
		Type: typ,
	}
}
