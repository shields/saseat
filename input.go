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
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type Guest struct {
	Name   string
	Gender string
}

func ReadGuests(in io.Reader) ([]Guest, error) {
	guests := []Guest{}
	names := make(map[string]bool)

	r := csv.NewReader(in)
	r.FieldsPerRecord = -1

	for {
		record, err := r.Read()
		if err == io.EOF {
			return guests, nil
		}
		if err != nil {
			return nil, err
		}
		var g Guest
		switch len(record) {
		case 2:
			g.Gender = record[1]
			fallthrough
		case 1:
			g.Name = record[0]
			if g.Gender == "" {
				if g.Name[:4] == "Mr. " {
					g.Gender = "male"
				}
				if g.Name[:4] == "Ms. " || g.Name[:5] == "Mrs. " || g.Name[:5] == "Miss" {
					g.Gender = "female"
				}
			}
			if names[g.Name] {
				return nil, fmt.Errorf("duplicate name %q", g.Name)
			}
			names[g.Name] = true
		case 0:
			return nil, errors.New("empty line")
		default:
			return nil, errors.New("too many records")
		}

		guests = append(guests, g)
	}
}

type prefKey struct {
	a, b string // a < b
}

type Prefs map[prefKey]float64

func ReadPrefs(in io.Reader) (Prefs, error) {
	p := make(Prefs)

	r := csv.NewReader(in)
	r.FieldsPerRecord = 3

	for {
		record, err := r.Read()
		if err == io.EOF {
			return p, nil
		}
		if err != nil {
			return nil, err
		}

		a, b := record[0], record[1]
		if a == "" || b == "" {
			return nil, errors.New("empty name")
		}
		if a == b {
			return nil, errors.New("cannot prefer self")
		}
		if a > b {
			a, b = b, a
		}

		c, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}

		key := prefKey{a, b}
		if _, ok := p[key]; ok {
			return nil, errors.New("duplicate pref")
		}
		p[key] = c
	}
}

func (p Prefs) Set(a, b string, pref float64) {
	if a > b {
		a, b = b, a
	}
	p[prefKey{a, b}] = pref
}

func (p Prefs) Pref(a, b string) float64 {
	if a > b {
		a, b = b, a
	}
	return p[prefKey{a, b}]
}

// CheckGuests verifies that each preference is for a named guest.
func (p Prefs) CheckGuests(guests []Guest) error {
	names := make(map[string]bool)
	for _, g := range guests {
		names[g.Name] = true
	}
	for key, _ := range p {
		if !names[key.a] {
			return fmt.Errorf("pref name %q not found in guest list", key.a)
		}
		if !names[key.b] {
			return fmt.Errorf("pref name %q not found in guest list", key.b)
		}
	}
	return nil
}
