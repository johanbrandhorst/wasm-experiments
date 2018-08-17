// +build js,wasm

package main

import "syscall/js"

var document js.Value

func init() {
	document = js.Global().Get("document")
}

func main() {
	div := document.Get("body")

	node := document.Call("createElement", "div")
	node.Set("innerText", "Hello jsgo.io!")

	div.Call("appendChild", node)
}
