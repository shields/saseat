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

// XXX
// XXX INCOMPLETE HACKED VERSION FOR ROUND TABLES, WITH NO TESTS
// XXX

package saseat

import (
	"math"
	"math/rand"
)

const (
	// Apply a penalty if guests at adjacent seats both have
	// genders and they are the same.
	SameGenderPref = -10
)

type Table struct {
	Guests []Guest
	// Each guest's total preference, including gender penalty.
	Scores []float64
	// The sum of the guest scores.
	Score float64
}

// NewTable returns a new Table with capacity seats.  capacity must be
// a positive number.
func NewTable(capacity int) Table {
	return Table{
		Guests: make([]Guest, capacity),
		Scores: make([]float64, capacity),
	}
}

// Rescore recalculates the scores based on the current guests.
func (t *Table) Rescore(p Prefs) {
	segment := 2 * math.Pi / float64(len(t.Guests))

	t.Score = 0
	t.Scores = make([]float64, len(t.Scores))

	for i, g1 := range t.Guests {
		for j, g2 := range t.Guests {
			if i >= j {
				continue
			}

			x1 := math.Sin(segment * float64(i))
			x2 := math.Sin(segment * float64(j))
			y1 := math.Cos(segment * float64(i))
			y2 := math.Cos(segment * float64(j))

			// Magically apply exponentially weighted distance: squared
			// Euclidean distance is just not taking the root of the sum of
			// the squared distance.
			d := (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
			// d can be up to 4, which is one diameter (2) squared.
			weight := (4 - d) / 4
			if weight < 0.1 {
				weight = 0.1
			}
			s := p.Pref(g1.Name, g2.Name) * weight
			t.Scores[i] += s
			t.Scores[j] += s
			t.Score += s

			if (i+1 == j || (i == 0 && j == len(t.Guests))) &&
				g1.Gender != "" && g2.Gender != "" && g1.Gender == g2.Gender {
				t.Score += SameGenderPref
			}
		}
	}
}

// Swap randomly swaps some contiguous portion of the guests between this
// table and another table.  If t2 is this table, then guests will be swapped
// between the left and right sides.  If max is -1, up to an entire side may be
// swapped.
func (t *Table) Swap(t2 *Table, max int) {
	var side1, side2 []Guest
	if t == t2 {
		side1 = t.Guests[0 : len(t.Guests)/2]
		side2 = t.Guests[len(t.Guests)/2 : len(t.Guests)]
	} else {
		side1 = t.Guests
		side2 = t2.Guests
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
}
