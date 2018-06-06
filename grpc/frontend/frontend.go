package main

import (
	"context"
	"fmt"
	"io"

	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"

	grpc "github.com/johanbrandhorst/grpc-wasm"
	server "github.com/johanbrandhorst/wasm-experiments/grpc/proto/client"
)

// Build with Go WASM fork
//go:generate rm -f ./html/test.wasm
//go:generate bash -c "GOOS=js GOARCH=wasm GOROOT=$GOPATH/src/github.com/neelance/go/ $GOPATH/src/github.com/neelance/go/bin/go build -o ./html/test.wasm frontend.go"

// Integrate generated JS into a Go file for static loading.
//go:generate bash -c "go run assets_generate.go"

func main() {
	cc, _ := grpc.Dial("")
	client := server.NewBackendClient(cc)
	resp, err := client.GetUser(context.Background(), &server.GetUserRequest{
		UserId: "1234",
	})
	if err != nil {
		st := status.Convert(err)
		fmt.Println(st.Code(), st.Message(), st.Details())
	} else {
		fmt.Println(resp.GetId())
	}
	resp, err = client.GetUser(context.Background(), &server.GetUserRequest{
		UserId: "123",
	})
	if err != nil {
		st := status.Convert(err)
		fmt.Println(st.Code(), st.Message(), st.Details())
	} else {
		fmt.Println(resp.GetId())
	}

	srv, err := client.GetUsers(context.Background(), &server.GetUsersRequest{
		NumUsers: 3,
	})
	if err != nil {
		st := status.Convert(err)
		fmt.Println(st.Code(), st.Message(), st.Details())
	} else {
		for {
			user, err := srv.Recv()
			if err != nil {
				if err != io.EOF {
					st := status.Convert(err)
					fmt.Println(st.Code(), st.Message(), st.Details())
				}
				break
			}

			fmt.Println(user.GetId())
		}
	}

	fmt.Println("finished")
}
