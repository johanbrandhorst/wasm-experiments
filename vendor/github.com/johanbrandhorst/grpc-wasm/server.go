// +build js,wasm

package grpc

// Server is a gRPC server to serve RPC requests.
type Server struct{}

func (s *Server) RegisterService(sd *ServiceDesc, ss interface{}) {}
