generate:
	protoc -I. ./proto/web.proto \
		--go_out=plugins=grpc:$$GOPATH/src \
		--wasm_out=plugins=grpc:./
	mv ./proto/web.wasm.pb.go ./proto/client/
	go generate -x ./frontend/

serve:
	go run main.go
