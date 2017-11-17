package main

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestExtractUrls(t *testing.T) {
	for _, c := range []struct {
		html    string
		numUrls int
	}{
		{`<a href="https://google.com">Google</a>`, 1},
		{
			`
				<div>
					<a href="https://google.com">Google</a>
					<a href="https://google.com">Google</a>
				</div>
			`,
			1,
		},
		{
			`
				<div>
					<a href="https://google.com">Google</a>
					<a href="https://yahoo.com">Yahoo!</a>
				</div>
			`,
			2,
		},
	} {
		n, err := html.Parse(strings.NewReader(c.html))

		assert.Equal(t, nil, err)
		assert.Equal(t, c.numUrls, len(extractUrls(n)))
	}
}

func TestUrlParse(t *testing.T) {
	u, err := url.Parse("file-path")

	assert.Equal(t, nil, err)
	assert.Equal(t, "", u.Scheme)
}

func TestIsUrl(t *testing.T) {
	for _, s := range []string{"http://google.com", "https://google.com"} {
		assert.True(t, isUrl(s))
	}

	for _, s := range []string{"", "file-path"} {
		assert.False(t, isUrl(s))
	}
}
