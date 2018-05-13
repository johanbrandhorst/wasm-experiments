# Frontend

This folder contains the entire source for the frontend app hosted by the server.

## bundle

The `bundle` package is a `vfsgen` generated package, created from the contents of
the `html` folder. It serves as the interface that the `main.go` server uses to serve
the GopherJS frontend without the need for a `static` directory. The generation is done
via `go:generate` in `frontend.go`.

## html

The `html` folder contains the static sources used.
