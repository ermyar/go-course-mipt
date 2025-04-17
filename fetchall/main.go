//go:build !solution

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type request struct {
	url string
	err error
	dur time.Duration
}

func fetch(url string, c chan request) {
	timer := time.Now()
	_, err := http.Get(url)
	c <- request{url, err, time.Since(timer)}
}

func main() {
	args := os.Args[1:]
	c := make(chan request, len(args))
	timer := time.Now()

	for _, url := range args {
		go fetch(url, c)
	}

	for times := 0; times < len(args); times++ {
		tmp := <-c
		if tmp.err != nil {
			fmt.Println("failed", tmp.url, tmp.err)
		} else {
			fmt.Println(tmp.dur, tmp.url)
		}
	}
	fmt.Println("elapsed", time.Since(timer))
}
