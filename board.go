// Copyright (c) 2018 codeliveroil. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	minCodeSize = 4
	maxCodeSize = 10
	maxTries    = 10
	maxRepeats  = 10
)

var (
	nullPeg = peg{id: -1, color: 1, name: "nullPeg"}
	basePeg = allPegs[0]

	allPegs []peg = []peg{
		peg{id: 0, color: 9, name: "red"},
		peg{id: 1, color: 10, name: "green"},
		peg{id: 2, color: 11, name: "yellow"},
		peg{id: 3, color: 27, name: "blue"},
		peg{id: 4, color: 13, name: "pink"},
		peg{id: 5, color: 208, name: "orange"},
		peg{id: 6, color: 14, name: "cyan"},
		peg{id: 7, color: 7, name: "gray"},
	}
)

type board struct {
	codeSize int
	pegSet   []peg
	code     row
	rows     []*row
	currRow  int
	lost     bool
}

type row struct {
	pegs    []peg
	result  rowResult
	currPeg int
}

type peg struct {
	id    int
	color uint8
	name  string
}

type rowResult struct {
	correct int
	badLoc  int
}

// newRow returns a new row with all
// except the 0th peg set to nullPeg
// pegs
func newRow(codeSize int) *row {
	r := &row{}
	for i := 0; i < codeSize; i++ {
		r.pegs = append(r.pegs, nullPeg)
	}
	r.pegs[0] = basePeg
	return r
}

// newBoard initializes a new board with a code of 'codeSize',
// with 'uniquePegCount' unique pegs allowing 'maxRepeatPerPeg'
// repeats per peg
func newBoard(codeSize, uniquePegCount, maxRepeatPerPeg int) (*board, error) {
	if codeSize < minCodeSize || codeSize > maxCodeSize {
		return nil, errors.New(fmt.Sprintf("code size should be between %v and %v, both inclusive", minCodeSize, maxCodeSize))
	}
	if l := len(allPegs); uniquePegCount > l || uniquePegCount < 1 {
		return nil, errors.New(fmt.Sprintf("peg count should be between 1 and %v, both inclusive", l))
	}
	if maxRepeatPerPeg > maxRepeats || maxRepeatPerPeg < 1 {
		return nil, errors.New(fmt.Sprintf("repeat count should be between 1 and %v, both inclusive", maxRepeats))
	}
	if uniquePegCount*maxRepeatPerPeg < codeSize {
		return nil, errors.New("combination of peg count and repeat count does not yield sufficient pegs to play")
	}

	pegSet := allPegs[:uniquePegCount]

	// create random code
	rand.Seed(time.Now().Unix())
	var pegs []peg
	for i := 0; i < uniquePegCount; i++ {
		for j := 0; j < maxRepeatPerPeg; j++ {
			pegs = append(pegs, pegSet[i])
		}
	}
	var code row
	for i := 0; i < codeSize; i++ {
		p := rand.Intn(len(pegs))
		code.pegs = append(code.pegs, pegs[p])
		pegs = append(pegs[:p], pegs[p+1:]...)
	}

	b := &board{
		codeSize: codeSize,
		pegSet:   pegSet,
		code:     code,
		rows:     []*row{newRow(codeSize)},
		currRow:  0,
	}
	return b, nil
}

// row returns the current row being played.
func (b *board) row() *row {
	return b.rows[b.currRow]
}

// won returns true if the player won the game.
func (b *board) won() bool {
	return b.row().result.correct == b.codeSize
}

// selectPeg changes the current peg with the  next peg in the unique
// peg set. The peg to the right will be chosen if 'clockwise' is true
// and left otherwise. The selection will wrap if it hits the boundary of the
// peg set.
func (b *board) selectPeg(clockwise bool) {
	r := b.row()
	p := r.pegs[r.currPeg]
	var i int
	if clockwise {
		i = p.id + 1
		if i == len(b.pegSet) {
			i = 0
		}
	} else {
		i = p.id - 1
		if i == -1 {
			i = len(b.pegSet) - 1
		}
	}
	r.pegs[r.currPeg] = b.pegSet[i]
}

// confirmPeg confirms the peg selection and moves onto the next code index.
func (b *board) confirmPeg() {
	b.row().currPeg++
	if b.row().currPeg == b.codeSize {
		b.validate()
		if b.won() {
			return
		}
		if b.currRow == maxTries-1 {
			b.lost = true
			return
		}
		b.currRow++
		b.rows = append(b.rows, newRow(b.codeSize))
	}
	b.row().pegs[b.row().currPeg] = basePeg
}

// deletePeg deletes the current peg selection and moves to the previous peg.
func (b *board) deletePeg() {
	r := b.row()
	if r.currPeg == 0 {
		return
	}
	r.pegs[r.currPeg] = nullPeg
	r.currPeg--
}

// validate checks the state of the current row.
func (b *board) validate() {
	r := b.row()
	remCode := make(map[int]int)
	var remRow []peg
	for i, cp := range b.code.pegs {
		if rp := r.pegs[i]; rp == cp {
			r.result.correct++
		} else {
			count := remCode[cp.id]
			count++
			remCode[cp.id] = count
			remRow = append(remRow, rp)
		}
	}

	for _, p := range remRow {
		c, ok := remCode[p.id]
		if ok && c > 0 {
			r.result.badLoc++
			c--
			remCode[p.id] = c
		}
	}
}
