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

package db_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
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

// TODO: split this bad boy up
func (s *MySuite) TestDatabaseEndToEnd(c *C) {
	baseDir := s.createTempDir(c)
	defer os.RemoveAll(baseDir)

	foo := filepath.Join(baseDir, "foo")
	c.Assert(os.MkdirAll(foo, 0755), IsNil)

	handle := db.NewGobDatabase(strings.NewReader(""), db.Options{})

	handle.AdjustWeight(foo, 1)
	w := handle.Weights[foo]
	c.Assert(w.Value > 0, Equals, true)
	c.Assert(handle.Weights, HasLen, 1)
	buf := new(bytes.Buffer)
	c.Assert(handle.Save(buf), IsNil)

	handle = db.NewGobDatabase(buf, db.Options{TimeMatching: true})
	c.Assert(handle.Weights, HasLen, 1)
	c.Assert(w == handle.Weights[foo], Equals, true)

	entry := handle.Search("nomatch")
	c.Assert(entry, Equals, db.Entry{})
	for _, query := range []string{"f", "foo", "oo"} {
		entry = handle.Search(query)
		c.Assert(entry.Path, Equals, foo)
	}

	// remove the non-existent directory
	handle.Prune(100)
	c.Assert(handle.Weights, HasLen, 1)

	buf = new(bytes.Buffer)
	c.Assert(handle.Dump(buf), IsNil)
	c.Assert(buf.String(), Not(Equals), "")

	handle.AdjustWeight(foo, -0.5)
	c.Assert(handle.Weights, HasLen, 1)
	handle.AdjustWeight(foo, -2)
	c.Assert(handle.Weights, HasLen, 0)

	for i := 0; i < 10; i++ {
		dirName := filepath.Join(os.TempDir(), "dbtest", strconv.Itoa(i))
		c.Assert(os.MkdirAll(dirName, 0755), IsNil)
		handle.AdjustWeight(dirName, 1)
	}

	nonDir := filepath.Join(baseDir, "foo.txt")
	g, err := os.Create(nonDir)
	c.Assert(err, IsNil)
	c.Assert(g.Close(), IsNil)
	handle.AdjustWeight(nonDir, 1)

	handle.Prune(3)
	c.Assert(handle.Weights, HasLen, 3)
}

var _ = Suite(&MySuite{})
