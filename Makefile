SHELL=/bin/bash

all: unit_test build system_test

unit_test:
	go test ./...

build_all: build_linux build_darwin build_windows build_freebsd

build_linux:
	GOOS=linux GOARCH=amd64 go build -o ./artefacts/modcop_linux_amd64 cmd/modcop/main.go

build_darwin:
	GOOS=darwin GOARCH=amd64 go build -o ./artefacts/modcop_darwin_amd64 cmd/modcop/main.go

build_windows:
	GOOS=windows GOARCH=amd64 go build -o ./artefacts/modcop_windows_amd64 cmd/modcop/main.go

build_freebsd:
	GOOS=freebsd GOARCH=amd64 go build -o ./artefacts/modcop_freebsd_amd64 cmd/modcop/main.go

build:
	go build -o ./artefacts/modcop cmd/modcop/main.go

tests:
	go test ./...

system_test: build
	go test -tags cli_tests ./cli_tests