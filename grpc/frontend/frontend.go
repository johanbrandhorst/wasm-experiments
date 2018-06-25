package main

import (
	"context"
	"fmt"
	"io"
	"os"

	_ "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"

	web "github.com/johanbrandhorst/wasm-experiments/grpc/proto"
)

// Build with Go WASM fork
//go:generate rm -f ./html/test.wasm
//go:generate bash -c "GOOS=js GOARCH=wasm GOROOT=$GOPATH/src/github.com/johanbrandhorst/go/ $GOPATH/src/github.com/johanbrandhorst/go/bin/go build -o ./html/test.wasm frontend.go"

// Integrate generated JS into a Go file for static loading.
//go:generate bash -c "go run assets_generate.go"

func init() {
	// Should only be done from init functions
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout))
}

func main() {
	cc, err := grpc.Dial("")
	if err != nil {
		fmt.Println(err)
		return
	}
	client := web.NewBackendClient(cc)
	resp, err := client.GetUser(context.Background(), &web.GetUserRequest{
		UserId: "1234",
	})
	if err != nil {
		st := status.Convert(err)
		fmt.Println(st.Code(), st.Message(), st.Details())
	} else {
		fmt.Println(resp.GetId())
	}
	resp, err = client.GetUser(context.Background(), &web.GetUserRequest{
		UserId: "123",
	})
	if err != nil {
		st := status.Convert(err)
		fmt.Println(st.Code(), st.Message(), st.Details())
	} else {
		fmt.Println(resp.GetId())
	}

	srv, err := client.GetUsers(context.Background(), &web.GetUsersRequest{
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
