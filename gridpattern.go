//WIP

package main

/*
This is a TextGrid (usually 3x3) that contains the equivalent
of a 2D reqular expression (which uses custom syntax to make
things more visual, but standard syntax is also possible).

The custom syntax is:

	. means anything
	b means any boundary (any of - = / \ + | :)
	! means not boundary (none of - = / \ + | :)
	- means - or =
	| means | or :
	[ means not | nor :
	~ means not - nor =
	^ means a boundary but not - nor =
	( means a boundary but not | nor :
	s means a straight boundary (one of - = + | :)
	S means not a straight boundary (none of - = + | :)

	1 means a cell that has entry point 1
	2 means a cell that has entry point 2
	3 means a cell that has entry point 3
	etc. up to number 8

	%1 means a cell that does not have entry point 1 etc.

See below for an explanation of entry points

+, \, / and the space are literal (as is any other character)


Entry points

	1   2   3
	 *--*--*
	 |     |
	8*     *4
	 |     |
	 *--*--*
	7   6   5

We number the entry points for each cell as in the diagram above.
If a cell is occupied by a character, we define as entry points
the points of the above diagram that the character can touch with
the end of its lines. For example:

	- has entry points 8 and 4,
	| and : have entry points 2 and 6,
	/ has 3 and 7,
	\ has 1 and 5,
	+ has 2, 6, 8 and 4
	etc.

*/
var gridPatterns = map[byte]string{
	'[':  "[^|:]",
	'|':  "[|:]",
	'-':  "[-=]",
	'!':  "[^-=\\/\\\\+|:]",
	'b':  "[-=\\/\\\\+|:]",
	'^':  "[\\/\\\\+|:]",
	'(':  "[-=\\/\\\\+]",
	'~':  ".",
	'+':  "\\+",
	'\\': "\\\\",
	's':  "[-=+|:]",
	'S':  "[\\/\\\\]",
	'*':  "\\*",

	//entry points
	'1': "[\\\\]",
	'2': "[|:+\\/\\\\]",
	'3': "[\\/]",
	'4': "[-=+\\/\\\\]",
	'5': "[\\\\]",
	'6': "[|:+\\/\\\\]",
	'7': "[\\/]",
	'8': "[-=+\\/\\\\]",
}

var gridPatternsInv = map[byte]string{
	'1': "[^\\\\]",
	'2': "[^|:+\\/\\\\]",
	'3': "[^\\/]",
	'4': "[^-=+\\/\\\\]",
	'5': "[^\\\\]",
	'6': "[^|:+\\/\\\\]",
	'7': "[^\\/]",
	'8': "[^-=+\\/\\\\]",
}
