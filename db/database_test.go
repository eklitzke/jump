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

package db_test

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/eklitzke/jump/db"
	"github.com/rs/zerolog"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	TestingT(t)
}

type MySuite struct{}

func (s *MySuite) createTempDir(c *C) string {
	baseDir, err := ioutil.TempDir("", "jump-test-")
	c.Assert(err, IsNil)
	c.Assert(baseDir, Not(Equals), "")
	return baseDir
}

func (s *MySuite) TestLoadDatabase(c *C) {
	handle := db.NewDatabase(strings.NewReader(""), db.Options{})
	c.Assert(handle, Not(Equals), nil)
}

var _ = Suite(&MySuite{})
