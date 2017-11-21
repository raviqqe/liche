package main

import "path/filepath"

var extensions = map[string]bool{
	".htm":  true,
	".html": true,
	".md":   true,
}

func isMarkupFile(f string) bool {
	_, ok := extensions[filepath.Ext(f)]
	return ok
}
