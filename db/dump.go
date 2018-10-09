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
	"os/user"
	"sort"
	"strings"

	"github.com/rs/zerolog/log"
)

// DumpOpts represents the dump options
type DumpOpts struct {
	Short bool // shorten output paths by abbreviating tilde
}

// shorten a path
func shortenPath(me *user.User, path string) string {
	if strings.HasPrefix(path, me.HomeDir) {
		return "~" + path[len(me.HomeDir):]
	}
	return path
}

// Dump returns a JSON serializable representation of the database weights.
func Dump(d Database, opts DumpOpts) interface{} {
	weights := d.GetWeights()
	if opts.Short {
		if me, err := user.Current(); err != nil {
			log.Error().Err(err).Msg("failed to lookup user")
		} else {
			var newWeights []Entry
			for _, w := range weights {
				w.Path = shortenPath(me, w.Path)
				newWeights = append(newWeights, w)
			}
			weights = newWeights
		}
	}

	output := struct {
		Format  string  `json:"format"`
		Weights []Entry `json:"weights"`
	}{
		Format:  "gob",
		Weights: weights,
	}
	sort.Sort(descendingWeight(output.Weights))
	return output

}
