package main

import (
	"bytes"
	"io/ioutil"
	"net/url"
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

func extractURLs(n *html.Node) []string {
	ss := make(map[string]bool)
	ns := make([]*html.Node, 0, 1024)
	ns = append(ns, n)

	for len(ns) > 0 {
		i := len(ns) - 1
		n := ns[i]
		ns = ns[:i]

		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				for _, a := range n.Attr {
					if a.Key == "href" && isURL(a.Val) {
						ss[a.Val] = true
						break
					}
				}
			case "img":
				for _, a := range n.Attr {
					if a.Key == "src" && isURL(a.Val) {
						ss[a.Val] = true
						break
					}
				}
			}
		}

		for n := n.FirstChild; n != nil; n = n.NextSibling {
			ns = append(ns, n)
		}
	}

	return stringSetToSlice(ss)
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}
