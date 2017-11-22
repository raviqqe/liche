package main

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringSetToSlice(t *testing.T) {
	ss := stringSetToSlice(map[string]bool{"foo": true, "bar": false})

	sort.Strings(ss)

	assert.Equal(t, []string{"bar", "foo"}, ss)
}

func TestIndent(t *testing.T) {
	for _, c := range []struct {
		source, target string
	}{
		{"foo", "\tfoo"},
		{"foo\nbar", "\tfoo\n\tbar"},
	} {
		assert.Equal(t, c.target, indent(c.source))
	}
}

func TestCapitalizeFirst(t *testing.T) {
	for _, ss := range [][2]string{
		{"foo", "Foo"},
		{"foo bar", "Foo bar"},
	} {
		assert.Equal(t, ss[1], capitalizeFirst(ss[0]))
	}
}
