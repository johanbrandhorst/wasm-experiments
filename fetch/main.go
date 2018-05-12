// +build js,wasm

package main

import (
	"net/http"
	"net/url"

	"github.com/johanbrandhorst/wasm-experiments/fetch/fetch"
)

func main() {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	req := fetch.Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   "httpbin.org",
			Path:   "anything",
		},
		Method:  fetch.POST,
		Body:    []byte(`{"key": "value"}`),
		Headers: headers,
	}
	fetch.Fetch(req)
}
