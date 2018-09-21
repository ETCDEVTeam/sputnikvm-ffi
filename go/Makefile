.PHONY : build run clean

ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

build:
	make -C ../c build
	rm -rf build/
	CGO_LDFLAGS="$(ROOT_DIR)../c/libsputnikvm.a -ldl" go build -o build/sample  ./main

run: build
	./build/sample

clean: 
	$(shell rm -fr build)