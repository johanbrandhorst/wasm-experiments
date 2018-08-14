// +build js,wasm

package main

import "syscall/js"

var document js.Value

func init() {
	document = js.Global().Get("document")
}

func main() {
	div := document.Call("getElementById", "target")

	node := document.Call("createElement", "div")
	node.Set("innerHTML", "Hello World")

	div.Call("appendChild", node)
}
