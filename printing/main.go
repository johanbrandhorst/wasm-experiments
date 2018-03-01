package main

import (
	"fmt"
	"runtime/js"
)

func main() {
	fmt.Println("Test")
	js.Global.Get("console").Call("log", "yadayada")
}
