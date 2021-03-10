package memfs

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/test"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type MemorySuite struct {
	test.FilesystemSuite
	path string
}

var _ = Suite(&MemorySuite{})

func (s *MemorySuite) SetUpTest(c *C) {
	s.FilesystemSuite = test.NewFilesystemSuite(New())
}

func (s *MemorySuite) TestCapabilities(c *C) {
	_, ok := s.FS.(billy.Capable)
	c.Assert(ok, Equals, true)

	caps := billy.Capabilities(s.FS)
	c.Assert(caps, Equals, billy.DefaultCapabilities&^billy.LockCapability)
}

func (s *MemorySuite) TestNegativeOffsets(c *C) {
	f, err := s.FS.Create("negative")
	c.Assert(err, IsNil)

	buf := make([]byte, 100)
	_, err = f.ReadAt(buf, -100)
	c.Assert(err, ErrorMatches, "readat negative: negative offset")

	_, err = f.Seek(-100, io.SeekCurrent)
	c.Assert(err, IsNil)
	_, err = f.Write(buf)
	c.Assert(err, ErrorMatches, "writeat negative: negative offset")
}


func (s *MemorySuite) TestExclusive(c *C) {
	f, err := s.FS.OpenFile("exclusive", os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	c.Assert(err, IsNil)

	fmt.Fprint(f, "mememememe")

	err = f.Close()
	c.Assert(err, IsNil)

	_, err = s.FS.OpenFile("exclusive", os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	c.Assert(err, ErrorMatches, os.ErrExist.Error())

func (s *MemorySuite) TestOrder(c *C) {
	var err error

	files := []string{
		"a",
		"b",
		"c",
	}
	for _, f := range files {
		_, err = s.FS.Create(f)
		c.Assert(err, IsNil)
	}

	attemps := 30
	for n := 0; n < attemps; n++ {
		actual, err := s.FS.ReadDir("")
		c.Assert(err, IsNil)

		for i, f := range files {
			c.Assert(actual[i].Name(), Equals, f)
		}
	}
}
