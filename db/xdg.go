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
	"os"
	"path/filepath"
)

const (
	vendorName = "jump"       // the xdg application name
	dbName     = "db.gob"     // name of the database file
	configName = "config.yml" // the config file name
)

func dirOrTmp(dir string, err error) string {
	if err != nil {
		return "/tmp"
	}
	return dir
}

// DatabasePath looks up the full path to the database file.
func DatabasePath() string {
	return filepath.Join(dirOrTmp(os.UserCacheDir()), vendorName, dbName)
}

// ConfigPath returns the path to the jump config file.
func ConfigPath() string {
	return filepath.Join(dirOrTmp(os.UserConfigDir()), vendorName, configName)
}
