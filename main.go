package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/a8m/mark"
	"github.com/docopt/docopt-go"
	"github.com/fatih/color"
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

	bs, err := ioutil.ReadFile(args["<filename>"].(string))

	if err != nil {
		panic(err)
	}

	n, err := html.Parse(strings.NewReader(mark.Render(string(bs))))

	if err != nil {
		panic(err)
	}

	checkURLs(extractURLs(n), args["--verbose"].(bool))
}

func checkURLs(ss []string, verbose bool) {
	client := &http.Client{Timeout: 5 * time.Second}
	bs := make(chan bool, len(ss))

	for _, s := range ss {
		go checkURL(client, s, bs, verbose)
	}

	ok := true

	for i := 0; i < len(ss); i++ {
		ok = ok && <-bs
	}

	if !ok {
		os.Exit(1)
	}
}

func checkURL(client *http.Client, s string, bs chan bool, verbose bool) {
	_, err := client.Get(s)

	if s := color.New(color.FgCyan).SprintFunc()(s); err != nil {
		printToStderr(
			color.New(color.FgRed).SprintFunc()("ERROR") + "\t" + s + "\t" + err.Error())
	} else if err == nil && verbose {
		printToStderr(color.New(color.FgGreen).SprintFunc()("OK") + "\t" + s)
	}

	bs <- err == nil
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

func stringSetToSlice(s2b map[string]bool) []string {
	ss := make([]string, 0, len(s2b))

	for s := range s2b {
		ss = append(ss, s)
	}

	return ss
}

func isURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https")
}

func getArgs() map[string]interface{} {
	usage := `Link checker for Markdown and HTML

Usage:
	linkcheck [-v] <filename>

Options:
	-v, --verbose  Be verbose`

	args, err := docopt.Parse(usage, nil, true, "linkcheck", true)

	if err != nil {
		panic(err)
	}

	return args
}

func printToStderr(xs ...interface{}) {
	fmt.Fprintln(os.Stderr, xs...)
}
