// +build js,wasm

package main

import (
	"log"
	"os"

	"github.com/vugu/vugu"
)

func main() {
	rootInst, err := vugu.New(&Root{}, nil)
	if err != nil {
		log.Fatal(err)
	}

	env := vugu.NewJSEnv("#target", rootInst, vugu.RegisteredComponentTypes())
	env.DebugWriter = os.Stdout

	for ok := true; ok; ok = env.EventWait() {
		err = env.Render()
		if err != nil {
			log.Fatal(err)
		}
	}
}
