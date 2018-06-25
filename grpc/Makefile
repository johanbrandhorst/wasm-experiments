generate:
	protoc -I. ./proto/web.proto \
		--go_out=plugins=grpc:$$GOPATH/src
	go generate -x ./frontend/

to_frontend:
	cd ../../../../google.golang.org/grpc && \
		mv -f parse_js.gox parse_js.go && \
		mv -f parse.go parse.gox && \
		mv -f newstream_js.gox newstream_js.go && \
		mv -f newstream.go newstream.gox && \
		mv -f dial_js.gox dial_js.go && \
		mv -f dial.go dial.gox && \
		mv -f transport/stream_js.gox transport/stream_js.go && \
		mv -f transport/stream.go transport/stream.gox && \
		mv -f transport/client_transport_js.gox transport/client_transport_js.go && \
		mv -f transport/client_transport.go transport/client_transport.gox && \
		mv -f transport/content_type_js.gox transport/content_type_js.go && \
		mv -f transport/content_type.go transport/content_type.gox

to_backend:
	cd ../../../../google.golang.org/grpc && \
		mv -f parse_js.go parse_js.gox && \
		mv -f parse.gox parse.go && \
		mv -f newstream_js.go newstream_js.gox && \
		mv -f newstream.gox newstream.go && \
		mv -f dial_js.go dial_js.gox && \
		mv -f dial.gox dial.go && \
		mv -f transport/stream_js.go transport/stream_js.gox && \
		mv -f transport/stream.gox transport/stream.go && \
		mv -f transport/client_transport_js.go transport/client_transport_js.gox && \
		mv -f transport/client_transport.gox transport/client_transport.go && \
		mv -f transport/content_type_js.go transport/content_type_js.gox && \
		mv -f transport/content_type.gox transport/content_type.go

serve:
	go run main.go
