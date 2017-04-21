test:
	# Run test in current dir and all subdirectories
	go test ./...

test-race:
	# Run test in current dir and all subdirectories
	go test -race -bench . ./...

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
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/fzipp/gocyclo
	#go get -u github.com/barakmich/go-nyet
	#go get -u github.com/golang/lint/golint
	go get github.com/opennota/check/cmd/defercheck
	go get github.com/opennota/check/cmd/structcheck
	go get github.com/opennota/check/cmd/varcheck

cover:
	go test -cover ./...

# Coverage for a particular package
# go test -coverprofile=coverage.out ./...

# Show package coverage in web browser
# go tool cover -html=coverage.out

quality: mccabe  defercheck structcheck varcheck

mccabe:
	gocyclo -over 9 .

nyet:
	go-nyet ./...

defercheck:
	defercheck ./...

structcheck:
	structcheck ./...

varchack:
	vearcheck ./...

clean:
	go clean ./...
	git clean -f -d
