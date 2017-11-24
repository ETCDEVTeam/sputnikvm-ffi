# SputnikVM FFI Bindings

This repo contains C and Go bindings for the SputnikVM core
library.

## Usage

### C

In `c` folder, run `make build`. It will generate an object file
`libsputnikvm.so`, and you can use the header file `sputnikvm.h` to
interact with it.

### Go

See `go/main.go` for example usage, note that when building a Go
application, you need to add the `-ldflags` option to the folder that
contains the `libsputnikvm.so` file.
