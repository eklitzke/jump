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

import (
	"io"
)

// Database represents the database.
type Database interface {
	// Adjust the weight for a given path; weight can be positive or
	// negative.
	AdjustWeight(string, float64)

	// Does the database need to be saved?
	// TODO: this is a little hacky.
	Dirty() bool

	// Return the list of weights in the database.
	GetWeights() []Entry

	// Remove a path from the database.
	Remove(string)

	// Replace the current weights.
	Replace([]Entry)

	// Prune the database.
	Prune(int)

	// Save the database to a writer.
	Save(io.Writer) error

	// Search for a query and find the best match.
	// TODO: allow this to return multiple results.
	Search(int, ...string) []Entry
}

// NewDatabase loads a database file.
func NewDatabase(r io.Reader, opts Options) Database {
	return NewGobDatabase(r, opts)
}
