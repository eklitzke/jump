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
	"math"
	"sort"
	"time"

	"github.com/rs/zerolog/log"
)

// StringCompare is a string comparison function.
type StringCompare func(string, string) bool

// Searcher implements the matching algorithm.
type Searcher struct {
	input  weightMap // read-only input weights
	output weightMap // output weights
	opts   Options   // options
}

// Search searches for the needle in the input list using the given comparator,
// and modulates the weight by alpha.
func (s *Searcher) Search(needle string, cmp StringCompare, alpha float64) {
	log.Debug().Str("needle", needle).Float64("alpha", alpha).Msg("doing search")
	for path, inputWeight := range s.input {
		if cmp(path, needle) {
			if w, ok := s.output[path]; ok {
				w.Value *= alpha
				s.output[path] = w
			} else {
				beta := alpha
				if s.opts.TimeMatching {
					elapsed := time.Since(inputWeight.UpdatedAt).Seconds()
					if elapsed > 0 {
						beta /= math.Log1p(elapsed)
					}
				}
				log.Debug().Float64("alpha", alpha).Float64("beta", beta).Str("path", path).Float64("initial_weight", inputWeight.Value).Msg("new search candidate")
				inputWeight.Value *= beta
				s.output[path] = inputWeight
			}
		}
	}
}

// Best returns the best matching entry that is actually a directory.
func (s *Searcher) Best(count int) ([]Entry, []string) {
	var errorPaths []string
	entries := toEntryList(s.output)
	sort.Sort(descendingWeight(entries))
	if s.opts.Debug {
		for rank, entry := range entries {
			log.Debug().Int("rank", rank).Float64("score", entry.Weight).Str("path", entry.Path).Msg("final search candidate")
		}
	}

	var results []Entry
	for _, entry := range entries {
		if err := CheckIsDir(entry.Path); err != nil {
			errorPaths = append(errorPaths, entry.Path)
			continue
		}
		results = append(results, entry)
		if len(results) >= count {
			break
		}
	}
	return results, errorPaths
}

// NewSearcher creates a new searcher instance.
func NewSearcher(input weightMap, opts Options) *Searcher {
	return &Searcher{
		input:  input,
		output: make(weightMap),
		opts:   opts,
	}
}
