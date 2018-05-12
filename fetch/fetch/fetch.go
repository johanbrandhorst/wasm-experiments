// +build js,wasm

package fetch

//go:generate bash -c "GOROOT=$GOPATH/src/github.com/neelance/go/ GOOS=js GOARCH=wasm stringer -type Method"

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"syscall/js"
	"syscall/js/callback"
)

type Method int

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods
const (
	GET Method = iota
	HEAD
	POST
	PUT
	DELETE
	CONNECT
	OPTIONS
	TRACE
	PATCH
)

type Request struct {
	Method  Method
	URL     *url.URL
	Headers http.Header
	Body    io.Reader
}

type Response struct {
	URL        *url.URL
	Headers    http.Header
	Body       io.ReadCloser
	StatusCode int
}

func Fetch(req *Request) (*Response, error) {
	init := js.Global.Get("Object").New()
	init.Set("method", req.Method.String())
	if req.Body != nil {
		switch req.Method {
		case GET, HEAD:
			return nil, errors.New("cannot use body with GET or HEAD HTTP methods")
		}
		init.Set("body", req.Body)
	}

	headers := js.Global.Get("Headers").New()
	for header, values := range req.Headers {
		for _, value := range values {
			headers.Call("append", header, value)
		}
	}
	init.Set("headers", headers)

	cb := callback.New(func(args []js.Value) {
		fmt.Println(args)
	})

	promise := js.Global.Call("fetch", req.URL.String(), init)
	promise.Call("then", cb)
	return &Response{}, nil
}
