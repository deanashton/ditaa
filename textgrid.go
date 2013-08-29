package main

import (
	"bufio"
	"bytes"
	"io"
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
	Rows [][]byte // FIXME: change to [][]rune and handle UTF-8 input
}

func NewTextGrid() *TextGrid {
	return &TextGrid{}
}

func (t *TextGrid) LoadFrom(r io.Reader) error {
	lines := [][]byte{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := append([]byte(nil), scanner.Bytes()...)
		lines = append(lines, line)
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// strip trailing blank lines
	for i := len(lines) - 1; i >= 0; i-- {
		tmp := bytes.TrimSpace(lines[i])
		if len(tmp) != 0 {
			lines = lines[:i+1]
			break
		}
	}

	t.Rows = lines

	t.fixTabs(DEFAULT_TAB_SIZE)

	// make all lines of equal length
	// add blank outline around the buffer to prevent fill glitch
	// convert tabs to spaces (or remove them if setting is 0)

	maxLen := 0
	for _, row := range t.Rows {
		if len(row) > maxLen {
			maxLen = len(row)
		}
	}

	newrows := make([][]byte, 0, len(t.Rows)+2*blankBorderSize)
	for i := 0; i < blankBorderSize; i++ {
		newrows = append(newrows, bytes.Repeat(_SPACE, maxLen+2*blankBorderSize))
	}
	border := bytes.Repeat(_SPACE, blankBorderSize)
	for _, row := range t.Rows {
		newrow := make([]byte, 0, maxLen+2*blankBorderSize)
		newrow = append(newrow, border...)
		newrow = append(newrow, row...)
		for i := len(newrow); i < cap(newrow); i++ {
			newrow = append(newrow, ' ')
		}
	}
	for i := 0; i < blankBorderSize; i++ {
		newrows = append(newrows, bytes.Repeat(_SPACE, maxLen+2*blankBorderSize))
	}
	t.Rows = newrows

	t.replaceBullets()
	t.replaceHumanColorCodes()

	return nil
}

func (t *TextGrid) fixTabs(tabSize int) {
	for y, row := range t.Rows {
		newrow := make([]byte, 0, len(row))
		for _, c := range row {
			if c == '\t' {
				newrow = append(newrow, bytes.Repeat(_SPACE, tabSize-len(newrow)%tabSize)...)
			} else {
				newrow = append(newrow, c)
			}
		}
		t.Rows[y] = newrow
	}
}
