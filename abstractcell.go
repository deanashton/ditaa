// WIP

package main

type AbstractCell [9]bool

func (c *AbstractCell) Get(x, y int) bool {
	return (*c)[y*3+x]
}

func abpix(source, mask int32) bool {
	switch source & mask {
	case 0:
		return false
	case mask:
		return true
	}
	panic("bad abstract pixel")
}

func PaintAbCell(hextop, hexmid, hexbot int32) AbstractCell {
	return AbstractCell{
		abpix(hextop, 0x100), abpix(hextop, 0x010), abpix(hextop, 0x001),
		abpix(hexmid, 0x100), abpix(hexmid, 0x010), abpix(hexmid, 0x001),
		abpix(hexbot, 0x100), abpix(hexbot, 0x010), abpix(hexbot, 0x001),
	}
}

var (
	abHLine   = PaintAbCell(0x000, 0x111, 0x000)
	abVLine   = PaintAbCell(0x010, 0x010, 0x010)
	abCorner1 = PaintAbCell(0x000, 0x011, 0x010)
	abCorner2 = PaintAbCell(0x000, 0x110, 0x010)
	abCorner3 = PaintAbCell(0x010, 0x110, 0x000)
	abCorner4 = PaintAbCell(0x010, 0x011, 0x000)
	abT       = PaintAbCell(0x000, 0x111, 0x010)
	abInvT    = PaintAbCell(0x010, 0x111, 0x000)
	abK       = PaintAbCell(0x010, 0x011, 0x010)
	abInvK    = PaintAbCell(0x010, 0x110, 0x010)
	abCross   = PaintAbCell(0x010, 0x111, 0x010)
	abStar    = PaintAbCell(0x111, 0x111, 0x111)
)
