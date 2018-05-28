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
	req, err := http.NewRequest(
		"POST",
		"https://httpbin.org/anything",
		strings.NewReader(`{"test":"test"}`),
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	/*
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
		defer cancel()
		req = req.WithContext(ctx)
	*/
	resp, err := c.Do(req)
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
