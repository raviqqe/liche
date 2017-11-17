package main

import (
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
