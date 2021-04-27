BINARY_NAME=company-api
VERSION=0.6.0
IMAGE_NAME := "deadpyxel/company-api"
 
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

package:
	docker build --build-arg VERSION=${VERSION} --build-arg GIT_COMMIT=$(GIT_COMMIT) -t $(IMAGE_NAME):local .