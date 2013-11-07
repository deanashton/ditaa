package main

type Diagram struct {
	W, H int
}

/*
An outline of the inner workings of this very important (and monstrous)
constructor is presented here. Boundary processing is the first step
of the process:

  1. Copy the grid into a work grid and remove all type-on-line
     and point markers from the work grid
  2. Split grid into distinct shapes by plotting the grid
     onto an AbstractionGrid and its getDistinctShapes() method.
  3. Find all the possible boundary sets of each of the
     distinct shapes. This can produce duplicate shapes (if the boundaries
     are the same when filling from the inside and the outside).
  4. Remove duplicate boundaries.
  5. Remove obsolete boundaries. Obsolete boundaries are the ones that are
     the sum of their parts when plotted as filled shapes. (see method
     removeObsoleteShapes())
  6. Seperate the found boundary sets to open, closed or mixed
     (See CellSet class on how its done).
  7. Are there any closed boundaries?
        * YES. Subtract all the closed boundaries from each of the
          open ones. That should convert the mixed shapes into open.
        * NO. In this (harder) case, we use the method
          breakTrulyMixedBoundaries() of CellSet to break boundaries
          into open and closed shapes (would work in any case, but it's
          probably slower than the other method). This method is based
          on tracing from the lines' ends and splitting when we get to
          an intersection.
  8. If we had to eliminate any mixed shapes, we seperate the found
     boundary sets again to open, closed or mixed.

At this stage, the boundary processing is all complete and we
proceed with using those boundaries to create the shapes:

  1. Create closed shapes.
  2. Create open shapes. That's when the line end corrections are
     also applied, concerning the positioning of the ends of lines
     see methods connectEndsToAnchors() and moveEndsToCellEdges() of
     DiagramShape.
  3. Assign color codes to closed shapes.
  4. Assign extended markup tags to closed shapes.
  5. Create arrowheads.
  6. Create point markers.

Finally, the text processing occurs: [pending]

*/
func NewDiagram(grid *TextGrid) *Diagram {
	d := Diagram{}
	d.W, d.H = len(grid.Rows[0])*CELL_WIDTH, len(grid.Rows)*CELL_HEIGHT

	workGrid := CopyTextGrid(grid)
	//TODO: workGrid.replaceTypeOnLine()
	//TODO: workGrid.replacePointMarkersOnLine()

	boundaries := getAllBoundaries(workGrid)
	boundarySetsStep1 := getDistinctShapes(NewAbstractionGrid(workGrid, boundaries))
	_ = boundarySetsStep1

	if DEBUG {
		println("******* Distinct shapes found using AbstractionGrid *******")
		for _, cells := range boundarySetsStep1 {
			cells.printAsGrid()
		}
	}

	//Find all the boundaries by using the special version of the filling method
	//(fills in a different buffer than the buffer it reads from)
	w, h := grid.Width(), grid.Height()
	boundarySetsStep2 := []*CellSet{}
	for _, cells := range boundarySetsStep1 {
		//the fill buffer keeps track of which cells have been
		//filled already
		fillBuffer := NewTextGrid()
		fillBuffer.Rows = BlankRows(3*w, 3*h)

		for yi := 0; yi < 3*h; yi++ {
			for xi := 0; xi < 3*w; xi++ {
				if !fillBuffer.IsBlankXY(xi, yi) {
					continue
				}

				copyGrid := NewTextGrid()
				copyGrid.Rows = NewAbstractionGrid(workGrid, cells).Rows

				boundaries := findBoundariesExpandingFrom(copyGrid, Cell{xi, yi})
				if len(boundaries.Set) == 0 {
					continue //i'm not sure why these occur
				}
				boundarySetsStep2 = append(boundarySetsStep2, makeScaledOneThirdEquivalent(boundaries))

				copyGrid.Rows = NewAbstractionGrid(workGrid, cells).Rows
				filled := copyGrid.fillContinuousArea(xi, yi, '*')
				FillCellsWith(fillBuffer.Rows, filled, '*')
				FillCellsWith(fillBuffer.Rows, boundaries, '-')

				if DEBUG {
					makeScaledOneThirdEquivalent(boundaries).printAsGrid()
					println("----------------------------------------")
				}
			}
		}
	}

	//TODO: rest...

	return &d
}

func makeScaledOneThirdEquivalent(cells *CellSet) *CellSet {
	bb := cells.Bounds()
	gridBig := NewTextGrid()
	gridBig.Rows = BlankRows(bb.Max.X+2, bb.Max.Y+2)
	FillCellsWith(gridBig.Rows, cells, '*')

	gridSmall := NewTextGrid()
	gridSmall.Rows = BlankRows((bb.Max.X+2)/3, (bb.Max.Y+2)/3)
	for y := 0; y < gridBig.Height(); y++ {
		for x := 0; x < gridBig.Width(); x++ {
			if !gridBig.IsBlank(Cell{x, y}) {
				gridSmall.Set(x/3, y/3, '*')
			}
		}
	}
	return gridSmall.GetAllNonBlank()
}

func findBoundariesExpandingFrom(grid *TextGrid, seed Cell) *CellSet {
	boundaries := NewCellSet()
	if grid.IsOutOfBounds(seed) {
		return boundaries
	}
	oldChar := grid.GetCell(seed)
	newChar := rune(1) //TODO: kludge
	stack := []Cell{seed}
	expand := func(c Cell) {
		switch grid.GetCell(c) {
		case oldChar:
			stack = append(stack, c)
		case '*':
			boundaries.Add(c)
		}
	}
	for len(stack) > 0 {
		var c Cell
		c, stack = stack[len(stack)-1], stack[:len(stack)-1]
		grid.SetCell(c, newChar)
		expand(c.North())
		expand(c.South())
		expand(c.East())
		expand(c.West())
	}
	return boundaries
}

func getDistinctShapes(g *AbstractionGrid) []*CellSet {
	result := []*CellSet{}

	tg := TextGrid{Rows: g.Rows}
	nonBlank := tg.GetAllNonBlank()

	distinct := breakIntoDistinctBoundaries(nonBlank)
	for _, set := range distinct {
		temp := EmptyAbstractionGrid(g.Width(), g.Height())
		FillCellsWith(temp.Rows, set, '*')
		result = append(result, temp.GetAsTextGrid().GetAllNonBlank())
	}
	return result
}

func breakIntoDistinctBoundaries(cells *CellSet) []*CellSet {
	result := []*CellSet{}
	bb := cells.Bounds()
	boundaryGrid := NewTextGrid()
	boundaryGrid.Rows = BlankRows(bb.Max.X+2, bb.Max.Y+2)
	FillCellsWith(boundaryGrid.Rows, cells, '*')

	for c := range cells.Set {
		if boundaryGrid.IsBlankXY(c.X, c.Y) {
			continue
		}
		boundarySet := boundaryGrid.fillContinuousArea(c.X, c.Y, ' ')
		result = append(result, boundarySet)
	}
	return result
}

func getAllBoundaries(g *TextGrid) *CellSet {
	set := NewCellSet()
	for y, row := range g.Rows {
		for x := range row {
			c := Cell{x, y}
			if g.IsBoundary(c) {
				set.Add(c)
			}
		}
	}
	return set
}
