//go:build !solution

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func fetchURL(url string, ch chan bool) {
	start := time.Now()

	resp, err := http.Get(url)
	end := time.Since(start)

	if err == nil {
		defer resp.Body.Close()
		fmt.Printf("%v\t%v %v\n", end, resp.ContentLength, url)
	} else {
		fmt.Printf("error in %v", url)
	}

	ch <- true
}

func main() {
	start := time.Now()

	urls := os.Args[1:]

	ch := make(chan bool)

	for _, url := range urls {
		go fetchURL(url, ch)
	}

	n := len(urls)

	for i := 0; i < n; i++ {
		<-ch
	}

	end := time.Since(start)

	fmt.Printf("%v elapsed", end)
}
