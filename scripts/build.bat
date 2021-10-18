cd %~dp0
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w" -i -o ../bin/github-dns.exe ../cmd/github-dns/github-dns.go
go build -ldflags "-s -w" -i -o ../bin/ghosts-cli.exe ../cmd/ghosts-cli/ghosts-cli.go
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w" -i -o ../bin/github-dns ../cmd/github-dns/github-dns.go
go build -ldflags "-s -w" -i -o ../bin/ghosts-cli ../cmd/ghosts-cli/ghosts-cli.go