# wasm-experiments

Just some playing around with the experimental Go WebAssembly arch target.

## Requirements

Requires `go` >= 1.11.0.
[TinyGo](https://github.com/tinygo-org/tinygo) examples require `docker`.

## Basic instructions

Choose your target from the experiments. Compile with `make <target>` and serve
with `make serve`. This starts a local web server which serves the `./html`
directory. It ensures `.wasm` files are served with their appropriate
content-type.

For example:

```bash
$ make hello
rm -f ./html/*
GOOS=js GOARCH=wasm go build -o ./html/test.wasm ./hello/main.go
cp $(go env GOROOT)/misc/wasm/wasm_exec.html ./html/index.html
cp $(go env GOROOT)/misc/wasm/wasm_exec.js ./html/wasm_exec.js
$ make serve
go run main.go
2019/02/24 14:11:51 Serving on http://localhost:8080
```

Navigate to http://localhost:8080 to load the page. Some examples require opening
the browser console to be seen.

## TinyGo targets

To compile something with the TinyGo WebAssembly compiler, simply choose the
target and invoke the `tinygo` make rule with the target specified and
serve as usual. For example:

```bash
$ make tinygo target=hello
docker run --rm -v $$(pwd):/go/src/github.com/johanbrandhorst/wasm-experiments tinygo/tinygo:0.5.0 /bin/bash -c "\
        cd /go/src/github.com/johanbrandhorst/wasm-experiments && \
        tinygo build -o ./html/test.wasm -target wasm --no-debug ./hello/main.go && \
        cp /usr/local/tinygo/targets/wasm_exec.js ./html/wasm_exec.js\
"
cp $$(go env GOROOT)/misc/wasm/wasm_exec.html ./html/index.html
sed -i -e 's;</button>;</button>\n\t<div id=\"target\"></div>;' ./html/index.html
$ make serve
go run main.go
2019/02/24 14:33:58 Serving on http://localhost:8080
```

Note that some of the targets either do not compile or panic at runtime when
compiled with TinyGo. The following targets have been tested to work with
TinyGo:

- `hello`
- `channels`
- `js`

## Experiments

### Hello

A simple `Hello World` example that prints to the browser console.

### Channels

Showing basic first class support for channel operations in compiled WebAssembly.

### JS

A simple example of how to interact with the JavaScript world from WebAssembly.

### Fetch

A more complicated example showing how to use `net/http` `DefaultClient` to
send a HTTP request, parse the result and write it to the DOM.

