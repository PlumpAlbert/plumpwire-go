include .env
export $(shell sed 's/=.*//' .env)

BINARY_NAME=plumpwire
BINARY_PATH=build/${BINARY_NAME}

build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_PATH} main.go

run: build
	./${BINARY_PATH}

clean:
	go clean
	rm ${BINARY_PATH}

test:
	env
