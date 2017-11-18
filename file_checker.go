package main

import (
	"bytes"
	"io/ioutil"
	"strings"
	"time"

	"golang.org/x/net/html"
	"gopkg.in/russross/blackfriday.v2"
)

type fileChecker struct {
	urlChecker urlChecker
}

func newFileChecker(timeout time.Duration) fileChecker {
	return fileChecker{newURLChecker(timeout)}
}

func (c fileChecker) Check(f string) ([]urlResult, error) {
	bs, err := ioutil.ReadFile(f)

	if err != nil {
		return nil, err
	}

	n, err := html.Parse(bytes.NewReader(blackfriday.Run(bs)))

	if err != nil {
		return nil, err
	}

	us := extractURLs(n)
	rc := make(chan urlResult, len(us))
	rs := make([]urlResult, 0, len(us))

	go c.urlChecker.CheckMany(us, rc)

	for r := range rc {
		rs = append(rs, r)
	}

	return rs, nil
}

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
