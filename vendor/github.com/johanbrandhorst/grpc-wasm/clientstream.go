// +build js,wasm

package grpc

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

type clientStream struct {
	ctx    context.Context
	req    *http.Request
	client *http.Client
	errCh  chan error
	msgCh  chan []byte
}

func newStream(ctx context.Context, client *http.Client, endpoint string) (*clientStream, error) {
	cs := &clientStream{
		ctx:    ctx,
		client: client,
	}

	req, err := http.NewRequest(
		"POST",
		endpoint,
		nil,
	)
	if err != nil {
		return nil, status.New(codes.Unavailable, err.Error()).Err()
	}

	cs.req = req.WithContext(ctx)
	return cs, nil
}

func (c *clientStream) Header() (metadata.MD, error) {
	return nil, nil
}

func (c *clientStream) Trailer() metadata.MD {
	return nil
}

func (c *clientStream) Context() context.Context {
	return c.ctx
}

func (c *clientStream) RecvMsg(reply interface{}) error {
	select {
	case <-c.ctx.Done():
		return c.ctx.Err()
	case err := <-c.errCh:
		return err
	case msg, ok := <-c.msgCh:
		if !ok {
			return io.EOF
		}
		err := proto.Unmarshal(msg, reply.(proto.Message))
		return err
	}
}

func (c *clientStream) SendMsg(req interface{}) error {
	msg, err := proto.Marshal(req.(proto.Message))
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	bufHeader := make([]byte, 5)

	// Write length of b into buf
	binary.BigEndian.PutUint32(bufHeader[1:], uint32(len(msg)))

	c.req.Body = ioutil.NopCloser(bytes.NewBuffer(append(bufHeader, msg...)))
	addHeaders(c.req)

	resp, err := c.client.Do(c.req)
	if err != nil {
		return status.Error(codes.Unavailable, err.Error())
	}

	st := statusFromHeaders(resp.Header)
	if st.Code() != codes.OK {
		resp.Body.Close()
		return st.Err()
	}

	c.errCh = make(chan error, 1)
	c.msgCh = make(chan []byte, 1)

	// Read response asynchronously
	go func() {
		defer resp.Body.Close()

		msgHeader := make([]byte, 5)
		for {
			_, err := io.ReadFull(resp.Body, msgHeader)
			if err != nil {
				c.errCh <- status.Error(codes.Internal, err.Error())
				return
			}
			// 1 in MSB signifies that this is the trailer. Break loop.
			// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-WEB.md#protocol-differences-vs-grpc-over-http2
			if msgHeader[0]>>7 == 1 {
				break
			}

			msgLen := binary.BigEndian.Uint32(msgHeader[1:])

			msg := make([]byte, msgLen)
			_, err = io.ReadFull(resp.Body, msg)
			if err != nil {
				c.errCh <- status.Error(codes.Internal, err.Error())
				return
			}
			c.msgCh <- msg
		}

		if msgHeader[0]&1 == 0 {
			trailers, err := readTrailers(resp.Body)
			if err != nil {
				c.errCh <- status.Error(codes.Internal, err.Error())
				return
			}
			st = statusFromHeaders(trailers)
		} else {
			// TODO(johanbrandhorst): Support compressed trailers
		}

		if st.Code() != codes.OK {
			c.errCh <- st.Err()
			return
		}

		close(c.msgCh)
	}()

	return nil
}

func (c *clientStream) CloseSend() error {
	return nil
}
