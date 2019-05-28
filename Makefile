.PHONY: hello
hello: clean
	GOOS=js GOARCH=wasm go build -o ./html/test.wasm ./hello/main.go
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.html ./html/index.html
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js ./html/wasm_exec.js

.PHONY: tinygo
tinygo: clean
	#docker run --rm -v $$(pwd):/go/src/github.com/johanbrandhorst/wasm-experiments tinygo/tinygo:0.5.0 /bin/bash -c "\
	#	cd /go/src/github.com/johanbrandhorst/wasm-experiments && \
	#	tinygo build -o ./html/test.wasm -target wasm --no-debug ./$(target)/main.go && \
	#	cp /usr/local/tinygo/targets/wasm_exec.js ./html/wasm_exec.js\
	#"
	tinygo build -o ./html/test.wasm -target wasm --no-debug ./$(target)/main.go
	cp ~/go/src/github.com/tinygo-org/tinygo/targets/wasm_exec.js ./html/wasm_exec.js
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.html ./html/index.html
	sed -i -e 's;</button>;</button>\n\t<div id=\"target\"></div>;' ./html/index.html

.PHONY: channels
channels: clean
	GOOS=js GOARCH=wasm go build -o ./html/test.wasm ./channels/main.go
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.html ./html/index.html
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js ./html/wasm_exec.js

.PHONY: js
js: clean
	GOOS=js GOARCH=wasm go build -o ./html/test.wasm ./js/main.go
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.html ./html/index.html
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js ./html/wasm_exec.js
	sed -i -e 's;</button>;</button>\n\t<div id=\"target\"></div>;' ./html/index.html

.PHONY: fetch
fetch: clean
	GOOS=js GOARCH=wasm go build -o ./html/test.wasm ./fetch/main.go
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.html ./html/index.html
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js ./html/wasm_exec.js
	sed -i -e 's;</button>;</button>\n\t<div id=\"target\"></div>;' ./html/index.html

test: clean
	GOOS=js GOARCH=wasm go test -c -o ./html/test.wasm ./test/

clean:
	rm -f ./html/*

serve:
	go run main.go

install-test:
	go get github.com/agnivade/wasmbrowsertest
	mv $$GOPATH/bin/wasmbrowsertest $$GOPATH/bin/go_js_wasm_exec
