build: clean compile

clean:
	rm -f ../html/test.wasm

compile:
	GOOS=js GOARCH=wasm GOROOT=$(GOPATH)/src/github.com/johanbrandhorst/go/ $(GOPATH)/src/github.com/johanbrandhorst/go/bin/go build -o ../html/test.wasm main.go
