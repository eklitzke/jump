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
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
)

// Database represents the database.
type Database struct {
	path    string             // path to the underlying file
	Weights map[string]float64 // map of entry to weight
}

// AdjustWeight adjusts the weight of a path.
func (d *Database) AdjustWeight(path string, weight float64) {
	if weight >= 0 {
		// increase the weight
		current := d.Weights[path]
		d.Weights[path] = math.Sqrt(current*current + weight*weight)
		return
	}

	// decrease the weight
	newWeight := d.Weights[path] + weight
	if newWeight <= 0 {
		// if the weight is negative or zero, delete it
		d.Remove(path)
		return
	}
	d.Weights[path] = newWeight
}

// Remove removes a path from the database.
func (d *Database) Remove(path string) {
	delete(d.Weights, path)
}

// Dump prints the database to the specified writer.
func (d *Database) Dump(w io.Writer) error {
	entries := []Entry{}
	for path, weight := range d.Weights {
		entries = append(entries, Entry{Path: path, Weight: weight})
	}
	sort.Sort(byWeight(entries))
	for _, entry := range entries {
		if _, err := fmt.Fprintf(w, "%.9f %s\n", entry.Weight, entry.Path); err != nil {
			return err
		}
	}
	return nil
}

// Prune removes entries from the database that no longer exist.
func (d *Database) Prune() {
	var removePaths []string
	for path := range d.Weights {
		st, err := os.Stat(path)
		if err != nil {
			log.Debug().Err(err).Msg("failed to stat file")
			removePaths = append(removePaths, path)
		}
		if !st.IsDir() {
			log.Debug().Msg("removing non-directory entry")
			removePaths = append(removePaths, path)
		}
	}
	for _, path := range removePaths {
		delete(d.Weights, path)
	}
}

// SumWeights computes the sum of weights in the database.
func (d *Database) SumWeights() float64 {
	var sum float64
	for _, weight := range d.Weights {
		sum += weight
	}
	return sum
}

// Save atomically saves the database.
func (d *Database) Save() error {
	dir := filepath.Dir(d.path)
	temp, err := ioutil.TempFile(dir, fmt.Sprintf(".%s-", dbName))
	if err != nil {
		log.Error().Err(err).Str("dir", dir).Msg("failed to create temporary save file")
		return err
	}

	// clean up the temporary file when we're done with it
	tempName := temp.Name()
	defer func() {
		if err := temp.Close(); err != nil {
			log.Error().Err(err).Str("path", tempName).Msg("error closing temporary file")
		}
		if err := os.Remove(tempName); err != nil && !os.IsNotExist(err) {
			log.Error().Err(err).Str("path", tempName).Msg("failed to close temporary file")
		}
	}()

	w := bufio.NewWriter(temp)
	enc := gob.NewEncoder(w)
	if err := enc.Encode(d); err != nil {
		log.Error().Err(err).Msg("failed to gob encode database")
		return err
	}
	if err := w.Flush(); err != nil {
		log.Error().Err(err).Msg("failed to flush temporary file")
		return err
	}
	return os.Rename(tempName, d.path)
}

var errNotFound = errors.New("entry not found")

type stringCompare func(string, string) bool

func (d *Database) search(cmp stringCompare, needle string) (Entry, error) {
	var candidates []Entry
	for path, weight := range d.Weights {
		if cmp(path, needle) {
			candidates = append(candidates, Entry{Path: path, Weight: weight})
		}
	}
	if candidates == nil {
		return Entry{}, errNotFound
	}
	return FindHighestWeight(candidates), nil
}

// Search searches for the best database entry.
func (d *Database) Search(needle string) Entry {
	// first check exact suffix matches
	if entry, err := d.search(strings.HasSuffix, needle); err == nil {
		return entry
	}

	// next try any contains
	if entry, err := d.search(strings.Contains, needle); err == nil {
		return entry
	}

	// return an empty entry by default
	return Entry{}
}

// NewDatabase loads a database file.
func NewDatabase(path string) *Database {
	db := &Database{path: path, Weights: make(map[string]float64)}
	dbFile, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug().Err(err).Str("path", path).Msg("database file not found")
		} else {
			// a more serious error
			log.Error().Err(err).Str("path", path).Msg("failed to open database file")
		}
		return db
	}
	defer func() {
		if err := dbFile.Close(); err != nil {
			log.Warn().Err(err).Str("path", path).Msg("failed to close db file")
		}
	}()

	dec := gob.NewDecoder(dbFile)
	if err := dec.Decode(db); err != nil {
		log.Error().Err(err).Msg("failed to decode database file")
	}
	return db
}

// defaultDB is a handle to a default database
var defaultDB *Database

// LoadDefaultDatabase loads the default database handle.
func LoadDefaultDatabase() *Database {
	if defaultDB == nil {
		defaultDB = NewDatabase(DatabasePath())
	}
	return defaultDB
}
