ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

build:
	make -C ../c build
	go build -ldflags="-r $(ROOT_DIR)../c" main.go

run: build
	./main
