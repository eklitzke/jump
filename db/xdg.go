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
	"os"
	"path/filepath"

	"github.com/OpenPeeDeeP/xdg"
	"github.com/rs/zerolog/log"
)

const (
	xdgName = "jump"   // the xdg application name
	dbName  = "db.gob" // name of the database file
)

// full path to the database file
var dbPath string

// DatabasePath looks up the full path to the database file.
func DatabasePath() string {
	if dbPath == "" {
		x := xdg.New(xdgName, "")
		dataHome := x.DataHome()
		_, err := os.Stat(dataHome)
		if err != nil {
			if os.IsNotExist(err) {
				// create the data directory
				if err := os.Mkdir(dataHome, 0700); err != nil {
					log.Fatal().Err(err).Str("dir", dataHome).Msg("failed to create data directory")
				}
			} else {
				log.Fatal().Err(err).Str("dir", dataHome).Msg("failed to stat data directory")
			}
		}
		dbPath = filepath.Join(dataHome, dbName)
	}
	return dbPath
}
