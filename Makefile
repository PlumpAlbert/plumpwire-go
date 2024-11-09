include .env
export $(shell sed 's/=.*//' .env)

BINARY_NAME=plumpwire
BINARY_PATH=.bin/${BINARY_NAME}

build:
	GOARCH=amd64 GOOS=linux go build -a -o ${BINARY_PATH} *.go

run: build
	./${BINARY_PATH}

clean:
	go clean
	rm ${BINARY_PATH}

test:
	env
