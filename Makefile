.PHONY: hello
hello: clean
	GOOS=js GOARCH=wasm go build -o ./html/test.wasm ./hello/main.go
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.html ./html/index.html
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js ./html/wasm_exec.js

.PHONY: tinygo
tinygo: clean
	docker run --rm -v $$(pwd):/go/src/github.com/johanbrandhorst/wasm-experiments tinygo/tinygo:0.6.1 /bin/bash -c "\
			cd /go/src/github.com/johanbrandhorst/wasm-experiments && \
			tinygo build -o ./html/test.wasm -target wasm --no-debug ./$(target)/main.go && \
			cp /usr/local/tinygo/targets/wasm_exec.js ./html/wasm_exec.js\
	"
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

.PHONY: canvas
canvas: clean
	GOOS=js GOARCH=wasm go build -o ./html/test.wasm ./canvas/main.go
	cp ./canvas/index.html ./html/index.html
	cp ./canvas/main.go ./html/main.go
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js ./html/wasm_exec.js

tinygo-canvas: clean
	#docker run --rm -v $$(pwd):/go/src/github.com/johanbrandhorst/wasm-experiments tinygo/tinygo:0.6.1 /bin/bash -c "\
	#		cd /go/src/github.com/johanbrandhorst/wasm-experiments && \
	#		tinygo build -o ./html/test.wasm -target wasm --no-debug ./canvas/main.go && \
	#		cp /usr/local/tinygo/targets/wasm_exec.js ./html/wasm_exec.js\
	#"
	# Requires tinygo with experimental GC
	# https://github.com/tinygo-org/tinygo/pull/350
	tinygo build -o ./html/test.wasm -target wasm --no-debug ./tinygo-canvas/main.go
	cp ~/go/src/github.com/tinygo-org/tinygo/targets/wasm_exec.js ./html/
	cp ./tinygo-canvas/index.html ./html/index.html
	cp ./tinygo-canvas/main.go ./html/main.go

.PHONY: ebiten
ebiten: clean
	GO111MODULE=on GOOS=js GOARCH=wasm go build -o ./html/ebiten.wasm ./ebiten/main.go
	cp ./ebiten/index.html ./html/index.html
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js ./html/wasm_exec.js

.PHONY: vugu
vugu: clean
	GO111MODULE=on go get github.com/vugu/vugu/cmd/vugugen
	vugugen --skip-go-mod --skip-main ./vugu/
	GOOS=js GOARCH=wasm go build -o ./html/main.wasm ./vugu/
	cp ./vugu/index.html ./html/
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js ./html/wasm_exec.js

.PHONY: vecty
vecty: clean
	GOOS=js GOARCH=wasm go build -o ./html/test.wasm ./vecty/main.go
	cp ./vecty/index.html ./html/index.html
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js ./html/wasm_exec.js

test: clean
	GOOS=js GOARCH=wasm go test -c -o ./html/test.wasm ./test/

clean:
	rm -f ./html/*

install-test:
	go get github.com/agnivade/wasmbrowsertest
	mv $$GOPATH/bin/wasmbrowsertest $$GOPATH/bin/go_js_wasm_exec
