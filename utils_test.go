package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListDirectory(t *testing.T) {
	fc := make(chan string, 1024)
	err := listDirectory(".", fc)
	close(fc)

	assert.Equal(t, nil, err)
	assert.NotEqual(t, 0, len(fc))

	for f := range fc {
		i, err := os.Stat(f)

		assert.True(t, isMarkupFile(f))
		assert.Equal(t, nil, err)
		assert.False(t, i.IsDir())
	}
}
