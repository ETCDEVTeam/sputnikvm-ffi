build:
	cd ffi && cargo build --release
	cp ffi/target/release/libsputnikvm_ffi.so libsputnikvm.so

doxygen:
	doxygen

deploy-docs:
	git add -f ./html
	$treeObjId = git write-tree --prefix=./html
	git reset -- ./html
	$commitId = git commit-tree -p gh-pages -m "Doc deployments" $treeObjId
	git update-ref refs/heads/gh-pages $commitId
