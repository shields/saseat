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

package saseat_test

import (
	"reflect"
	"strings"
	"testing"

	"msrl.com/hacks/saseat"
)

var guestsFile = `Mr. Abraham Lincoln,
Mrs. Mary Lincoln
Dr. Woodrow Wilson
X
`

var expectedGuests = []saseat.Guest{
	{"Mr. Abraham Lincoln", "male"},
	{"Mrs. Mary Lincoln", "female"},
	{"Dr. Woodrow Wilson", ""},
	{"X", ""},
}

func TestReadGuests(t *testing.T) {
	g, err := saseat.ReadGuests(strings.NewReader(guestsFile))
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(g, expectedGuests) {
		t.Errorf("got: %q\nwant: %q\n", g, expectedGuests)
	}
}

var prefsFile = `Abe,Mary,999
Woody,Mary,42.1
`

var expectedPrefs = []struct {
	a, b     string
	expected float64
}{
	{"Abe", "Mary", 999},
	{"Mary", "Abe", 999},
	{"Mary", "Woody", 42.1},
	{"Woody", "Mary", 42.1},
	{"Abe", "Woody", 0},
	{"Jane", "Tarzan", 0},
}

func testPrefs(t *testing.T) {
	p, err := saseat.ReadPrefs(strings.NewReader(prefsFile))
	if err != nil {
		t.Error(err)
	}
	for _, tt := range expectedPrefs {
		actual := p.Pref(tt.a, tt.b)
		if actual != tt.expected {
			t.Errorf("Pref(%q, %q) = %v, want %v", tt.a, tt.b, actual, tt.expected)
		}
	}
}
