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
	"encoding/json"
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

// GobDatabase represents the database.
type GobDatabase struct {
	dirty   bool      // dirty bit
	path    string    // path to the underlying file
	opts    Options   // database options
	Weights WeightMap // map of entry to weight
}

// AdjustWeight adjusts the weight of a path. The adjusted weight value is
// returned.
func (d *GobDatabase) AdjustWeight(path string, weight float64) {
	d.dirty = true

	var newWeight float64
	if weight >= 0 {
		// increase the weight
		current := d.Weights[path].Value
		d.Weights[path] = NewWeight(math.Sqrt(current*current + weight*weight))
		return
	}

	// decrease the weight
	newWeight = d.Weights[path].Value + weight
	if newWeight <= 0 {
		// if the weight is negative or zero, delete it
		d.Remove(path)
		return
	}
	d.Weights[path] = NewWeight(newWeight)
}

// Remove removes a path from the database.
func (d *GobDatabase) Remove(path string) {
	d.dirty = true
	delete(d.Weights, path)
}

// Dump prints the database to the specified writer.

// Prune removes entries from the database that no longer exist.
func (d *GobDatabase) Prune(maxEntries int) {
	// delete non-existent entries
	for path := range d.Weights {
		st, err := os.Stat(path)
		if err != nil {
			log.Debug().Err(err).Msg("failed to stat file")
			delete(d.Weights, path)
			d.dirty = true
			continue
		}
		if !st.IsDir() {
			log.Debug().Msg("removing non-directory entry")
			delete(d.Weights, path)
			d.dirty = true
		}
	}

	// delete low weight entries
	// FIXME: this isn't the best heuristic and could be improved
	if maxEntries <= 0 {
		return // ignore zero/negative value
	}
	deleteCount := len(d.Weights) - maxEntries
	if deleteCount > 0 {
		entries := toEntryList(d.Weights)
		sort.Sort(ascendingWeight(entries))
		for i, entry := range entries {
			delete(d.Weights, entry.Path)
			if i == deleteCount-1 {
				break
			}
		}
		d.dirty = true
	}
}

// SumWeights computes the sum of weights in the database.
func (d *GobDatabase) SumWeights() float64 {
	var sum float64
	for _, weight := range d.Weights {
		sum += weight.Value
	}
	return sum
}

// Save atomically saves the database.
func (d *GobDatabase) Save() error {
	if !d.dirty {
		log.Debug().Msg("database not dirty, skipping save")
		return nil
	}

	// ensure the directory exists
	dir := filepath.Dir(d.path)
	ensureDirectory(dir)

	// create the temporary file in the same directory as the destination
	// file, to ensure that the rename operation is atomic
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

	// encode and flush the file
	w := bufio.NewWriter(temp)
	enc := gob.NewEncoder(w)
	if err := enc.Encode(d.Weights); err != nil {
		log.Error().Err(err).Msg("failed to gob encode database")
		return err
	}
	if err := w.Flush(); err != nil {
		log.Error().Err(err).Msg("failed to flush temporary file")
		return err
	}

	// atomic rename
	if err := os.Rename(tempName, d.path); err != nil {
		log.Error().Err(err).Str("dbpath", d.path).Str("tempfile", tempName).Msg("failed to rename db file")
		return err
	}

	d.dirty = false
	return nil
}

// Search searches for the best database entry.
func (d *GobDatabase) Search(needle string) Entry {
	s := NewSearcher(d.Weights, d.opts)

	// first check exact suffix matches
	exact := needle
	if !strings.HasPrefix(exact, "/") {
		exact = "/" + needle
	}
	s.Search(exact, strings.HasSuffix, 10.)

	// next check regular suffix matches
	s.Search(needle, strings.HasSuffix, 2.5)

	// next try any contains matches
	s.Search(needle, strings.Contains, 1.)

	// TODO: implement time relevance as well.

	// find the best match
	best, errorPaths := s.Best()

	// if any errors were encountered, remove those paths
	for _, path := range errorPaths {
		log.Warn().Str("path", path).Msg("removing bad path")
		d.Remove(path)
	}

	return best
}

// Dump prints the database to the specified writer.
func (d *GobDatabase) Dump(w io.Writer) error {
	output := dumpOutput{
		Format:  "gob",
		Weights: toEntryList(d.Weights),
	}
	sort.Sort(descendingWeight(output.Weights))
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}

func (d *GobDatabase) Replace(weights WeightMap) {
	d.Weights = weights
	d.dirty = true
}

// NewGobDatabase loads a database file.
func NewGobDatabase(path string, opts Options) *GobDatabase {
	db := &GobDatabase{
		path:    path,
		opts:    opts,
		Weights: make(WeightMap),
	}
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
	if err := dec.Decode(&db.Weights); err != nil {
		log.Error().Err(err).Str("path", path).Msg("failed to decode database file")
	}
	return db
}

type dumpOutput struct {
	Format  string  `json:"format"`
	Weights []Entry `json:"weights"`
}
