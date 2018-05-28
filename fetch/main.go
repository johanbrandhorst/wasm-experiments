// +build js,wasm

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/johanbrandhorst/fetch"
)

func main() {
	c := http.Client{
		Transport: &fetch.Transport{},
	}
	resp, err := c.Post(
		"https://httpbin.org/anything",
		"application/json",
		strings.NewReader(`{"test":"test"}`),
	)
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
