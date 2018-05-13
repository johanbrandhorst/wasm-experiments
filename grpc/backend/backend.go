package backend

import (
	"context"

	"github.com/johanbrandhorst/wasm-experiments/grpc/proto/server"
)

// Backend should be used to implement the server interface
// exposed by the generated server proto.
type Backend struct {
}

// Ensure struct implements interface
var _ server.BackendServer = (*Backend)(nil)

func (b Backend) GetUser(ctx context.Context, req *server.GetUserRequest) (*server.User, error) {
	return &server.User{
		Id: "1234",
	}, nil
}
