FROM golang:1.6
MAINTAINER bernhard.biskup@gmx.de

WORKDIR /go/src/github.com/bbiskup/edify

ADD scripts scripts
RUN ./scripts/get_test_deps.sh
COPY . /go/src/github.com/bbiskup/edify
RUN go get -t ./...
RUN go build -v

ENTRYPOINT ["/bin/bash", "-c"]