package main

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

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

const blankBorderSize = 2

var humanColorCodes = map[string]string{
	"GRE": "9D9",
	"BLU": "55B",
	"PNK": "FAA",
	"RED": "E32",
	"YEL": "FF3",
	"BLK": "000",
}

var markupTags = map[string]struct{}{
	"d":  struct{}{},
	"s":  struct{}{},
	"io": struct{}{},
	"c":  struct{}{},
	"mo": struct{}{},
	"tr": struct{}{},
	"o":  struct{}{},
}

var _SPACE = []byte{' '}

type TextGrid struct {
	Rows [][]rune
}

func NewTextGrid() *TextGrid {
	return &TextGrid{}
}

func CopyTextGrid(other *TextGrid) *TextGrid {
	t := TextGrid{}
	t.Rows = make([][]rune, len(other.Rows))
	for y, row := range other.Rows {
		t.Rows[y] = append([]rune(nil), row...)
	}
	return &t
}

func (t *TextGrid) LoadFrom(r io.Reader) error {
	lines := [][]rune{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := []rune(scanner.Text())
		lines = append(lines, line)
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// strip trailing blank lines
	for i := len(lines) - 1; i >= 0; i-- {
		if !onlyWhitespaces(lines[i]) {
			lines = lines[:i+1]
			break
		}
	}

	fixTabs(lines, DEFAULT_TAB_SIZE)
	t.Rows = lines

	// make all lines of equal length
	// add blank outline around the buffer to prevent fill glitch
	// convert tabs to spaces (or remove them if setting is 0)

	maxLen := 0
	for _, row := range t.Rows {
		if len(row) > maxLen {
			maxLen = len(row)
		}
	}

	newrows := make([][]rune, 0, len(t.Rows)+2*blankBorderSize)
	for i := 0; i < blankBorderSize; i++ {
		newrows = append(newrows, appendSpaces(nil, maxLen+2*blankBorderSize))
	}
	for _, row := range t.Rows {
		newrow := make([]rune, 0, maxLen+2*blankBorderSize)
		newrow = appendSpaces(newrow, blankBorderSize)
		newrow = append(newrow, row...)
		newrow = appendSpaces(newrow, cap(newrow)-len(newrow))
	}
	for i := 0; i < blankBorderSize; i++ {
		newrows = append(newrows, appendSpaces(nil, maxLen+2*blankBorderSize))
	}
	t.Rows = newrows

	t.replaceBullets()
	t.replaceHumanColorCodes()

	return nil
}

func onlyWhitespaces(rs []rune) bool {
	for _, r := range rs {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func fixTabs(rows [][]rune, tabSize int) {
	for y, row := range rows {
		newrow := make([]rune, 0, len(row))
		for _, c := range row {
			if c == '\t' {
				newrow = appendSpaces(newrow, tabSize-len(newrow)%tabSize)
			} else {
				newrow = append(newrow, c)
			}
		}
		rows[y] = newrow
	}
}

func appendSpaces(row []rune, n int) []rune {
	for i := 0; i < n; i++ {
		row = append(row, ' ')
	}
	return row
}

func (t *TextGrid) replaceBullets() {
	for y, row := range t.Rows {
		for x, _ := range row {
			if t.IsBullet(x, y) {
				t.Set(x, y, ' ')
				t.Set(x+1, y, '\u2022')
			}
		}
	}
}

func (t *TextGrid) replaceHumanColorCodes() {
	for y, row := range t.Rows {
		s := string(row)
		for k, v := range humanColorCodes {
			k, v = "c"+k, "c"+v
			s = strings.Replace(s, k, v, -1)
		}
		t.Rows[y] = []rune(s)
	}
}

func (t *TextGrid) IsBullet(x, y int) bool {
	ch := t.Get(x, y)
	return (ch == 'o' || ch == '*') &&
		t.IsBlankNon0(x+1, y) &&
		t.IsBlankNon0(x-1, y) &&
		isAlphNum(t.Get(x+2, y))
}

func (t *TextGrid) IsBlankNon0(x, y int) bool { return t.Get(x, y) == ' ' }

func isAlphNum(ch rune) bool { return unicode.IsLetter(ch) || unicode.IsDigit(ch) }

func (t *TextGrid) Set(x, y int, ch rune) { t.Rows[y][x] = ch }
func (t *TextGrid) Get(x, y int) rune     { return t.Rows[y][x] }

type Cell struct{ X, Y int }
