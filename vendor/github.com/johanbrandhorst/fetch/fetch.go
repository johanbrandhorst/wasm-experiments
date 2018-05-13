// +build js,wasm

package fetch

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"syscall/js"
)

// Adapted for syscall/js from
// https://github.com/gopherjs/gopherjs/blob/8dffc02ea1cb8398bb73f30424697c60fcf8d4c5/compiler/natives/src/net/http/fetch.go

// streamReader implements an io.ReadCloser wrapper for ReadableStream of https://fetch.spec.whatwg.org/.
type streamReader struct {
	pending []byte
	stream  js.Value
}

func (r *streamReader) Read(p []byte) (n int, err error) {
	if len(r.pending) == 0 {
		var (
			bCh   = make(chan []byte)
			errCh = make(chan error)
		)
		r.stream.Call("read").Call("then",
			js.NewCallback(func(args []js.Value) {
				result := args[0]
				if result.Get("done").Bool() {
					errCh <- io.EOF
					return
				}
				value := make([]byte, result.Get("value").Get("byteLength").Int())
				js.ValueOf(value).Call("set", result.Get("value"))
				bCh <- value
			}),
			js.NewCallback(func(args []js.Value) {
				// Assumes it's a DOMException.
				errCh <- errors.New(args[0].Get("message").String())
			}),
		)
		select {
		case b := <-bCh:
			r.pending = b
		case err := <-errCh:
			return 0, err
		}
	}
	n = copy(p, r.pending)
	r.pending = r.pending[n:]
	return n, nil
}

func (r *streamReader) Close() error {
	// This ignores any error returned from cancel method. So far, I did not encounter any concrete
	// situation where reporting the error is meaningful. Most users ignore error from resp.Body.Close().
	// If there's a need to report error here, it can be implemented and tested when that need comes up.
	r.stream.Call("cancel")
	return nil
}

// Transport is a RoundTripper that is implemented using the WHATWG Fetch API.
// It supports streaming response bodies.
type Transport struct{}

// RoundTrip performs a full round trip of a request.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	headers := js.Global.Get("Headers").New()
	for key, values := range req.Header {
		for _, value := range values {
			headers.Call("append", key, value)
		}
	}

	ac := js.Global.Get("AbortController").New()

	opt := js.Global.Get("Object").New()
	opt.Set("headers", headers)
	opt.Set("method", req.Method)
	opt.Set("credentials", "same-origin")
	opt.Set("signal", ac.Get("signal"))

	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			_ = req.Body.Close() // RoundTrip must always close the body, including on errors.
			return nil, err
		}
		_ = req.Body.Close()
		opt.Set("body", body)
	}
	respPromise := js.Global.Call("fetch", req.URL.String(), opt)
	if respPromise == js.Undefined {
		return nil, errors.New("your browser does not support the Fetch API, please upgrade")
	}

	var (
		respCh = make(chan *http.Response)
		errCh  = make(chan error)
	)
	respPromise.Call("then",
		js.NewCallback(func(args []js.Value) {
			result := args[0]
			header := http.Header{}
			result.Get("headers").Call("forEach", js.NewCallback(func(args []js.Value) {
				key, value := args[0].String(), args[1].String()
				ck := http.CanonicalHeaderKey(key)
				header[ck] = append(header[ck], value)
			}))

			contentLength := int64(-1)
			if cl, err := strconv.ParseInt(header.Get("Content-Length"), 10, 64); err == nil {
				contentLength = cl
			}

			b := result.Get("body")
			if b == js.Undefined {
				errCh <- errors.New("your browser does not support the ReadableStream API, please upgrade")
				return
			}

			select {
			case respCh <- &http.Response{
				Status:        result.Get("status").String() + " " + http.StatusText(result.Get("status").Int()),
				StatusCode:    result.Get("status").Int(),
				Header:        header,
				ContentLength: contentLength,
				Body:          &streamReader{stream: b.Call("getReader")},
				Request:       req,
			}:
			case <-req.Context().Done():
			}
		}),
		js.NewCallback(func(args []js.Value) {
			select {
			case errCh <- fmt.Errorf("net/http: fetch() failed: %s", args[0].String()):
			case <-req.Context().Done():
			}
		}),
	)
	select {
	case <-req.Context().Done():
		// Abort the Fetch request
		ac.Call("abort")
		return nil, errors.New("net/http: request canceled")
	case resp := <-respCh:
		return resp, nil
	case err := <-errCh:
		return nil, err
	}
}
