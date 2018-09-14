ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

build:
	make -C ../c build
	CGO_LDFLAGS="$(ROOT_DIR)../c/libsputnikvm.a -ldl" go build -o build/sample  ./main
	
run: build
	./build/sample
