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
	"errors"
	"os"
	"sort"

	"github.com/rs/zerolog/log"
)

// StringCompare is a string comparison function.
type StringCompare func(string, string) bool

// Searcher implements the matching algorithm.
type Searcher struct {
	input  WeightMap // read-only input weights
	output WeightMap // output weights
}

// Search searches for the needle in the input list using the given comparator,
// and modulates the weight by alpha.
func (s *Searcher) Search(needle string, cmp StringCompare, alpha float64) {
	for path, inputWeight := range s.input {
		if cmp(path, needle) {
			if w, ok := s.output[path]; ok {
				w.Value *= alpha
				s.output[path] = w
			} else {
				inputWeight.Value *= alpha
				s.output[path] = inputWeight
			}
		}
	}
}

// Best returns the best matching entry that is actually a directory.
func (s *Searcher) Best() (Entry, []string) {
	var errorPaths []string
	entries := toEntryList(s.output)
	sort.Sort(descendingWeight(entries))
	for _, entry := range entries {
		if err := CheckIsDir(entry.Path); err != nil {
			errorPaths = append(errorPaths, entry.Path)
			continue
		}
		return entry, errorPaths
	}
	log.Debug().Msg("no entries found for query")
	return Entry{}, errorPaths
}

// NewSearcher creates a new searcher instance.
func NewSearcher(input WeightMap) *Searcher {
	return &Searcher{
		input:  input,
		output: make(WeightMap),
	}
}

// ErrNotDir is returned by CheckIsDir when the path is not a directory.
var ErrNotDir = errors.New("path is not a directory")

// CheckIsDir checks that the input path is a directory.
func CheckIsDir(path string) error {
	st, err := os.Stat(path)
	if err != nil {
		log.Warn().Err(err).Str("path", path).Msg("failed to stat path")
		return err
	}
	if !st.IsDir() {
		log.Warn().Str("path", path).Msg("specified file is not a directory")
		return ErrNotDir
	}
	return nil
}
