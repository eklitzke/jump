// Copyright 2019 Evan Klitzke <evan@eklitzke.org>
//
// This file is part of jump.
//
// jump is free software: you can redistribute it and/or modify it under
// the terms of the GNU General Public License as published by the Free Software
// Foundation, either version 3 of the License, or (at your option) any later
// version.
//
// jump is distributed in the hope that it will be useful, but WITHOUT ANY
// WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
// A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with
// jump. If not, see <http://www.gnu.org/licenses/>.

package db

import (
	"time"
)

// Weight represents a weight value with a timestamp.
type Weight struct {
	Value     float64
	UpdatedAt time.Time
}

// NewWeight creates a new weight value with the current timestamp.
func NewWeight(val float64) Weight {
	return Weight{
		Value:     val,
		UpdatedAt: time.Now().UTC(),
	}
}

// WeightMap is a map from string paths to weights.
type WeightMap map[string]Weight
