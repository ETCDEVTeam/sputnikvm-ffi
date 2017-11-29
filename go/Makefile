ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

build:
	make -C ../c build
	go build main.go

run: build
	./main
