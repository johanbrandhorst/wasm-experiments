// +build js,wasm

package main

import (
	"syscall/js"
	"syscall/js/callback"
)

func main() {
	wait := make(chan struct{})
	promise := js.Global.Call("fetch", "https://api.github.com")
	promise.Call("then", callback.New(func(args []js.Value) {
		response := args[0]
		r := response.Get("body").Call("getReader")
		r.Call("read")
		close(wait)
	}).Value)
	<-wait
}
