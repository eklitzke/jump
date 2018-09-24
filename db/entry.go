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

import "time"

// Entry represents a database entry.
type Entry struct {
	Path      string
	Weight    float64
	UpdatedAt time.Time
}

type descendingWeight []Entry

func (b descendingWeight) Len() int           { return len(b) }
func (b descendingWeight) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b descendingWeight) Less(i, j int) bool { return b[i].Weight > b[j].Weight }

type ascendingWeight []Entry

func (b ascendingWeight) Len() int           { return len(b) }
func (b ascendingWeight) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ascendingWeight) Less(i, j int) bool { return b[i].Weight > b[j].Weight }

// FindHighestWeight searches an entry list for the entry with the highest weight.
func FindHighestWeight(entries []Entry) Entry {
	var best Entry
	for _, entry := range entries {
		if entry.Weight > best.Weight {
			best = entry
		}
	}
	return best
}

func toEntryList(w WeightMap) []Entry {
	var entries []Entry
	for path, weight := range w {
		entries = append(entries, Entry{
			Path:      path,
			Weight:    weight.Value,
			UpdatedAt: weight.UpdatedAt,
		})
	}
	return entries
}
