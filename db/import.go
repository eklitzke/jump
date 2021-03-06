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
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	autojumpVendor = "autojump"
	autojumpDbFile = "autojump.txt"
)

func FindAutojumpDatabase() string {
	// XXX: Or is it config dir? I forget
	return filepath.Join(dirOrTmp(os.UserCacheDir()), autojumpVendor, autojumpDbFile)
}

// LoadAutojumpDatabase loads the autojump database file
func LoadAutojumpDatabase(r io.Reader) ([]Entry, error) {
	var entries []Entry
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		sep := strings.IndexAny(line, " \t")
		if sep == -1 {
			log.Warn().Str("line", line).Msg("failed to split line")
			continue
		}

		stringWeight := line[:sep]
		weight, err := strconv.ParseFloat(stringWeight, 64)
		if err != nil {
			log.Warn().Str("line", line).Msg("failed to parse weight as float64")
			continue
		}
		entries = append(entries, Entry{
			Path:      strings.TrimSpace(line[sep:]),
			Weight:    weight,
			UpdatedAt: time.Now().UTC(),
		})
	}
	if err := scanner.Err(); err != nil {
		log.Error().Err(err).Msg("error scanning file")
		return nil, err
	}
	return entries, nil
}
