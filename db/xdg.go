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
	"path/filepath"

	goxdg "github.com/OpenPeeDeeP/xdg"
)

const (
	vendorName = "jump"     // the xdg application name
	dbName     = "db.gob"   // name of the database file
	configName = "jump.yml" // the config file name
)

// a handle to an xdg instance
var xdg *goxdg.XDG

func init() {
	xdg = newXDG(vendorName)
}

func newXDG(vendor string) *goxdg.XDG {
	return goxdg.New(vendor, "")
}

// DatabasePath looks up the full path to the database file. The method ensures
// the data directory exists, so writers to the database don't need to handle
// this case.
func DatabasePath() string {
	// We don't use xdg.QueryData here because it checks that the containing
	// directory exists; we defer directory creation until we actually save
	// the database.
	return filepath.Join(xdg.DataHome(), dbName)
}

// ConfigPath returns the path to the jump config file.
func ConfigPath() string {
	return filepath.Join(xdg.ConfigHome(), configName)
}
