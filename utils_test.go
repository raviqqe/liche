package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListFiles(t *testing.T) {
	fs, err := listFiles(".")

	assert.Equal(t, nil, err)
	assert.NotEqual(t, 0, len(fs))

	for _, f := range fs {
		i, err := os.Stat(f)

		assert.True(t, isMarkupFile(f))
		assert.Equal(t, nil, err)
		assert.False(t, i.IsDir())
	}
}
