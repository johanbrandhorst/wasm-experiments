package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/NYTimes/gziphandler"
)

var gz = flag.Bool("gzip", false, "enable automatic gzip compression")

func main() {
	flag.Parse()
	h := wasmContentTypeSetter(http.FileServer(http.Dir("./html")))
	if *gz {
		h = gziphandler.GzipHandler(h)
	}

	log.Print("Serving on http://localhost:8080")
	err := http.ListenAndServe(":8080", h)
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func wasmContentTypeSetter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".wasm") {
			w.Header().Set("content-type", "application/wasm")
		}
		h.ServeHTTP(w, r)
	})
}
