package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
	"gopkg.in/russross/blackfriday.v2"
)

type fileChecker struct {
	urlChecker   urlChecker
	documentRoot string
}

func newFileChecker(timeout time.Duration, r string, s semaphore) fileChecker {
	return fileChecker{newURLChecker(timeout, s), r}
}

func (c fileChecker) Check(f string) ([]urlResult, error) {
	n, err := parseFile(f)

	if err != nil {
		return nil, err
	}

	us, err := c.extractURLs(n)

	if err != nil {
		return nil, err
	}

	rc := make(chan urlResult, len(us))
	rs := make([]urlResult, 0, len(us))

	go c.urlChecker.CheckMany(us, f, rc)

	for r := range rc {
		rs = append(rs, r)
	}

	return rs, nil
}

func (c fileChecker) CheckMany(fc <-chan string, rc chan<- fileResult) {
	wg := sync.WaitGroup{}

	for f := range fc {
		wg.Add(1)

		go func(f string) {
			if rs, err := c.Check(f); err == nil {
				rc <- fileResult{filename: f, urlResults: rs}
			} else {
				rc <- fileResult{filename: f, err: err}
			}

			wg.Done()
		}(f)
	}

	wg.Wait()
	close(rc)
}

func parseFile(f string) (*html.Node, error) {
	bs, err := ioutil.ReadFile(f)

	if err != nil {
		return nil, err
	}

	if !isHTMLFile(f) {
		bs = blackfriday.Run(bs)
	}

	n, err := html.Parse(bytes.NewReader(bs))

	if err != nil {
		return nil, err
	}

	return n, nil
}

func (c fileChecker) extractURLs(n *html.Node) ([]string, error) {
	us := make(map[string]bool)
	ns := []*html.Node{}
	ns = append(ns, n)

	addURL := func(u string) error {
		u, err := c.resolveURL(u)

		if err != nil {
			return err
		}

		us[u] = true

		return nil
	}

	for len(ns) > 0 {
		i := len(ns) - 1
		n := ns[i]
		ns = ns[:i]

		if n.Type == html.ElementNode {
			switch n.Data {
			case "a":
				for _, a := range n.Attr {
					if a.Key == "href" && isURL(a.Val) {
						if err := addURL(a.Val); err != nil {
							return nil, err
						}

						break
					}
				}
			case "img":
				for _, a := range n.Attr {
					if a.Key == "src" && isURL(a.Val) {
						if err := addURL(a.Val); err != nil {
							return nil, err
						}

						break
					}
				}
			}
		}

		for n := n.FirstChild; n != nil; n = n.NextSibling {
			ns = append(ns, n)
		}
	}

	return stringSetToSlice(us), nil
}

func (c fileChecker) resolveURL(u string) (string, error) {
	abs := strings.HasPrefix(u, "/")

	if abs && c.documentRoot != "" {
		return path.Join(c.documentRoot, u), nil
	} else if abs {
		return "", errors.New("document root directory is not specified")
	}

	return u, nil
}

func isURL(s string) bool {
	if strings.HasPrefix(s, "#") {
		return false
	}

	u, err := url.Parse(s)
	return err == nil && (u.Scheme == "" || u.Scheme == "http" || u.Scheme == "https")
}
