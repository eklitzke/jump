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

// Database represents the database.
type Database interface {
	AdjustWeight(path string, weight float64)
	Remove(path string)
	Weights() WeightMap // FIXME: this is hacky and overly tied to the gob implementation
	Replace(WeightMap)
	Prune(int)
	Save() error
	Search(needle string) Entry
}

// NewDatabase loads a database file.
func NewDatabase(path string, opts Options) Database {
	return NewGobDatabase(path, opts)
}
