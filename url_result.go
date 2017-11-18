package main

import "github.com/fatih/color"

type urlResult struct {
	url string
	err error
}

func (r urlResult) String() string {
	if r.err == nil {
		return color.GreenString("OK") + "\t" + r.url
	}

	return color.RedString("ERROR") + "\t" + r.url + "\t" + color.YellowString(r.err.Error())
}
