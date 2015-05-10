test:
	# Run test in current dir and all subdirectories
	go test ./...

test-verbose:
	go test -v ./...


bench:
	# Run test in current dir and all subdirectories
	go test -bench . ./...

check:
	go vet -x ./...

get-deps:
	go get -t ./...

get-test-deps:
	go get golang.org/x/tools/cmd/cover

cover:
	go test -cover ./...