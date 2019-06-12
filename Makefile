SHELL=/bin/bash

all: unit_test build system_test

unit_test:
	go test ./...

build_all: build_linux build_darwin build_windows build_freebsd

build_linux:
	GOOS=linux GOARCH=amd64 go build -o ./artefacts/modcons_linux_amd64 cmd/modcons/main.go

build_darwin:
	GOOS=darwin GOARCH=amd64 go build -o ./artefacts/modcons_darwin_amd64 cmd/modcons/main.go

build_windows:
	GOOS=windows GOARCH=amd64 go build -o ./artefacts/modcons_windows_amd64 cmd/modcons/main.go

build_freebsd:
	GOOS=freebsd GOARCH=amd64 go build -o ./artefacts/modcons_freebsd_amd64 cmd/modcons/main.go

build:
	go build -o ./artefacts/modcons cmd/modcons/main.go

tests:
	go test ./...

system_test: build
	go test -tags cli_tests ./cli_tests