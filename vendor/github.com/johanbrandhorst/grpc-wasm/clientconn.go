// +build js,wasm

package grpc

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/johanbrandhorst/fetch"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ClientConn struct {
	// TODO(johanbrandhorst): Remove once we can rely on http.DefaultClient
	client *http.Client

	target string
}

// Dial creates a client connection to the target. The target string should
// be a URL with scheme HTTP or HTTPS, or a FQDN to infer the scheme.
func Dial(target string, opts ...DialOption) (*ClientConn, error) {
	return DialContext(context.Background(), target, opts...)
}

func DialContext(ctx context.Context, target string, opts ...DialOption) (conn *ClientConn, err error) {
	return &ClientConn{
		client: &http.Client{
			Transport: &fetch.Transport{},
		},
		target: target,
	}, nil
}

func (cc *ClientConn) NewStream(ctx context.Context, desc *StreamDesc, method string, opts ...CallOption) (ClientStream, error) {
	if desc.ClientStreams {
		return nil, status.Error(codes.Unimplemented, "client-side streaming is not supported by grpc-web")
	}

	endpoint := cc.target + "/" + method
	if cc.target == "" {
		endpoint = method
	}

	return newStream(ctx, cc.client, endpoint)
}

func (cc *ClientConn) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...CallOption) error {
	b, err := proto.Marshal(args.(proto.Message))
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	bufHeader := make([]byte, 5)

	// Write length of b into buf
	binary.BigEndian.PutUint32(bufHeader[1:], uint32(len(b)))

	endpoint := cc.target + "/" + method
	if cc.target == "" {
		endpoint = method
	}
	req, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer(append(bufHeader, b...)),
	)
	if err != nil {
		return status.Error(codes.Unavailable, err.Error())
	}
	req = req.WithContext(ctx)
	addHeaders(req)

	resp, err := cc.client.Do(req)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer resp.Body.Close()

	st := statusFromHeaders(resp.Header)
	if st.Code() != codes.OK {
		return st.Err()
	}

	msgHeader := make([]byte, 5)
	for {
		_, err := resp.Body.Read(msgHeader)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		// 1 in MSB signifies that this is the trailer. Break loop.
		// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-WEB.md#protocol-differences-vs-grpc-over-http2
		if msgHeader[0]>>7 == 1 {
			break
		}

		msgLen := binary.BigEndian.Uint32(msgHeader[1:])

		msg := make([]byte, msgLen)
		_, err = resp.Body.Read(msg)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		err = proto.Unmarshal(msg, reply.(proto.Message))
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}

	if msgHeader[0]&1 == 0 {
		trailers, err := readTrailers(resp.Body)
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
		st = statusFromHeaders(trailers)
	} else {
		// TODO(johanbrandhorst): Support compressed trailers
	}

	return st.Err()
}

func addHeaders(req *http.Request) {
	// TODO: Add more headers
	// https://github.com/grpc/grpc-go/blob/590da37e2dfb4705d8ebd9574ce4cb75295d9674/transport/http2_client.go#L356
	req.Header.Add("content-type", "application/grpc-web+proto")
	req.Header.Add("x-grpc-web", "1")
	if dl, ok := req.Context().Deadline(); ok {
		timeout := dl.Sub(time.Now())
		req.Header.Add("grpc-timeout", encodeTimeout(timeout))
	}
	md, ok := metadata.FromOutgoingContext(req.Context())
	if ok {
		for h, vs := range md {
			for _, v := range vs {
				req.Header.Add(h, v)
			}
		}
	}
}

const maxTimeoutValue int64 = 100000000 - 1

// Copied from grpc-go
// https://github.com/grpc/grpc-go/blob/590da37e2dfb4705d8ebd9574ce4cb75295d9674/transport/http_util.go#L388
// div does integer division and round-up the result. Note that this is
// equivalent to (d+r-1)/r but has less chance to overflow.
func div(d, r time.Duration) int64 {
	if m := d % r; m > 0 {
		return int64(d/r + 1)
	}
	return int64(d / r)
}

// Copied from grpc-go
// https://github.com/grpc/grpc-go/blob/590da37e2dfb4705d8ebd9574ce4cb75295d9674/transport/http_util.go#L398
func encodeTimeout(t time.Duration) string {
	if t <= 0 {
		return "0n"
	}
	if d := div(t, time.Nanosecond); d <= maxTimeoutValue {
		return strconv.FormatInt(d, 10) + "n"
	}
	if d := div(t, time.Microsecond); d <= maxTimeoutValue {
		return strconv.FormatInt(d, 10) + "u"
	}
	if d := div(t, time.Millisecond); d <= maxTimeoutValue {
		return strconv.FormatInt(d, 10) + "m"
	}
	if d := div(t, time.Second); d <= maxTimeoutValue {
		return strconv.FormatInt(d, 10) + "S"
	}
	if d := div(t, time.Minute); d <= maxTimeoutValue {
		return strconv.FormatInt(d, 10) + "M"
	}
	// Note that maxTimeoutValue * time.Hour > MaxInt64.
	return strconv.FormatInt(div(t, time.Hour), 10) + "H"
}

// Copied from grpc-go
// https://github.com/grpc/grpc-go/blob/b94ea975f3beb73799fac17cc24ee923fcd3cb5c/transport/http_util.go#L213
func decodeBinHeader(v string) ([]byte, error) {
	if len(v)%4 == 0 {
		// Input was padded, or padding was not necessary.
		return base64.StdEncoding.DecodeString(v)
	}
	return base64.RawStdEncoding.DecodeString(v)
}

func readTrailers(in io.Reader) (http.Header, error) {
	s := bufio.NewScanner(in)
	trailers := http.Header{}
	for s.Scan() {
		v := s.Text()
		kv := strings.SplitN(v, ": ", 2)
		if len(kv) != 2 {
			return nil, errors.New("malformed header: " + v)
		}
		trailers.Add(kv[0], kv[1])
	}

	return trailers, s.Err()
}

func statusFromHeaders(h http.Header) *status.Status {
	details := h.Get("grpc-status-details-bin")
	if details != "" {
		b, err := decodeBinHeader(details)
		if err != nil {
			return status.New(codes.Internal, "malformed grps-status-details-bin header: "+err.Error())
		}
		s := &spb.Status{}
		err = proto.Unmarshal(b, s)
		if err != nil {
			return status.New(codes.Internal, "malformed grps-status-details-bin header: "+err.Error())
		}
		return status.FromProto(s)
	}
	sh := h.Get("grpc-status")
	if sh != "" {
		val, err := strconv.Atoi(sh)
		if err != nil {
			return status.New(codes.Internal, "malformed grpc-status header: "+err.Error())
		}
		return status.New(codes.Code(val), h.Get("grpc-message"))
	}
	return status.New(codes.OK, "")
}
