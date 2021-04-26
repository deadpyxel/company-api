BINARY_NAME=company-api
 
all: build test

build:
	go build -o ${BINARY_NAME} main.go

run:
	go build -o ${BINARY_NAME} main.go
	./${BINARY_NAME}

test:
	go test -v

clean:
	go clean