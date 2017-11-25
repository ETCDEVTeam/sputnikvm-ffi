build:
	cd ffi && cargo build --release
	cp ffi/target/release/libsputnikvm_ffi.so libsputnikvm.so

doxygen:
	doxygen

deploy-docs: doxygen
	git add -f ./html
	$(eval TREE_OBJ_ID := $(shell git write-tree --prefix=c/html))
	git reset -- ./html
	$(eval COMMIT_ID := $(shell git commit-tree -p gh-pages -m "Doc deployments" $(TREE_OBJ_ID)))
	git update-ref refs/heads/gh-pages $(COMMIT_ID)
