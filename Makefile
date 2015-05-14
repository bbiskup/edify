test:
	# Run test in current dir and all subdirectories
	go test ./...

test-multi-cpu:
	go test -cpu=1,2,3,4,5,6,7,8 ./...

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

# Coverage for a particular package
# go test -coverprofile=coverage.out ./...

# Show package coverage in web browser
# go tool cover -html=coverage.out