package db_test

import (
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
	c.Assert(1, Equals, len(handle.Weights))
	c.Assert(w == handle.Weights["foo"], Equals, true)
}
