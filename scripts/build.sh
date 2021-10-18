#!/bin/bash
cd "$( dirname "$0"  )"
GOOS=windows GOARCH=arm64 go build -ldflags "-s -w" -i -o ../bin/github-dns ../cmd/github-dns/github-dns.go
GOOS=windows GOARCH=arm64 go build -ldflags "-s -w" -i -o ../bin/ghosts-cli ../cmd/ghosts-cli/ghosts-cli.go

GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -i -o ../bin/github-dns ../cmd/github-dns/github-dns.go
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w" -i -o ../bin/ghosts-cli ../cmd/ghosts-cli/ghosts-cli.go