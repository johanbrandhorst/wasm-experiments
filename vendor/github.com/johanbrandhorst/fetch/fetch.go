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

// Transport is a RoundTripper that is implemented using the WHATWG Fetch API.
// It supports streaming response bodies.
type Transport struct{}

// RoundTrip performs a full round trip of a request.
func (*Transport) RoundTrip(req *http.Request) (*http.Response, error) {
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

	var (
		respCh = make(chan *http.Response, 1)
		errCh  = make(chan error, 1)
	)
	if req.Body != nil {
		/* Streaming request bodies are not supported yet
		body := js.Global.Get("ReadableStream")
		if body != js.Undefined {
			source := js.Global.Get("Object").New()
			// TODO(johanbrandhorst): Use ReadableByteStreamController.
			// Currently Unsupported: https://developer.mozilla.org/en-US/docs/Web/API/ReadableByteStreamController#Browser_Compatibility
			start := js.NewCallback(func(args []js.Value) {
				fmt.Println("start called")
				controller := args[0]
				w := &streamWriter{controller: controller}
				_, err := io.Copy(w, req.Body)
				if err != nil {
					errCh <- err
					return
				}
			})
			defer start.Close()
			source.Set("start", start)
			body = js.Global.Get("ReadableStream").New(source)
		}
		*/
		content, err := ioutil.ReadAll(req.Body)
		if err != nil {
			req.Body.Close() // RoundTrip must always close the body, including on errors.
			return nil, err
		}
		req.Body.Close()
		opt.Set("body", js.ValueOf(content))
	}
	respPromise := js.Global.Call("fetch", req.URL.String(), opt)
	if respPromise == js.Undefined {
		return nil, errors.New("your browser does not support the Fetch API, please upgrade")
	}

	success := js.NewCallback(func(args []js.Value) {
		result := args[0]
		header := http.Header{}
		writeHeaders := js.NewCallback(func(args []js.Value) {
			key, value := args[0].String(), args[1].String()
			ck := http.CanonicalHeaderKey(key)
			header[ck] = append(header[ck], value)
		})
		defer writeHeaders.Close()
		result.Get("headers").Call("forEach", writeHeaders)

		contentLength := int64(-1)
		if cl, err := strconv.ParseInt(header.Get("Content-Length"), 10, 64); err == nil {
			contentLength = cl
		}

		b := result.Get("body")
		var body io.ReadCloser
		if b != js.Undefined {
			body = &streamReader{stream: b.Call("getReader")}
		} else {
			// Fall back to using the arrayBuffer
			// https://developer.mozilla.org/en-US/docs/Web/API/Body/arrayBuffer
			body = &arrayReader{arrayPromise: result.Call("arrayBuffer")}
		}
		select {
		case respCh <- &http.Response{
			Status:        result.Get("status").String() + " " + http.StatusText(result.Get("status").Int()),
			StatusCode:    result.Get("status").Int(),
			Header:        header,
			ContentLength: contentLength,
			Body:          body,
			Request:       req,
		}:
		case <-req.Context().Done():
		}
	})
	defer success.Close()
	failure := js.NewCallback(func(args []js.Value) {
		select {
		case errCh <- fmt.Errorf("net/http: fetch() failed: %s", args[0].String()):
		case <-req.Context().Done():
		}
	})
	defer failure.Close()
	respPromise.Call("then", success, failure)
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

// streamReader implements an io.ReadCloser wrapper for ReadableStream of https://fetch.spec.whatwg.org/.
type streamReader struct {
	pending []byte
	stream  js.Value
}

func (r *streamReader) Read(p []byte) (n int, err error) {
	if len(r.pending) == 0 {
		var (
			bCh   = make(chan []byte, 1)
			errCh = make(chan error, 1)
		)
		success := js.NewCallback(func(args []js.Value) {
			result := args[0]
			if result.Get("done").Bool() {
				errCh <- io.EOF
				return
			}
			bCh <- copyBytes(result.Get("value"))
		})
		defer success.Close()
		failure := js.NewCallback(func(args []js.Value) {
			// Assumes it's a DOMException.
			errCh <- errors.New(args[0].Get("message").String())
		})
		defer failure.Close()
		r.stream.Call("read").Call("then", success, failure)
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

// arrayReader implements an io.ReadCloser wrapper for arrayBuffer
// https://developer.mozilla.org/en-US/docs/Web/API/Body/arrayBuffer.
type arrayReader struct {
	arrayPromise js.Value
	pending      []byte
	read         bool
}

func (r *arrayReader) Read(p []byte) (n int, err error) {
	if !r.read {
		r.read = true
		var (
			bCh   = make(chan []byte, 1)
			errCh = make(chan error, 1)
		)
		success := js.NewCallback(func(args []js.Value) {
			// Wrap the input ArrayBuffer with a Uint8Array
			uint8arrayWrapper := js.Global.Get("Uint8Array").New(args[0])
			bCh <- copyBytes(uint8arrayWrapper)
		})
		defer success.Close()
		failure := js.NewCallback(func(args []js.Value) {
			// Assumes it's a DOMException.
			errCh <- errors.New(args[0].Get("message").String())
		})
		defer failure.Close()
		r.arrayPromise.Call("then", success, failure)
		select {
		case b := <-bCh:
			r.pending = b
		case err := <-errCh:
			return 0, err
		}
	}
	if len(r.pending) == 0 {
		return 0, io.EOF
	}
	n = copy(p, r.pending)
	r.pending = r.pending[n:]
	return n, nil
}

func (r *arrayReader) Close() error {
	// This is a noop
	return nil
}

/*
// streamWriter exposes a ReadableStreamDefaultController as an io.Writer
// https://developer.mozilla.org/en-US/docs/Web/API/ReadableStreamDefaultController
type streamWriter struct {
	controller js.Value
}

func (w *streamWriter) Write(p []byte) (int, error) {
	w.controller.Call("enqueue", p)
	return len(p), nil
}
*/

func copyBytes(in js.Value) []byte {
	value := make([]byte, in.Get("byteLength").Int())
	js.ValueOf(value).Call("set", in)
	return value
}
