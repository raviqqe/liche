package main

import (
	"path/filepath"
	"strings"
)

var extensions = map[string]bool{
	".htm":  true,
	".html": true,
	".md":   true,
}

func isMarkupFile(f string) bool {
	_, ok := extensions[filepath.Ext(f)]
	return ok
}

func isHTMLFile(f string) bool {
	return strings.HasSuffix(f, ".html") || strings.HasSuffix(f, ".htm")
}
