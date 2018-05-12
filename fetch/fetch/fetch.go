// +build js,wasm

package fetch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"syscall/js"
	"syscall/js/callback"
)

type Method int

func (m Method) String() string {
	switch m {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case PATCH:
		return "PATCH"
	}

	panic("unsupported method")
}

const (
	GET Method = iota
	POST
	PUT
	DELETE
	PATCH
)

type Request struct {
	Method  Method
	URL     *url.URL
	Headers http.Header
	Body    []byte
}

type Response struct {
	URL        url.URL
	Headers    http.Header
	Body       io.Reader
	StatusCode int
}

func Fetch(req Request) (Response, error) {
	init := js.Global.Get("Object").New()
	init.Set("method", req.Method.String())
	init.Set("body", req.Body)

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
	return Response{}, nil
}
