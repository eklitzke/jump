package db_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/eklitzke/jump/db"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestDatabaseEndToEnd(c *C) {
	// TODO: refactor this to not require a temporary file
	f, err := ioutil.TempFile("", "test.db-")
	c.Assert(err, IsNil)
	c.Assert(f.Close(), IsNil)

	defer os.Remove(f.Name())

	handle := db.LoadDatabase(f.Name())
	handle.AdjustWeight("foo", 1)
	w := handle.Weights["foo"]
	c.Assert(w > 0, Equals, true)
	c.Assert(handle.Weights, HasLen, 1)
	c.Assert(handle.Save(), IsNil)

	handle = db.LoadDatabase(f.Name())
	c.Assert(handle.Weights, HasLen, 1)
	c.Assert(w == handle.Weights["foo"], Equals, true)

	entry := handle.Search("nomatch")
	c.Assert(entry, Equals, db.Entry{})
	for _, query := range []string{"f", "foo", "oo"} {
		entry = handle.Search(query)
		c.Assert(entry.Path, Equals, "foo")
	}

	// remove the non-existent directory
	handle.Prune(100)
	c.Assert(handle.Weights, HasLen, 0)

	// add one that does exist
	handle.AdjustWeight("/tmp", 1)
	handle.Prune(100)
	c.Assert(handle.Weights, HasLen, 1)

	// force remove it
	handle.Remove("/tmp")
	c.Assert(handle.Weights, HasLen, 0)
	c.Assert(handle.SumWeights(), Equals, 0.)

	buf := bytes.Buffer{}
	c.Assert(handle.Dump(&buf), IsNil)
	c.Assert(buf.String(), Equals, "")
}

func (s *MySuite) TestConfig(c *C) {
	c.Assert(db.DatabasePath(), Not(Equals), "")
	c.Assert(db.ConfigPath(), Not(Equals), "")
}
