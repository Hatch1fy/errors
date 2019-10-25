all:

install: export GO111MODULE=on
install:
	go install github.com/Hatch1fy/errors

lint:
	golangci-lint run --enable-all

test:
	go test
