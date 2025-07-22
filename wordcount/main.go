//go:build !solution

package main

import (
	"bufio"
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fileNames := os.Args[1:]

	m := make(map[string]int)

	for _, fileName := range fileNames {
		file, err := os.Open(fileName)
		check(err)

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()

			m[line]++
		}

		file.Close()
	}

	for k, v := range m {
		if v > 1 {
			fmt.Printf("%v\t%v\n", v, k)
		}
	}
}
