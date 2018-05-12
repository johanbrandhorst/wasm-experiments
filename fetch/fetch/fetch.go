// +build js,wasm

package fetch

//go:generate bash -c "GOROOT=$GOPATH/src/github.com/neelance/go/ GOOS=js GOARCH=wasm stringer -type Method"

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
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

/*

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
	Ok         bool
}

func Fetch(ctx context.Context, req *Request) (*Response, error) {
	init := js.Global.Get("Object").New()
	init.Set("method", req.Method.String())
	if req.Body != nil {
		switch req.Method {
		case GET, HEAD:
			return nil, errors.New("cannot use body with GET or HEAD HTTP methods")
		}
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		init.Set("body", b)
	}

	headers := js.Global.Get("Headers").New()
	for header, values := range req.Headers {
		for _, value := range values {
			headers.Call("append", header, value)
		}
	}
	init.Set("headers", headers)

	ac := js.Global.Get("AbortController").New()
	init.Set("signal", ac.Get("signal"))

	promise := js.Global.Call("fetch", req.URL.String(), init)

	wait := make(chan Response)
	errChan := make(chan error)
	cb := callback.New(func(args []js.Value) {
		jsResp := args[0]
		resp := Response{
			Headers: http.Header{},
		}
		u, err := url.Parse(jsResp.Get("url").String())
		if err != nil {
			errChan <- err
		}
		jsResp.Get("headers").Call("forEach", callback.New(func(args []js.Value) {
			key, val := args[0].String(), args[1].String()
			resp.Headers[key] = append(resp.Headers[key], val)
		}).Value)


		wait <- Response{
			Ok:  jsResp.Get("ok").Bool(),
			URL: u,
		}
	})

	promise.Call("then", cb.Value)

	select {
	case <-ctx.Done():
		// Abort the Fetch request
		ac.Call("abort")
		return nil, ctx.Err()
	case resp := <-wait:
		return &resp, nil
	case err := <-errChan:
		return nil, err
	}
}
*/

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
			func(result js.Value) {
				if result.Get("done").Bool() {
					errCh <- io.EOF
					return
				}
				var value []byte
				// TODO: Any way to avoid this copying?
				result.Get("value").Call("forEach", callback.New(func(args []js.Value) {
					value = append(value, args[0].String()[0])
				}).Value)
				bCh <- value
			},
			func(reason js.Value) {
				// Assumes it's a DOMException.
				errCh <- errors.New(reason.Get("message").String())
			},
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

// Transport is a RoundTripper that is implemented using Fetch API. It supports streaming
// response bodies.
type Transport struct{}

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
		// TODO: Find out if request body can be streamed into the fetch request rather than in advance here.
		//       See BufferSource at https://fetch.spec.whatwg.org/#body-mixin.
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			_ = req.Body.Close() // RoundTrip must always close the body, including on errors.
			return nil, err
		}
		_ = req.Body.Close()
		opt.Set("body", body)
	}
	respPromise := js.Global.Call("fetch", req.URL.String(), opt)

	var (
		respCh = make(chan *http.Response)
		errCh  = make(chan error)
	)
	respPromise.Call("then",
		callback.New(func(args []js.Value) {
			result := args[0]
			header := http.Header{}
			result.Get("headers").Call("forEach", callback.New(func(args []js.Value) {
				key, value := args[0].String(), args[1].String()
				ck := http.CanonicalHeaderKey(key)
				header[ck] = append(header[ck], value)
			}).Value)

			contentLength := int64(-1)
			if cl, err := strconv.ParseInt(header.Get("Content-Length"), 10, 64); err == nil {
				contentLength = cl
			}

			select {
			case respCh <- &http.Response{
				Status:        result.Get("status").String() + " " + http.StatusText(result.Get("status").Int()),
				StatusCode:    result.Get("status").Int(),
				Header:        header,
				ContentLength: contentLength,
				Body:          &streamReader{stream: result.Get("body").Call("getReader")},
				Request:       req,
			}:
			case <-req.Context().Done():
			}
		}).Value,
		callback.New(func(args []js.Value) {
			select {
			case errCh <- fmt.Errorf("net/http: fetch() failed: %s", args[0].String()):
			case <-req.Context().Done():
			}
		}).Value,
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
