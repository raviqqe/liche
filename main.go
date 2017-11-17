package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/a8m/mark"
	docopt "github.com/docopt/docopt-go"
	"golang.org/x/net/html"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, r.(error).Error())
			os.Exit(1)
		}
	}()

	args := getArgs()

	bs, err := ioutil.ReadFile(args["<filename>"].(string))

	if err != nil {
		panic(err)
	}

	html.Parse(strings.NewReader(mark.Render(string(bs))))
}

func extractUrls(n *html.Node) map[string]bool {
	ss := make(map[string]bool)
	ns := make([]*html.Node, 0, 1024)
	ns = append(ns, n)

	for len(ns) > 0 {
		i := len(ns) - 1
		n := ns[i]
		ns = ns[:i]

		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					ss[a.Val] = true
					break
				}
			}
		}

		for n := n.FirstChild; n != nil; n = n.NextSibling {
			ns = append(ns, n)
		}
	}

	return ss
}

func getArgs() map[string]interface{} {
	usage := `Link checker for Markdown and HTML

Usage:
	linkcheck <filename>`

	args, err := docopt.Parse(usage, nil, true, "linkcheck", true)

	if err != nil {
		panic(err)
	}

	return args
}
