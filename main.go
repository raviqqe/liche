package main

import (
	"net/url"
	"os"
	"time"

	"github.com/docopt/docopt-go"
	"golang.org/x/net/html"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			printToStderr(r.(error).Error())
			os.Exit(1)
		}
	}()

	args := getArgs()
	fs := args["<filenames>"].([]string)
	bs := make(chan bool, len(fs))
	c := newFileChecker(5*time.Second, args["--verbose"].(bool))

	for _, f := range fs {
		go func(f string) {
			bs <- c.Check(f)
		}(f)
	}

	ok := true

	for i := 0; i < len(fs); i++ {
		ok = <-bs && ok
	}

	if !ok {
		os.Exit(1)
	}
}

func extractURLs(n *html.Node) []string {
	ss := make(map[string]bool)
	ns := make([]*html.Node, 0, 1024)
	ns = append(ns, n)

	for len(ns) > 0 {
		i := len(ns) - 1
		n := ns[i]
		ns = ns[:i]

		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && isURL(a.Val) {
					ss[a.Val] = true
					break
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

func getArgs() map[string]interface{} {
	usage := `Link checker for Markdown and HTML

	Usage:
	linkcheck [-v] <filenames>...

	Options:
	-v, --verbose  Be verbose`

	args, err := docopt.Parse(usage, nil, true, "linkcheck", true)

	if err != nil {
		panic(err)
	}

	return args
}
