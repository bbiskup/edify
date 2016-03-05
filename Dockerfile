FROM golang:1.6
MAINTAINER bernhard.biskup@gmx.de

WORKDIR /go/src/github.com/bbiskup/edify

#RUN useradd -m dev

#RUN apt-get update -yy && apt-get install -yy wget
#RUN wget -v https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz
#RUN (cd /usr/local/ && tar xzvf go1.6.linux-amd64.tar.gz)

ADD scripts scripts
RUN ./scripts/get_test_deps.sh
COPY . /go/src/github.com/bbiskup/edify
RUN go get -t ./...
RUN go build -v

ENTRYPOINT ["/bin/bash", "-c"]