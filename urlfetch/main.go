//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func check(e error) {
	if e != nil {
		os.Exit(1)
	}
}

func main() {
	urls := os.Args[1:]

	for _, v := range urls {
		resp, err := http.Get(v)
		check(err)
		body, err := io.ReadAll(resp.Body)
		check(err)
		fmt.Println(string(body))

		resp.Body.Close()
	}
}
