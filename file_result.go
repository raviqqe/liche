package main

import (
	"strings"

	"github.com/fatih/color"
)

type fileResult struct {
	filename   string
	urlResults []urlResult
	err        error
}

func (r fileResult) String(verbose bool) string {
	ss := make([]string, 0, len(r.urlResults))

	if r.err != nil {
		ss = append(ss, "\t"+color.RedString(r.err.Error()))
	}

	for _, r := range r.urlResults {
		if r.err != nil || verbose {
			ss = append(ss, "\t"+r.String())
		}
	}

	return strings.Join(append([]string{r.filename}, ss...), "\n")
}

func (r fileResult) Ok() bool {
	if r.err != nil {
		return false
	}

	for _, r := range r.urlResults {
		if r.err != nil {
			return false
		}
	}

	return true
}
