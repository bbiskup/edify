test:
	# Run test in current dir and all subdirectories
	go test ./...

test-verbose:
	go test -v ./...


bench:
	# Run test in current dir and all subdirectories
	go test -bench . ./...

check:
	go vet