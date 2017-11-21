package main

import "os"

func main() {
	args, err := getArgs()

	if err != nil {
		printToStderr(err.Error())
		os.Exit(1)
	}

	if args.recursive {
		args.filenames, err = listFilesRecursively(args.filenames)

		if err != nil {
			printToStderr(err.Error())
			os.Exit(1)
		}
	}

	rc := make(chan fileResult, len(args.filenames))
	s := newSemaphore(args.concurrency)
	c := newFileChecker(args.timeout, s)

	go c.CheckMany(args.filenames, rc)

	ok := true

	for r := range rc {
		if !r.Ok() {
			ok = false
			printToStderr(r.String(args.verbose))
		} else if args.verbose {
			printToStderr(r.String(true))
		}
	}

	if !ok {
		os.Exit(1)
	}
}

func listFilesRecursively(fs []string) ([]string, error) {
	gs := []string{}

	for _, f := range fs {
		i, err := os.Stat(f)

		if err != nil {
			return nil, err
		}

		if i.IsDir() {
			gs = append(gs, listFiles(f)...)
		} else {
			gs = append(gs, f)
		}
	}

	return gs, nil
}
