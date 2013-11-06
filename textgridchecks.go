package main

import (
	"strings"
	"unicode"
)

type Cell struct{ X, Y int }

func (c Cell) North() Cell { return Cell{c.X, c.Y - 1} }
func (c Cell) South() Cell { return Cell{c.X, c.Y + 1} }
func (c Cell) East() Cell  { return Cell{c.X + 1, c.Y} }
func (c Cell) West() Cell  { return Cell{c.X - 1, c.Y} }

type CellSet struct {
	Set map[Cell]struct{}
}
type CellBounds struct{ Min, Max Cell }

func NewCellSet() *CellSet {
	return &CellSet{Set: make(map[Cell]struct{})}
}
func (s *CellSet) Add(c Cell) { s.Set[c] = struct{}{} }
func (s *CellSet) Bounds() CellBounds {
	var bb *CellBounds
	for c := range s.Set {
		if bb == nil {
			bb = &CellBounds{Min: c, Max: c}
			continue
		}
		if c.X < bb.Min.X {
			bb.Min.X = c.X
		}
		if c.X > bb.Max.X {
			bb.Max.X = c.X
		}
		if c.Y < bb.Min.Y {
			bb.Min.Y = c.Y
		}
		if c.Y > bb.Max.Y {
			bb.Max.Y = c.Y
		}
	}
	return *bb
}

func isAlphNum(ch rune) bool             { return unicode.IsLetter(ch) || unicode.IsDigit(ch) }
func isOneOf(ch rune, group string) bool { return strings.ContainsRune(group, ch) }

func (t *TextGrid) IsBullet(x, y int) bool {
	ch := t.Get(x, y)
	return (ch == 'o' || ch == '*') &&
		t.IsBlankNon0(x+1, y) &&
		t.IsBlankNon0(x-1, y) &&
		isAlphNum(t.Get(x+2, y))
}

func (t *TextGrid) IsOutOfBounds(c Cell) bool {
	return c.X >= t.Width() || c.Y >= t.Height() || c.X < 0 || c.Y < 0
}

func (t *TextGrid) IsBlankNon0(x, y int) bool { return t.Get(x, y) == ' ' }
func (t *TextGrid) IsBlank(c Cell) bool {
	ch := t.Get(c.X, c.Y)
	if ch == 0 {
		return false // FIXME: should this be false, or true (see 'isBlank(x,y)' in Java)
	}
	return ch == ' '
}
func (t *TextGrid) IsBlankXY(x, y int) bool {
	ch := t.Get(x, y)
	if ch == 0 {
		return true
	}
	return ch == ' '
}

func (t *TextGrid) IsBoundary(c Cell) bool {
	ch := t.Get(c.X, c.Y)
	switch ch {
	case 0:
		return false
	case '+', '\\', '/':
		return t.IsIntersection(c) ||
			t.IsCorner(c) ||
			t.IsStub(c) ||
			t.IsCrossOnLine(c)
	}
	return isOneOf(ch, text_boundaries) && !t.IsLoneDiagonal(c)
}

func (t *TextGrid) IsIntersection(c Cell) bool {
	return intersectionCriteria.AnyMatch(t.TestingSubGrid(c))
}
func (t *TextGrid) IsNormalCorner(c Cell) bool {
	return normalCornerCriteria.AnyMatch(t.TestingSubGrid(c))
}
func (t *TextGrid) IsRoundCorner(c Cell) bool {
	return roundCornerCriteria.AnyMatch(t.TestingSubGrid(c))
}
func (t *TextGrid) IsStub(c Cell) bool {
	return stubCriteria.AnyMatch(t.TestingSubGrid(c))
}
func (t *TextGrid) IsCrossOnLine(c Cell) bool {
	return crossOnLineCriteria.AnyMatch(t.TestingSubGrid(c))
}
func (t *TextGrid) IsLoneDiagonal(c Cell) bool {
	return loneDiagonalCriteria.AnyMatch(t.TestingSubGrid(c))
}
func (t *TextGrid) IsCross(c Cell) bool      { return crossCriteria.AnyMatch(t.TestingSubGrid(c)) }
func (t *TextGrid) IsT(c Cell) bool          { return TCriteria.AnyMatch(t.TestingSubGrid(c)) }
func (t *TextGrid) IsK(c Cell) bool          { return KCriteria.AnyMatch(t.TestingSubGrid(c)) }
func (t *TextGrid) IsInverseT(c Cell) bool   { return inverseTCriteria.AnyMatch(t.TestingSubGrid(c)) }
func (t *TextGrid) IsInverseK(c Cell) bool   { return inverseKCriteria.AnyMatch(t.TestingSubGrid(c)) }
func (t *TextGrid) IsCorner1(c Cell) bool    { return corner1Criteria.AnyMatch(t.TestingSubGrid(c)) }
func (t *TextGrid) IsCorner2(c Cell) bool    { return corner2Criteria.AnyMatch(t.TestingSubGrid(c)) }
func (t *TextGrid) IsCorner3(c Cell) bool    { return corner3Criteria.AnyMatch(t.TestingSubGrid(c)) }
func (t *TextGrid) IsCorner4(c Cell) bool    { return corner4Criteria.AnyMatch(t.TestingSubGrid(c)) }
func (t *TextGrid) IsStarOnLine(c Cell) bool { return starOnLineCriteria.AnyMatch(t.TestingSubGrid(c)) }

// func (t *TextGrid) Is(c Cell) bool { return .AnyMatch(t.TestingSubGrid(c))}
func (t *TextGrid) IsCorner(c Cell) bool { return t.IsNormalCorner(c) || t.IsRoundCorner(c) }
func (t *TextGrid) IsHorizontalLine(c Cell) bool {
	ch := t.Get(c.X, c.Y)
	if ch == 0 {
		return false
	}
	return isOneOf(ch, text_horizontalLines)
}
func (t *TextGrid) IsVerticalLine(c Cell) bool {
	ch := t.Get(c.X, c.Y)
	if ch == 0 {
		return false
	}
	return isOneOf(ch, text_verticalLines)
}

const (
	text_boundaries             = `/\|-*=:`
	text_undisputableBoundaries = `|-*=:`
	text_horizontalLines        = `-=`
	text_verticalLines          = `|:`
	text_arrowHeads             = `<>^vV`
	text_cornerChars            = `\/+`
	text_pointMarkers           = `*`
	text_dashedLines            = `:~=`
	text_entryPoints1           = `\`
	text_entryPoints2           = `|:+\/`
	text_entryPoints3           = `/`
	text_entryPoints4           = `-=+\/`
	text_entryPoints5           = `\`
	text_entryPoints6           = `|:+\/`
	text_entryPoints7           = `/`
	text_entryPoints8           = `-=+\/`
)
