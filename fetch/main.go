// +build js,wasm

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dennwc/dom"
	"github.com/johanbrandhorst/wasm-experiments/div"
)

func main() {
	t := dom.GetDocument().GetElementById("target")
	logger := log.New((*div.Writer)(t), "", log.LstdFlags)

	c := http.Client{}
	req, err := http.NewRequest(
		"POST",
		"https://httpbin.org/anything",
		strings.NewReader(`{"test":"test"}`),
	)
	if err != nil {
		logger.Fatal(err)
	}

	resp, err := c.Do(req)
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Print(string(b))
}
