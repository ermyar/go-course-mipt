//go:build !solution

package main

import (
	"fmt"
	"os"
)

var lookuptable []map[string]int = []map[string]int{
	{
		"lines":   1,
		"commits": 2,
		"files":   3,
	}, {
		"tabular":    1,
		"csv":        2,
		"json":       3,
		"json-lines": 4,
	},
}

type gitfame struct {
	path         string
	commit       string
	useCommitter bool
	order        string
	format       string
	exclude      []string
	restricted   []string
	extensions   map[string]bool
	result       map[string]stat
}

func main() {
	var gf gitfame
	err := gf.parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %s\n", err.Error())
		os.Exit(-1)
	}
	err = gf.compute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "compute error: %s\n", err.Error())
		os.Exit(-1)
	}
	err = gf.output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "output error: %s\n", err.Error())
	}
}
