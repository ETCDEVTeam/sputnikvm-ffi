ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

build:
	cd lib/sputnikvm-go && cargo build --release
	cp lib/sputnikvm-go/target/release/libsputnikvm_go.so lib/
	go build -ldflags="-r $(ROOT_DIR)lib" main.go

run: build
	./main
