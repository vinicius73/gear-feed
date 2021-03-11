include ./.env

export $(shell sed 's/=.*//' ./.env)

export CGO_ENABLED=0
export GOOS=linux

BUILD_FLAGS=-a -installsuffix cgo -ldflags '-s -w -extldflags "-static"'

init:
	- cd src/ && go get
	- cd src/ && go mod vendor

build:
	- cd src/ && go build ${BUILD_FLAGS} -o ../bin/gfeed
	- chmod +x bin/gfeed

build-functions:
	- cd src/functions && go build -tags scrapper ${BUILD_FLAGS} -o ../../bin/function-scrapper scrapper.go
	- chmod +x bin/function-scrapper

run:
	- @cd src/ && go run *.go --token=${TELEGRAM_TOKEN}

run-scrapper:
	- @cd src/ && go run *.go scrapper

release: build build-functions
	- upx --best bin/gfeed
	- upx --best bin/function-scrapper
	- ls -lh bin/