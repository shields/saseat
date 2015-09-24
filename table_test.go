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
	"math"
	"reflect"
	"testing"

	"msrl.com/hacks/saseat"
)

func TestScoring(t *testing.T) {
	table := saseat.NewTable(8)
	table.Left[0] = saseat.Guest{"Abe", "male"}
	table.Left[1] = saseat.Guest{"Mary", "female"}
	table.Left[2] = saseat.Guest{"Jane", "female"}
	table.Right[0] = saseat.Guest{"Alexander", "male"}
	table.Right[1] = saseat.Guest{"Kang", ""}
	expectedTableScore := saseat.CapacityPref*math.Expm1(1) + saseat.ImbalancePref*math.Expm1(1)

	prefs := saseat.Prefs{}
	prefs.Set("Abe", "Mary", 1000)
	prefs.Set("Kang", "Jane", -80)
	prefs.Set("Alexander", "Jane", 50)
	expectedLeft := []float64{
		1000 * saseat.AdjacentWeight,
		1000*saseat.AdjacentWeight + saseat.SameGenderPref,
		-80*saseat.DiagonalWeight + 50*saseat.SameTableWeight + saseat.SameGenderPref,
		0,
	}
	expectedRight := []float64{
		50 * saseat.SameTableWeight,
		-80 * saseat.DiagonalWeight,
		0,
		0,
	}

	table.Rescore(prefs)
	if table.TableScore != expectedTableScore {
		t.Errorf("got table score %v, want %v", table.TableScore, expectedTableScore)
	}
	if !reflect.DeepEqual(table.LeftScores, expectedLeft) {
		t.Errorf("got left prefs %v, want %v", table.LeftScores, expectedLeft)
	}
	if !reflect.DeepEqual(table.RightScores, expectedRight) {
		t.Errorf("got right prefs %v, want %v", table.RightScores, expectedRight)
	}
}
