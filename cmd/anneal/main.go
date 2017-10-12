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

package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"msrl.com/hacks/saseat"
)

var (
	guestFlag = flag.String("guests", "", "CSV file of guest names")
	prefFlag  = flag.String("prefs", "", "CSV file of preferences")
	tableFlag = flag.String("tables", "", "comma-separated list of table sizes")
)

func main() {
	rand.Seed(time.Now().UnixNano())

	flag.Parse()
	if *guestFlag == "" || *prefFlag == "" || *tableFlag == "" {
		fmt.Fprintln(os.Stderr, "all flags are required")
		os.Exit(2)
	}

	// Load guests file.
	f, err := os.Open(*guestFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	guests, err := saseat.ReadGuests(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Load prefs file.
	f, err = os.Open(*prefFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	prefs, err := saseat.ReadPrefs(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err = prefs.CheckGuests(guests); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Create tables, packing the guests to them in random order.
	var tables []saseat.Table
	perm := rand.Perm(len(guests))
	seated := 0
	for _, s := range strings.Split(*tableFlag, ",") {
		capacity, err := strconv.ParseInt(s, 0, 10)
		if err != nil {
			fmt.Fprintln(os.Stderr, "bad table capacity %q: %v", s, err)
			os.Exit(1)
		}
		if capacity <= 0 || capacity%2 != 0 {
			fmt.Fprintln(os.Stderr, "only positive even table sizes supported")
			os.Exit(1)
		}
		t := saseat.NewTable(int(capacity))
		for i := 0; i < int(capacity/2) && seated < len(guests); i++ {
			t.Left[i] = guests[perm[seated]]
			seated++
		}
		for i := 0; i < int(capacity/2) && seated < len(guests); i++ {
			t.Right[i] = guests[perm[seated]]
			seated++
		}
		t.Rescore(prefs)
		tables = append(tables, t)
	}
	if seated != len(guests) {
		fmt.Fprintf(os.Stderr, "seats for only %v of %v guests", seated, len(guests))
		os.Exit(1)
	}

	// Let's anneal!
	temp := 250.0
	reported := time.Now()
	iter := 1
	for ; ; iter++ {
		// Cooling.
		if iter%1000000 == 0 {
			temp *= 0.99
		}

		if time.Now().Sub(reported) > 1*time.Second {
			report(tables)
			fmt.Printf("Iteration %d, temperature %.1f\n\n\n", iter, temp)
			reported = time.Now()
		}

		// Spend more time trying to optimize within tables instead of
		// swapping people around the room.
		if rand.Float64() > 0.1 {
			i := rand.Intn(len(tables))
			t := copyTable(&tables[i])
			t.Swap(&t, 2)
			t.Rescore(prefs)
			if accept(t.Score, tables[i].Score, temp) {
				tables[i] = t
			}
		} else {
			i1 := rand.Intn(len(tables))
			i2 := rand.Intn(len(tables))
			if i1 == i2 {
				continue
			}
			t1 := copyTable(&tables[i1])
			t2 := copyTable(&tables[i2])
			t1.Swap(&t2, -1)
			t1.Rescore(prefs)
			t2.Rescore(prefs)
			if accept(t1.Score+t2.Score, tables[i1].Score+tables[i2].Score, temp) {
				tables[i1] = t1
				tables[i2] = t2
			}
		}
	}

	report(tables)
}

// Makes a partial copy of a table.
func copyTable(t *saseat.Table) saseat.Table {
	tt := saseat.Table{
		Left:        make([]saseat.Guest, len(t.Left)),
		Right:       make([]saseat.Guest, len(t.Right)),
		LeftScores:  make([]float64, len(t.LeftScores)),
		RightScores: make([]float64, len(t.RightScores)),
	}
	copy(tt.Left, t.Left)
	copy(tt.Right, t.Right)
	return tt
}

// Simulated annealing acceptance function.
func accept(new, old float64, temp float64) bool {
	if new >= old {
		return true
	}
	return math.Exp((new-old)/temp) > rand.Float64()
}

func report(tables []saseat.Table) {
	fmt.Printf("----------------------------------------\n\n")
	var Σ float64
	for _, t := range tables {
		Σ += t.Score
	}
	for i, t := range tables {
		printTable(i+1, t)
	}
	fmt.Printf("Score %.0f\n\n", Σ)
}

func printTable(n int, t saseat.Table) {
	fmt.Printf("Table %d -- subtotal %.0f; table %.0f\n\n", n, t.Score, t.TableScore)
	for i := 0; i < len(t.Left); i++ {
		fmt.Printf("%5.0f %-30s  %5.0f %-30s\n",
			t.LeftScores[i], t.Left[i].Name, t.RightScores[i], t.Right[i].Name)
	}
	fmt.Printf("\n\n")
}
