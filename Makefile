tag = $(shell git describe --tags --abbrev=0)

all: tests

tests:
	go test ./...

tests_coverage:
	go get golang.org/x/tools/cmd/cover
	go test ./... -test.coverprofile coverage.out
	go tool cover -html=coverage.out

build_local:
	go build dccm.go

build_all: build_linux build_windows build_osx

build_osx:
	env GOOS=darwin GOARCH=amd64 go build -o builds/ dccm.go
	zip -m -j builds/dccm_osx_amd64_${tag}.zip builds/dccm

build_windows:
	env GOOS=windows GOARCH=386 go build -o builds/ dccm.go
	zip -m -j builds/dccm_windows_386_${tag}.zip builds/dccm.exe

	env GOOS=windows GOARCH=amd64 go build -o builds/ dccm.go
	zip -m -j builds/dccm_windows_amd64_${tag}.zip builds/dccm.exe

build_linux:
	env GOOS=linux GOARCH=386 go build -o builds/ dccm.go
	zip -m -j builds/dccm_linux_386_${tag}.zip builds/dccm

	env GOOS=linux GOARCH=amd64 go build -o builds/ dccm.go
	zip -m -j builds/dccm_linux_amd64_${tag}.zip builds/dccm
