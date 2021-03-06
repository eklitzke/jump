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
	"encoding/gob"
	"io"
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
	opts    Options   // database options
	Weights weightMap // map of entry to weight
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

// Dirty checks the dirty bit.
func (d *GobDatabase) Dirty() bool {
	return d.dirty
}

// Remove removes a path from the database.
func (d *GobDatabase) Remove(path string) {
	d.dirty = true
	delete(d.Weights, path)
}

// Prune removes entries from the database that no longer exist.
func (d *GobDatabase) Prune(maxEntries int, excludePatterns []string) {
	// delete non-existent entries
	for path := range d.Weights {
		st, err := os.Stat(path)
		if err != nil {
			log.Debug().Err(err).Str("path", path).Msg("failed to stat file")
			delete(d.Weights, path)
			d.dirty = true
			continue
		}
		if !st.IsDir() {
			log.Debug().Str("path", path).Msg("removing non-directory entry")
			delete(d.Weights, path)
			d.dirty = true
			continue
		}
		for _, pattern := range excludePatterns {
			if strings.Contains(path, pattern) {
				log.Debug().Str("path", path).Str("pattern", pattern).Msg("removing path matching exclude pattern")
				delete(d.Weights, path)
				d.dirty = true
				break
			}
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

// Save atomically saves the database.
func (d *GobDatabase) Save(w io.Writer) error {
	enc := gob.NewEncoder(w)
	if err := enc.Encode(d.Weights); err != nil {
		log.Error().Err(err).Msg("failed to encode gob database")
		return err
	}
	d.dirty = false
	return nil
}

// Search searches for the best database entry.
func (d *GobDatabase) Search(count int, needles ...string) []Entry {
	s := NewSearcher(d.Weights, d.opts)

	// Assume all components form the suffix of the directory name.
	needle := filepath.Join(needles...)

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

	// find the best match
	results, errorPaths := s.Best(count)

	// if any errors were encountered, remove those paths
	for _, path := range errorPaths {
		log.Warn().Str("path", path).Msg("removing bad path")
		d.Remove(path)
	}

	return results
}

// Dump prints the database to the specified writer.
func (d *GobDatabase) GetWeights() []Entry {
	return toEntryList(d.Weights)
}

// Replace replaces the underlying weight map.
func (d *GobDatabase) Replace(entries []Entry) {
	d.Weights = make(weightMap)
	for _, entry := range entries {
		d.Weights[entry.Path] = Weight{Value: entry.Weight, UpdatedAt: entry.UpdatedAt}
	}
	d.dirty = true
}

// NewGobDatabase loads a database file.
func NewGobDatabase(r io.Reader, opts Options) *GobDatabase {
	db := &GobDatabase{
		opts:    opts,
		Weights: make(weightMap),
	}
	dec := gob.NewDecoder(r)
	if err := dec.Decode(&db.Weights); err != nil && err != io.EOF {
		log.Error().Err(err).Msg("failed to decode weights for gob database")
	}
	return db
}
