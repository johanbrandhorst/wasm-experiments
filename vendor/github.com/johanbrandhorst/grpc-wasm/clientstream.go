// +build js,wasm

package grpc

import (
	"context"

	"google.golang.org/grpc"
)

// ClientStream defines the interface a client stream has to satisfy.
type ClientStream grpc.ClientStream

// ServerStream defines the interface a server stream has to satisfy.
type ServerStream grpc.ServerStream

// StreamHandler defines the handler called by gRPC server to complete the
// execution of a streaming RPC. If a StreamHandler returns an error, it
// should be produced by the status package, or else gRPC will use
// codes.Unknown as the status code and err.Error() as the status message
// of the RPC.
type StreamHandler func(srv interface{}, stream ServerStream) error

// StreamDesc represents a streaming RPC service's method specification.
type StreamDesc struct {
	StreamName string
	Handler    StreamHandler

	// At least one of these is true.
	ServerStreams bool
	ClientStreams bool
}

type methodHandler func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor UnaryServerInterceptor) (interface{}, error)

// MethodDesc represents an RPC service's method specification.
type MethodDesc struct {
	MethodName string
	Handler    methodHandler
}

// ServiceDesc represents an RPC service's specification.
type ServiceDesc struct {
	ServiceName string
	// The pointer to the service interface. Used to check whether the user
	// provided implementation satisfies the interface requirements.
	HandlerType interface{}
	Methods     []MethodDesc
	Streams     []StreamDesc
	Metadata    interface{}
}
