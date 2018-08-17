// +build js,wasm

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/dennwc/dom"
)

type writer dom.Element

// Write implements io.Writer.
func (d writer) Write(p []byte) (n int, err error) {
	node := dom.GetDocument().CreateElement("div")
	node.SetTextContent(string(p))
	(*dom.Element)(&d).AppendChild(node)
	return len(p), nil
}

func main() {
	t := dom.GetDocument().GetElementById("target")
	logger := log.New((*writer)(t), "", log.LstdFlags)

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
