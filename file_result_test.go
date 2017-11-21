package main

import (
	"errors"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileResultString(t *testing.T) {
	err := errors.New("error")

	for _, r := range []fileResult{
		{"foo", nil, nil},
		{"foo", []urlResult{{"bar", nil}}, nil},
		{"foo", nil, err},
		{"foo", []urlResult{{"bar", err}}, nil},
		{"foo", []urlResult{{"bar", err}}, err},
		{"foo", []urlResult{{"bar", nil}, {"baz", err}}, nil},
	} {
		b, err := regexp.MatchString(r.filename, r.String(false))

		assert.Equal(t, nil, err)
		assert.Equal(t, true, b)
	}
}

func TestFileResultStringWithVerboseOption(t *testing.T) {
	err := errors.New("error")

	for _, r := range []fileResult{
		{"foo", nil, nil},
		{"foo", []urlResult{{"foo", nil}}, nil},
		{"foo", nil, err},
		{"foo", []urlResult{{"foo", err}}, nil},
		{"foo", []urlResult{{"foo", err}}, err},
		{"foo", []urlResult{{"foo", nil}, {"foo", err}}, nil},
	} {
		s := r.String(true)
		b, err := regexp.MatchString(r.filename, s)

		assert.Equal(t, nil, err)
		assert.True(t, b)
	}
}

func TestFileResultOk(t *testing.T) {
	for _, r := range []fileResult{
		{"foo", nil, nil},
		{"foo", []urlResult{{"foo", nil}}, nil},
	} {
		assert.True(t, r.Ok())
	}

	err := errors.New("error")

	for _, r := range []fileResult{
		{"foo", nil, err},
		{"foo", []urlResult{{"foo", err}}, nil},
		{"foo", []urlResult{{"foo", err}}, err},
		{"foo", []urlResult{{"foo", nil}, {"foo", err}}, nil},
	} {
		assert.False(t, r.Ok())
	}
}
