//go:build !solution

package main

import (
	"fmt"
	"os"
	"strings"
)

var m map[string]int

func Count(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	// for real tried to use string.Field
	for _, val := range strings.Split(string(data), "\n") {
		for _, str := range strings.Split(val, " ") {
			m[str]++
		}
	}
}

func main() {
	m = make(map[string]int)
	args := os.Args[1:]
	for _, file := range args {
		Count(file)
	}
	for word, cnt := range m {
		if cnt > 1 {
			fmt.Printf("%d\t%s\n", cnt, word)
		}
	}
}
