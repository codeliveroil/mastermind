// Copyright (c) 2018 codeliveroil. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package main

import "fmt"

const (
	clrBoard     = 3
	clrCorrect   = 15
	clrBadLoc    = 8
	clrHiddenPeg = 244
	boardPeg     = "●"
	debug        = false
	boardHeight  = maxTries + 4
)

// clear clears enough room to render the board and results and paints
// the background black so that the experience is the same regardless
// of the terminal color
func (b *board) clear() {
	write("\n", clrBoard)
	l := ""
	for i := 0; i < (b.codeSize*2+1)+(b.codeSize+5); i++ {
		l += " "
	}
	l += "\n"
	for i := 0; i < boardHeight; i++ {
		write(l, clrBoard)
	}
}

// draw renders the board onto the terminal.
//
// draw uses regular characters instead of line drawing chars because the
// latter introduces a slowdown each time the board is redrawn.
func (b *board) draw() {
	for i := 0; i < boardHeight; i++ {
		if !debug {
			lineUp()
		}
	}

	// top bar
	var bar string
	for i := 0; i < b.codeSize*2+1; i++ {
		bar += "="
	}
	bar += "\n"
	write(bar, clrBoard)

	printRow := func(r row) {
		for i := 0; i < b.codeSize; i++ {
			write("|", clrBoard)
			if p := r.pegs[i]; p == nullPeg {
				write("_", clrBoard)
			} else {
				write(boardPeg, p.color)
			}
			if i == b.codeSize-1 {
				write("| ", clrBoard)
			}
		}

		makeResult := func(c int) string {
			var res string
			if c%2 != 0 {
				res += "."
			}
			c--
			for i := 0; i < c; i = i + 2 {
				res += ":"
			}
			return res
		}

		write(makeResult(r.result.correct), clrCorrect)
		write(makeResult(r.result.badLoc), clrBadLoc)
		write("\n", clrBoard)
	}

	// code
	if b.won() || b.lost || debug {
		printRow(b.code)
	} else {
		for i := 0; i < b.codeSize; i++ {
			write("|", clrBoard)
			write(boardPeg, clrHiddenPeg)
			if i == b.codeSize-1 {
				write("|\n", clrBoard)
			}
		}
	}
	write(bar, clrBoard)

	// remaining attempts
	for i := 0; i < maxTries-b.currRow-1; i++ {
		r := *newRow(b.codeSize)
		r.pegs[0] = nullPeg
		printRow(r)
	}

	// current and existing attempts
	for i := len(b.rows) - 1; i >= 0; i-- {
		printRow(*b.rows[i])
	}

	// bottom bar
	write(bar, clrBoard)
}

// write writes the given string to the console with the foreground
// color specified in 'fgColor'
func write(str string, fgColor uint8) {
	nl := false
	if l := len(str); l > 0 && str[l-1] == '\n' {
		nl = true
		str = str[:l-1]
	}
	fmt.Print(fmt.Sprintf("\x1b[48;5;0m\x1b[38;5;%vm%v\x1b[0m", fgColor, str))

	if nl { //don't format the new line because we don't want the the black background to extend to the width of the terminal
		fmt.Print("\n")
	}
}
