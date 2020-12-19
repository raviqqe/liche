## DEPRECATION NOTICE

Sorry this project is not actively maintained anymore! 😢 Please consider migrating to one of the alternatives listed below.

### Alternatives

- [muffet](https://github.com/raviqqe/muffet)
  - Fast website link checker in Go
- [lychee](https://github.com/lycheeverse/lychee)
  - A glorious link checker
  - This tool supports testing links both in local files and on websites.

### Why is it not maintained anymore?

It's because we found several problems with the goals of the project and the amount of work it needs.

The goal of this software was originally to check links in Markdown files which are **compiled into HTML files and served via HTTP servers**. But that raises the following problems.

- We cannot test links not in the Markdown files.
  - For example, some markdown-based static site generators generate links at compile time.
  - e.g. Automatic generation of table of contents
- We cannot test the behaviour of HTTP servers.
  - Different HTTP servers handles URLs differently.
  - e.g. trailing slashes, inference of page file extensions, ...

It needs a lot of work to support all these different use cases. In short, we need to emulate different HTTP servers as well as web browsers.

### But I still want this...

If you think this software is still valuable for you even in comparison with the alternatives listed above and want it to be maintained, please let us know by posting a new issue.

# liche

[![Circle CI](https://img.shields.io/circleci/project/github/raviqqe/liche.svg?style=flat-square)](https://circleci.com/gh/raviqqe/liche)
[![Codecov](https://img.shields.io/codecov/c/github/raviqqe/liche.svg?style=flat-square)](https://codecov.io/gh/raviqqe/liche)
[![Go Report Card](https://goreportcard.com/badge/github.com/raviqqe/liche?style=flat-square)](https://goreportcard.com/report/github.com/raviqqe/liche)
[![License](https://img.shields.io/github/license/raviqqe/liche.svg?style=flat-square)](LICENSE)

[![asciicast](https://asciinema.org/a/148896.png)](https://asciinema.org/a/148896)

`liche` is a command to check links' connectivity in Markdown and HTML files.
It checks all `a` and `img` tags in specified files.

## Installation

```sh
go get -u github.com/raviqqe/liche
```

- requires [Go Modules]("https://github.com/golang/go/wiki/Modules#how-to-use-modules")

## Usage

```sh
> liche --help
Link checker for Markdown and HTML

Usage:
	liche [-c <num-requests>] [-d <directory>] [-r] [-t <timeout>] [-x <regex>] [-v] <filenames>...

Options:
	-c, --concurrency <num-requests>  Set max number of concurrent HTTP requests. [default: 512]
	-d, --document-root <directory>  Set document root directory for absolute paths.
	-r, --recursive  Search Markdown and HTML files recursively
	-t, --timeout <timeout>  Set timeout for HTTP requests in seconds. Disabled by default.
	-x, --exclude <regex>  Regex of links to exclude from checking.
	-v, --verbose  Be verbose.
```

## Examples

```sh
> liche file.md
> liche file1.md file2.md
> liche file.html
> liche file1.html file2.html
> liche -r directory # Search all Markdown and HTML files recursively.
```

## Supported properties

- File extensions: `.md`, `.html`, `.htm`
- HTML tags: `a`, `img`
- HTML attributes: `href`, `src`
- URL schemes: `http`, `https`

Also supports relative and absolute paths.
(Absolute paths need `--document-root` option.)

## License

[MIT](LICENSE)
