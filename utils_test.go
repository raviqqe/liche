package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringSetToSlice(t *testing.T) {
	assert.Equal(t, []string{"foo", "bar"}, stringSetToSlice(map[string]bool{"foo": true, "bar": false}))
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
