# Fetch
The Go http.Transport interface implemented over the WHATWG Fetch API using the WebAssembly arch target.

## Usage

This package requires the Go WASM compilation target to be supported.
Short example:

```go
c := http.Client{
    Transport: &fetch.Transport{},
}
resp, err := c.Get("https://api.github.com")
if err != nil {
    fmt.Println(err)
    return
}
defer resp.Body.Close()
b, err := ioutil.ReadAll(resp.Body)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(string(b))
```

See [my wasm-experiments repo](https://github.com/johanbrandhorst/wasm-experiments)
for the full example of its use.

## Attribution

The code is largely based on the Fetch API implementation
[in GopherJS](https://github.com/gopherjs/gopherjs/blob/8dffc02ea1cb8398bb73f30424697c60fcf8d4c5/compiler/natives/src/net/http/fetch.go).
