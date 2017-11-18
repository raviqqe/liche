package main

import "strings"

type fileResult struct {
	filename   string
	urlResults []urlResult
	err        error
}

func (r fileResult) String(verbose bool) string {
	ss := make([]string, 0, len(r.urlResults))

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
