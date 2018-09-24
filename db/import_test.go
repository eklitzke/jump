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
	"strings"

	"github.com/eklitzke/jump/db"
	. "gopkg.in/check.v1"
)

func (s *MySuite) TestImport(c *C) {
	r := strings.NewReader("1.0 foo\n2.0 bar\nbaz\nx y\n")
	weights, err := db.LoadAutojumpDatabase(r)
	c.Assert(err, IsNil)
	c.Assert(weights, HasLen, 2)
}
