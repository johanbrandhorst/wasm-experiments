// +build js,wasm

package grpc

// DialOption configures how we set up the connection.
type DialOption func(*dialOptions)

// dialOptions configure a Dial call. dialOptions are set by the DialOption
// values passed to Dial.
type dialOptions struct {
}
