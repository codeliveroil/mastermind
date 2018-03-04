// Copyright (c) 2018 codeliveroil. All rights reserved.
//
// This work is licensed under the terms of the MIT license.
// For a copy, see <https://opensource.org/licenses/MIT>.

package main

import "testing"

func errF(err error, t *testing.T) {
	if err != nil {
		t.Fatalf("Error running test: %v", err)
	}
}

func TestSelectPeg(t *testing.T) {
	b, err := newBoard(4, 4, 1)
	errF(err, t)
	b.selectPeg(true)
	if p := b.row().pegs[0]; p != b.pegSet[1] {
		t.Fatalf("expected: %v, got: %v", p, b.pegSet[1])
	}
	b.selectPeg(false)
	b.selectPeg(false)
	if p := b.row().pegs[0]; p != b.pegSet[b.codeSize-1] {
		t.Fatalf("expected: %v, got: %v", p, b.pegSet[b.codeSize-1])
	}
}

func TestConfirmPeg(t *testing.T) {
	b, err := newBoard(4, 4, 1)
	errF(err, t)
	b.code.pegs = []peg{b.pegSet[3], b.pegSet[2], b.pegSet[1], b.pegSet[0]} //to ensure that a win doesn't prematurely end the test case
	for i := 0; i < b.codeSize; i++ {
		for j := 0; j < i; j++ {
			b.selectPeg(true)
		}
		b.confirmPeg()
	}

	if r := b.currRow; r != 1 {
		t.Fatalf("expected: 1, got: %v", r)
	}

	r0 := b.rows[0].pegs

	for i := 0; i < b.codeSize; i++ {
		if r0[i] != b.pegSet[i] {
			t.Fatalf("expected: %v, got: %v", r0[i], b.pegSet[i])
		}
	}
}

func TestDeletePeg(t *testing.T) {
	b, err := newBoard(4, 4, 1)
	errF(err, t)
	b.confirmPeg()
	b.confirmPeg()
	b.confirmPeg()
	b.deletePeg()
	if i := b.row().currPeg; i != 2 {
		t.Fatalf("expected: 1, got: %v", i)
	}
}

func TestWin(t *testing.T) {
	b, err := newBoard(4, 4, 1)
	errF(err, t)
	b.code.pegs = []peg{b.pegSet[0], b.pegSet[1], b.pegSet[2], b.pegSet[3]}
	for i := 0; i < b.codeSize; i++ {
		b.confirmPeg()
	}
	if b.won() {
		t.Fatalf("expected: no win, got: win")
	}
	for i := 0; i < b.codeSize; i++ {
		for j := 0; j < i; j++ {
			b.selectPeg(true)
		}
		b.confirmPeg()
	}
	if !b.won() {
		t.Fatalf("expected: win, got: none")
	}
}

func TestLoss(t *testing.T) {
	b, err := newBoard(4, 4, 1)
	errF(err, t)
	b.code.pegs = []peg{b.pegSet[0], b.pegSet[1], b.pegSet[2], b.pegSet[3]}
	if b.lost {
		t.Fatalf("expected: no loss, got: loss")
	}
	for i := 0; i < b.codeSize*maxTries; i++ {
		b.confirmPeg()
	}
	if !b.lost {
		t.Fatalf("expected: loss, got: none")
	}
}
