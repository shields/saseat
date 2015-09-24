// Copyright 2015 Michael Shields
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package saseat

import (
	"math"
	"math/rand"
)

const (
	// The ideal table has two possible seats left unset.  Apply
	// an exponential penalty for having more or fewer seats.
	CapacityPref = -50

	// The ideal table has the same number of seats on each side.
	// Apply an exponential penalty if that doesn't happen.
	ImbalancePref = -400

	// Apply a penalty if guests at adjacent seats both have
	// genders and they are the same.
	SameGenderPref = -10

	// Weight how much to consider each neighbor.
	AdjacentWeight  = 1
	OppositeWeight  = 0.6
	DiagonalWeight  = 0.3
	SameTableWeight = 0.1
)

type Table struct {
	Left, Right []Guest
	// Each guest's total preference, including gender penalty.
	LeftScores, RightScores []float64
	// The table's score, not including guest scores.
	TableScore float64
	// The sum of the table scores and the guest scores.
	Score float64
}

// NewTable returns a new Table with capacity seats.  capacity must be
// a positive even number.
func NewTable(capacity int) Table {
	return Table{
		Left:        make([]Guest, capacity/2),
		Right:       make([]Guest, capacity/2),
		LeftScores:  make([]float64, capacity/2),
		RightScores: make([]float64, capacity/2),
	}
}

// Rescore recalculates the scores based on the current guests.
func (t *Table) Rescore(p Prefs) {
	var left, right int
	for _, g := range t.Left {
		if g.Name != "" {
			left++
		}
	}
	for _, g := range t.Right {
		if g.Name != "" {
			right++
		}
	}

	t.TableScore = 0
	t.TableScore += CapacityPref * math.Expm1(math.Abs(float64(
		(len(t.Left)+len(t.Right)-2)-(left+right))))
	t.TableScore += ImbalancePref * math.Expm1(math.Abs(float64(left-right)))

	t.scoreSide(t.Left, t.Right, t.LeftScores, p)
	t.scoreSide(t.Right, t.Left, t.RightScores, p)

	t.Score = t.TableScore
	for _, s := range t.LeftScores {
		if s != 0 {
			t.Score += s
		}
	}
	for _, s := range t.RightScores {
		if s != 0 {
			t.Score += s
		}
	}
}

func (t *Table) scoreSide(this, other []Guest, scores []float64, p Prefs) {
	for i, g := range this {
		var s float64
		if i > 0 {
			s += p.Pref(g.Name, this[i-1].Name) * AdjacentWeight
			s += p.Pref(g.Name, other[i-1].Name) * DiagonalWeight
			if g.Gender != "" && this[i-1].Gender != "" && g.Gender == this[i-1].Gender {
				s += SameGenderPref
			}
		}
		s += p.Pref(g.Name, other[i].Name) * OppositeWeight
		if i < len(this)-1 {
			s += p.Pref(g.Name, this[i+1].Name) * AdjacentWeight
			s += p.Pref(g.Name, other[i+1].Name) * DiagonalWeight
			if g.Gender != "" && this[i+1].Gender != "" && g.Gender == this[i+1].Gender {
				s += SameGenderPref
			}
		}
		// Check same-table scores for guests who are more than one
		// space away.
		for j := 0; j < i-1; j++ {
			s += p.Pref(g.Name, this[j].Name) * SameTableWeight
			s += p.Pref(g.Name, other[j].Name) * SameTableWeight
		}
		for j := i + 2; j < len(this); j++ {
			s += p.Pref(g.Name, this[j].Name) * SameTableWeight
			s += p.Pref(g.Name, other[j].Name) * SameTableWeight
		}
		scores[i] = s
	}
}

// SwapSide randomly swaps some contiguous portion of the guests between this
// table and another table.  If t2 is this table, then guests will be swapped
// between the left and right sides.  If max is -1, up to an entire side may be
// swapped.
func (t *Table) Swap(t2 *Table, max int) {
	var side1, side2 []Guest
	if t == t2 {
		side1 = t.Left
		side2 = t2.Right
	} else {
		switch rand.Intn(4) {
		case 0:
			side1 = t.Left
			side2 = t2.Left
		case 1:
			side1 = t.Left
			side2 = t2.Right
		case 2:
			side1 = t.Right
			side2 = t2.Left
		case 3:
			side1 = t.Right
			side2 = t2.Right
		}
	}

	// side1 is the shorter side.
	if len(side1) > len(side2) {
		side1, side2 = side2, side1
	}

	// How many can we swap?  At least one, at most the
	// number of seats on the smaller side.
	n := rand.Intn(len(side1)-1) + 1
	if max != -1 && n > max {
		n = max
	}
	offset := rand.Intn(len(side1) - n)

	// If side2 is longer, maybe index further into it.
	offset2 := offset
	if len(side1) != len(side2) {
		offset2 += rand.Intn(len(side2) - len(side1))
	}

	// Swap.
	for i := 0; i < n; i++ {
		side1[offset+i], side2[offset2+i] = side2[offset2+i], side1[offset+i]
	}
	closeGaps(side1)
	closeGaps(side2)
}

// closeGaps adjusts side such that all the guests have been moved to the
// beginning and all the empty seats to the end.
func closeGaps(side []Guest) {
	var gaps int
	for i := 0; i < len(side); i++ {
		if gaps > 0 {
			side[i-gaps] = side[i]
		}
		if side[i].Name == "" {
			gaps++
		}
	}
	for i := len(side) - gaps; i < len(side); i++ {
		side[i] = Guest{}
	}
}
