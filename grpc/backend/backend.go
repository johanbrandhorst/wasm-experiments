package backend

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/johanbrandhorst/wasm-experiments/grpc/proto/server"
)

// Backend should be used to implement the server interface
// exposed by the generated server proto.
type Backend struct {
}

// Ensure struct implements interface
var _ server.BackendServer = (*Backend)(nil)

func (b Backend) GetUser(ctx context.Context, req *server.GetUserRequest) (*server.User, error) {
	if req.GetUserId() != "1234" {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	return &server.User{
		Id: req.GetUserId(),
	}, nil
}
