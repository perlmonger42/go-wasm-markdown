# Go WebAssembly Markdown Example

This project reproduces the Markdown example in https://github.com/hexops/vecty/tree/main/example
using Vanilla Go WebAssembly (no Vecty, nor any other framework).


# Building

## Building the Markdown App

Simply run:
```bash
GOOS=js GOARCH=wasm go build -o cmd/markdown/main.wasm ./cmd/markdown
```

This creates `./cmd/markdown/main.wasm`, which is the Markdown App
compiled to webassembly. The provided `index.html` loads and executes
that code when is loaded itself.

## Development Server

This is the [basic Go web server](https://go.dev/play/p/pZ1f5pICVbV), enhanced in two ways:

1. It monitors the the server's `-dir` folder for changes, and recompiles
  the markdown app on any change.
2. It provides a `/ws` websocket-connection endpoint that pages can use to connect back
  to the server. The server sends "reload" when the watcher has detected a change and
  recompiled. That allows the app to automatically refresh on change.

Build the server into `bin/server` using:
```bash
    go build -o bin/server ./cmd/server.go
```

Then run it using `./bin/server`. Once it's running you can load the web page at
[http://localhost:8080](http://localhost:8080).

# References
[Go Wiki: WebAssembly](https://go.dev/wiki/WebAssembly#interacting-with-the-dom)
[Package syscall/js](https://pkg.go.dev/syscall/js)
