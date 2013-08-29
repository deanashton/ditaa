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
	workGrid.replaceTypeOnLine()
	workGrid.replacePointMarkersOnLine()

	return &d
}
