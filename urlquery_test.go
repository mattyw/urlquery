package urlquery_test

import (
	"testing"

	"github.com/mattyw/urlquery"
	gc "gopkg.in/check.v1"
)

func TestAll(t *testing.T) { gc.TestingT(t) }

type urlQuerySuite struct{}

var _ = gc.Suite(&urlQuerySuite{})

func (s *urlQuerySuite) TestMarshalStructTag(c *gc.C) {
	type S struct {
		A int    `urlquery:"foo"`
		B string `urlquery:"str"`
	}
	st := S{
		A: 1234,
		B: "abcdef",
	}
	str, err := urlquery.Marshal(st)
	c.Assert(err, gc.IsNil)
	expectedQuery := "foo=1234&str=abcdef"
	c.Assert(str, gc.Equals, expectedQuery)
}

func (s *urlQuerySuite) TestMarshalStruct(c *gc.C) {
	type S struct {
		A int
		B string
	}
	st := S{
		A: 1234,
		B: "abcdef",
	}
	str, err := urlquery.Marshal(st)
	c.Assert(err, gc.IsNil)
	expectedQuery := "a=1234&b=abcdef"
	c.Assert(str, gc.Equals, expectedQuery)
}

func (s *urlQuerySuite) TestUnmarshalMarshalStruct(c *gc.C) {
	type S struct {
		A int
		B string
	}
	expectedS := S{
		A: 1234,
		B: "abcdef",
	}
	var st S
	query := "a=1234&b=abcdef"
	err := urlquery.Unmarshal(query, &st)
	c.Assert(err, gc.IsNil)
	c.Assert(st, gc.DeepEquals, expectedS)
}

func (s *urlQuerySuite) TestUnmarshalMarshalStructTag(c *gc.C) {
	type S struct {
		A int    `urlquery:"foo"`
		B string `urlquery:"str"`
	}
	expectedS := S{
		A: 1234,
		B: "abcdef",
	}
	var st S
	query := "foo=1234&str=abcdef"
	err := urlquery.Unmarshal(query, &st)
	c.Assert(err, gc.IsNil)
	c.Assert(st, gc.DeepEquals, expectedS)
}
