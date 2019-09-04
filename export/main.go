package main

import (
	"syscall/js"
)

var document js.Value

func init() {
	document = js.Global().Get("document")
}

func main() {
	document.Set("appendText", js.FuncOf(appendText))

	// Prevent main from exiting
	select {}
}

func appendText(_ js.Value, _ []js.Value) interface{} {
	msg := document.Call("getElementById", "input").Get("value").String()

	p := document.Call("createElement", "p")
	p.Set("innerHTML", msg)
	document.Call("getElementById", "target").Call("appendChild", p)
	return nil
}
