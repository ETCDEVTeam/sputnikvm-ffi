ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

build:
	make -C ../c build
	CGO_LDFLAGS="$(ROOT_DIR)../c/libsputnikvm.a -ldl" go build -o build/sample

run: build
	./build/sample

debug:
	make -C ../c debug
	CGO_LDFLAGS="$(ROOT_DIR)../c/libsputnikvm.a -ldl" go build -o build/sample
	./build/sample