package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMarkupFileFinder(t *testing.T) {
	newMarkupFileFinder()
}

func TestMarkupFileFinderFindWithRecursiveOption(t *testing.T) {
	for _, fs := range [][]string{{"README.md"}, {"test"}, {"README.md", "test"}} {
		m := newMarkupFileFinder()
		m.Find(fs, true)

		assert.Equal(t, 0, len(m.Errors()))
		assert.NotEqual(t, 0, len(m.Filenames()))

		for f := range m.Filenames() {
			i, err := os.Stat(f)

			assert.True(t, isMarkupFile(f))
			assert.Equal(t, nil, err)
			assert.False(t, i.IsDir())
		}
	}
}

func TestMarkupFileFinderFindWithDirectory(t *testing.T) {
	m := newMarkupFileFinder()
	m.Find([]string{"test"}, false)

	assert.Equal(t, 1, len(m.Errors()))
	assert.Equal(t, 0, len(m.Filenames()))

	err := <-m.Errors()

	assert.NotEqual(t, nil, err)
}

func TestMarkupFileFinderListDirectory(t *testing.T) {
	m := newMarkupFileFinder()
	m.listDirectory("test")
	close(m.Filenames())

	assert.Equal(t, 0, len(m.Errors()))
	assert.NotEqual(t, 0, len(m.Filenames()))

	for f := range m.Filenames() {
		i, err := os.Stat(f)

		assert.True(t, isMarkupFile(f))
		assert.Equal(t, nil, err)
		assert.False(t, i.IsDir())
	}
}
