// +build js,wasm

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/johanbrandhorst/fetch"
)

func main() {
	c := http.Client{
		Transport: &fetch.Transport{},
	}
	resp, err := c.Get("https://api.github.com")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}
