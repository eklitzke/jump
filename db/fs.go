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

	"github.com/rs/zerolog/log"
)

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
