#!/bin/bash

go get -u golang.org/x/tools/cmd/cover
go get -u github.com/fzipp/gocyclo
go get -u github.com/barakmich/go-nyet
#go get -u github.com/golang/lint/golint
#go get github.com/opennota/check/cmd/defercheck
go get github.com/opennota/check/cmd/structcheck
go get github.com/opennota/check/cmd/varcheck