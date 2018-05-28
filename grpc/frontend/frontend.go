package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"

	"github.com/johanbrandhorst/fetch"
	"github.com/johanbrandhorst/wasm-experiments/grpc/proto/server"
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
	b, err := proto.Marshal(&server.GetUserRequest{
		UserId: "1234",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	bufHeader := make([]byte, 5)

	// Write length of b into buf
	binary.BigEndian.PutUint32(bufHeader[1:], uint32(len(b)))

	req, err := http.NewRequest("POST", "/web.Backend/GetUser", bytes.NewBuffer(append(bufHeader, b...)))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("content-type", "application/grpc-web+proto")
	//ctx, _ := context.WithTimeout(context.Background(), time.Second)
	//req = req.WithContext(ctx)

	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	for {
		header := make([]byte, 5)
		_, err := resp.Body.Read(header)
		if err != nil {
			fmt.Println(err)
			return
		}
		if header[0] == 0x80 {
			trailers, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(trailers))
			return
		}

		length := binary.BigEndian.Uint32(header[1:])

		message := make([]byte, length)
		_, err = resp.Body.Read(message)
		if err != nil {
			fmt.Println(err)
			return
		}
		/*
			status := resp.Header.Get("grpc-status")
			statusCode, err := strconv.Atoi(status)
			if err != nil {
				fmt.Println(err)
				return
			}
			code := codes.Code(statusCode)
			if code != codes.OK {
				msg := resp.Header.Get("grpc-message")
				fmt.Println(msg)
				return
			}
		*/
		user := new(server.User)
		err = proto.Unmarshal(message, user)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(user.Id)
	}
}
