// +build js,wasm

package main

import (
	"fmt"
	"net/http"

	"github.com/johanbrandhorst/wasm-experiments/fetch/fetch"
)

func main() {
	c := http.Client{
		Transport: &fetch.Transport{},
	}
	fmt.Println(c.Get("https://api.github.com"))
}
