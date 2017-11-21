package main

import (
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLResultString(t *testing.T) {
	for _, r := range []urlResult{
		{"https://google.com", nil},
		{"https://yahoo.com", errors.New("error")},
	} {
		p := "OK"

		if r.err != nil {
			p = "ERROR"
		}

		b, err := regexp.MatchString(p, r.String())

		assert.Equal(t, nil, err)
		assert.True(t, b)

		b, err = regexp.MatchString(r.url, r.String())

		assert.Equal(t, nil, err)
		assert.True(t, b)
	}
}
