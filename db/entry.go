// Copyright 2018 Evan Klitzke <evan@eklitzke.org>
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
	Path      string    `json:"path"`
	Weight    float64   `json:"weight"`
	UpdatedAt time.Time `json:"time,string"`
}

type descendingWeight []Entry

func (d descendingWeight) Len() int      { return len(d) }
func (d descendingWeight) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d descendingWeight) Less(i, j int) bool {
	if d[i].Weight == d[j].Weight {
		return d[i].UpdatedAt.After(d[j].UpdatedAt)
	}
	return d[i].Weight > d[j].Weight
}

type ascendingWeight []Entry

func (a ascendingWeight) Len() int      { return len(a) }
func (a ascendingWeight) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ascendingWeight) Less(i, j int) bool {
	if a[i].Weight == a[j].Weight {
		return a[j].UpdatedAt.After(a[i].UpdatedAt)
	}
	return a[j].Weight > a[i].Weight
}

func toEntryList(w weightMap) []Entry {
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
