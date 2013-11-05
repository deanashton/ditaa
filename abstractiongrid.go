package main

type AbstractionGrid struct {
	Rows [][]rune
}

func NewAbstractionGrid(t *TextGrid, cells *CellSet) *AbstractionGrid {
	g := &AbstractionGrid{}
	g.reset(3*t.Width(), 3*t.Height())
	for c := range cells.Set {
		if t.IsBlank(c) {
			continue
		}
		for _, x := range abstractionChecks {
			if x.check(t, c) {
				g.Set(c, x.result)
				break
			}
		}
	}
	return g
}

func (g *AbstractionGrid) reset(w, h int) {
	g.Rows = make([][]rune, h)
	for y := range g.Rows {
		g.Rows[y] = make([]rune, w)
		for x := range g.Rows[y] {
			g.Rows[y][x] = ' '
		}
	}
}

func (g *AbstractionGrid) Set(c Cell, brush AbstractCell) {
	x, y := 3*c.X, 3*c.Y
	for dy := 0; dy < 3; dy++ {
		for dx := 0; dx < 3; dx++ {
			if brush.Get(dx, dy) {
				g.Rows[y+dy][x+dx] = '*'
			}
		}
	}
}

var abstractionChecks = []struct {
	check  func(*TextGrid, Cell) bool
	result AbstractCell
}{
	{(*TextGrid).IsCross, abCross},
	{(*TextGrid).IsT, abT},
	{(*TextGrid).IsK, abK},
	{(*TextGrid).IsInverseT, abInvT},
	{(*TextGrid).IsInverseK, abInvK},
	{(*TextGrid).IsCorner1, abCorner1},
	{(*TextGrid).IsCorner2, abCorner2},
	{(*TextGrid).IsCorner3, abCorner3},
	{(*TextGrid).IsCorner4, abCorner4},
	{(*TextGrid).IsHorizontalLine, abHLine},
	{(*TextGrid).IsVerticalLine, abVLine},
	{(*TextGrid).IsCrossOnLine, abCross},
	{(*TextGrid).IsStarOnLine, abStar},
}
