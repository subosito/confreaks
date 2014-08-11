fmt:
	go fmt
	pushd bin; go fmt; popd

test:
	go test -v

lint:
	golint .
