build: clean compile

clean:
	rm -f ./html/test.wasm

compile:
	GOOS=js GOARCH=wasm GOROOT=$(GOPATH)/src/github.com/neelance/go/ $(GOPATH)/src/github.com/neelance/go/bin/go build -o ./html/test.wasm main.go
