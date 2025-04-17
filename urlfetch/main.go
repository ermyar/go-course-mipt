//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func fetch(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Printf("%s", body)
	return err
}

func main() {
	args := os.Args[1:]

	for _, val := range args {
		err := fetch(val)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
}
