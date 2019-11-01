#!/bin/sh

set -ex

go get github.com/golangci/golangci-lint/cmd/golangci-lint
go get -d ./...
