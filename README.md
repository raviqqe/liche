# liche

[![Circle CI](https://img.shields.io/circleci/project/github/raviqqe/liche.svg?style=flat-square)](https://circleci.com/gh/raviqqe/liche)
[![Go Report Card](https://goreportcard.com/badge/github.com/raviqqe/liche?style=flat-square)](https://goreportcard.com/report/github.com/raviqqe/liche)
[![License](https://img.shields.io/github/license/raviqqe/liche.svg?style=flat-square)](LICENSE)

[![asciicast](https://asciinema.org/a/148134.png)](https://asciinema.org/a/148134)

`liche` is a command to check links' connectivity in Markdown and HTML files.
It checks all `a` and `img` tags in specified files.

## Installation

```sh
go get -u github.com/raviqqe/liche
```

## Usage

```sh
> liche --help
Link checker for Markdown and HTML

Usage:
	liche [-c <num-requests>] [-t <timeout>] [-v] <filenames>...

Options:
	-c, --concurrency <num-requests>  Set max number of concurrent HTTP requests. [default: 32]
	-t, --timeout <timeout>  Set timeout for HTTP requests in seconds. Disabled by default.
	-v, --verbose  Be verbose.
```

### Markdown

```sh
> liche file.md
> liche file1.md file2.md
```

### HTML

```sh
> liche file.html
> liche file1.html file2.html
```

## Supported properties

- HTML tags: `a`, `img`
- HTML attributes: `href`, `src`
- URL schemes: `http`, `https`

## License

[MIT](LICENSE)
