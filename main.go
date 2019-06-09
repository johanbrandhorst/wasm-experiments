package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./html"))
	log.Print("Serving on http://localhost:8080")
	err := http.ListenAndServe(":8080", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/test.wasm" {
			resp.Header().Set("content-type", "application/wasm")
		}

		fs.ServeHTTP(resp, req)
	}))
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
