package backend

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
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
		st := status.New(codes.InvalidArgument, "invalid id")
		detSt, err := st.WithDetails(&errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "user",
					Description: "That user does not exist",
				},
			},
		})
		if err == nil {
			return nil, detSt.Err()
		}
		return nil, st.Err()
	}
	return &server.User{
		Id: req.GetUserId(),
	}, nil
}

func (b Backend) GetUsers(req *server.GetUsersRequest, srv server.Backend_GetUsersServer) error {
	for index := 0; index < int(req.GetNumUsers()); index++ {
		err := srv.Send(&server.User{
			Id: strconv.Itoa(index),
		})
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
	}

	return nil
}
