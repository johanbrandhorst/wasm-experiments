package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/johanbrandhorst/fetch"
)

// Build with Go WASM fork
//go:generate rm -f ./html/test.wasm
//go:generate bash -c "GOOS=js GOARCH=wasm GOROOT=$GOPATH/src/github.com/neelance/go/ $GOPATH/src/github.com/neelance/go/bin/go build -o ./html/test.wasm frontend.go"

// Integrate generated JS into a Go file for static loading.
//go:generate bash -c "go run assets_generate.go"

func main() {
	c := http.Client{
		Transport: &fetch.Transport{},
	}
	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Path: "/web.Backend/GetUser",
		},
		Header: http.Header{
			"Content-Type": []string{"application/grpc-web+proto"},
		},
	}
	//ctx, _ := context.WithTimeout(context.Background(), time.Second)
	//req = req.WithContext(ctx)

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
