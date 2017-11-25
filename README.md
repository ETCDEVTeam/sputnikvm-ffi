# SputnikVM FFI Bindings

This repo contains C and Go bindings for the SputnikVM core
library.

## Usage

### C

In `c` folder, run `make build`. It will generate an object file
`libsputnikvm.so`, and you can use the header file `sputnikvm.h` to
interact with it. You can find the generated documentation file for
`sputnikvm.h`
[here](https://ethereumproject.github.io/sputnikvm-ffi/sputnikvm_8h.html).

### Go

Go binding is currently work-in-progress. See `go/main.go` for example
usage, note that when building a Go application, you need to add the
`-ldflags` option to the folder that contains the `libsputnikvm.so`
file.
