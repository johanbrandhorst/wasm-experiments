.PHONY: hello
hello:
	rm -f ./html/*
	GOOS=js GOARCH=wasm go build -o ./html/test.wasm ./hello/main.go
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.html ./html/index.html
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js ./html/wasm_exec.js

.PHONY: tinygo
tinygo:
	rm -f ./html/*
	docker run --rm -v $$(pwd):/go/src/github.com/johanbrandhorst/wasm-experiments --entrypoint /bin/bash tinygo/tinygo:latest -c "\
		cd /go/src/github.com/johanbrandhorst/wasm-experiments && \
		tinygo build -o ./html/test.wasm -target wasm ./$(target)/main.go && \
		cp /go/src/github.com/tinygo-org/tinygo/targets/wasm_exec.js ./html/wasm_exec.js\
	"
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.html ./html/index.html

.PHONY: channels
channels:
	rm -f ./html/*
	GOOS=js GOARCH=wasm go build -o ./html/test.wasm ./channels/main.go
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.html ./html/index.html
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js ./html/wasm_exec.js

.PHONY: js
js:
	rm -f ./html/*
	GOOS=js GOARCH=wasm go build -o ./html/test.wasm ./js/main.go
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.html ./html/index.html
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js ./html/wasm_exec.js
	sed -i -e 's;</button>;</button>\n\t<div id=\"target\"></div>;' ./html/index.html

.PHONY: fetch
fetch:
	rm -f ./html/*
	GOOS=js GOARCH=wasm go build -o ./html/test.wasm ./fetch/main.go
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.html ./html/index.html
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js ./html/wasm_exec.js
	sed -i -e 's;</button>;</button>\n\t<div id=\"target\"></div>;' ./html/index.html

serve:
	go run main.go
