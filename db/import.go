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
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/rs/zerolog/log"
)

const (
	autojumpVendor = "autojump"
	autojumpDbFile = "autojump.txt"
)

var errNoAutojumpDatabase = errors.New("failed to find autojump database")

// LoadAutojumpDatabase loads the autojump database file
func LoadAutojumpDatabase() (map[string]float64, error) {
	x := xdg.New(autojumpVendor, "")
	dbPath := x.QueryData(autojumpDbFile)
	if dbPath == "" {
		return nil, errNoAutojumpDatabase
	}
	f, err := os.Open(dbPath)
	if err != nil {
		log.Error().Err(err).Str("path", dbPath).Msg("failed to open autojump database")
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Error().Err(err).Str("path", dbPath).Msg("failed to close autojump database")
		}
	}()

	weights := make(map[string]float64)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		sep := strings.IndexAny(line, " \t")
		if sep == -1 {
			log.Warn().Str("path", dbPath).Str("line", line).Msg("failed to split line")
			continue
		}

		stringWeight := line[:sep]
		weight, err := strconv.ParseFloat(stringWeight, 64)
		if err != nil {
			log.Warn().Str("path", dbPath).Str("line", line).Msg("failed to parse weight as float64")
			continue
		}
		path := strings.TrimSpace(line[sep:])
		weights[path] = weight
	}
	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Str("path", dbPath).Msg("error scanning file")
		return nil, err
	}
	return weights, nil
}
