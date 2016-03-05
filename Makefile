all: build up

build:
	docker-compose build

up:
	docker-compose up -d

bash:
	./docker-cmd.sh bash

test:
	# Run test in current dir and all subdirectories
	./docker-cmd.sh go test ./...

test-race:
	# Run test in current dir and all subdirectories
	./docker-cmd.sh go test -race -bench . ./...

test-multi-cpu:
	./docker-cmd.sh go test -cpu=1,2,3,4,5,6,7,8 ./...

test-verbose:
	./docker-cmd.sh go test -v ./...


bench:
	# Run test in current dir and all subdirectories
	./docker-cmd.sh go test -bench . ./...

check:
	./docker-cmd.sh go vet -x ./...

cover:
	./docker-cmd.sh go test -cover ./...

# Coverage for a particular package
# ./docker-cmd.sh go test -coverprofile=coverage.out ./...

# Show package coverage in web browser
# ./docker-cmd.sh go tool cover -html=coverage.out

quality: mccabe nyet defercheck structcheck varcheck

mccabe:
	./docker-cmd.sh gocyclo -over 9 .

nyet:
	./docker-cmd.sh go-nyet ./...

defercheck:
	./docker-cmd.sh defercheck ./...

structcheck:
	./docker-cmd.sh structcheck ./...

varchack:
	./docker-cmd.sh varcheck ./...

clean:
	./docker-cmd.sh go clean ./...
	git clean -f -d