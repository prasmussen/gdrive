#!/bin/bash
set -ex
go test ./...
if [ "$TRAVIS_GO_VERSION" = "1.6" ]; then
 GOARCH=386 go test ./...
 go tool vet -example .
 go get github.com/client9/misspell/cmd/misspell
 go get github.com/soniakeys/vetc
 misspell -error * */* */*/*
 vetc
fi
