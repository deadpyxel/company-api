BINARY_NAME=company-api
 
all: build test

build:
	mkdir uploads
	go build -o ${BINARY_NAME} main.go

run:
	go build -o ${BINARY_NAME} main.go
	rm uploads -r
	mkdir uploads
	./${BINARY_NAME}

test:
	go test -v

clean:
	go clean
	rm uploads -r
	rm *.db 