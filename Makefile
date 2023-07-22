include ./.env

export $(shell sed 's/=.*//' ./.env)

export CGO_ENABLED=0
export GOOS=linux
export ENTRIES_TABLE=gamer-feed-dev-entries
export LOG_LEVEL=debug

GIT_COMMIT=$(shell git rev-list -1 HEAD)
APP_VERSION=$(shell node -p "require('./package.json').version")
BUILD_DATE=$(shell date '+%Y-%m-%d__%H:%M:%S')

LDFLAGS=-ldflags "-X gfeed/domains.commit=${GIT_COMMIT} -X gfeed/domains.version=${APP_VERSION} -X gfeed/domains.buildDate=${BUILD_DATE} "
BUILD_FLAGS=-a -installsuffix cgo -ldflags '-s -w -extldflags "-static"' ${LDFLAGS}

version:
	- @echo Commit: ${GIT_COMMIT}
	- @echo Version: ${APP_VERSION}

init:
	- cd src/ && go get
	- cd src/ && go mod vendor

build: version
	- cd src/ && go build ${BUILD_FLAGS} -o ../bin/gfeed
	- chmod +x bin/gfeed

build-functions: version
	- cd src/functions && go build -tags scrapper ${BUILD_FLAGS} -o ../../bin/function-scrapper scrapper.go
	- chmod +x bin/function-scrapper

run:
	- @cd src/ && go run *.go --token=${TELEGRAM_TOKEN} --dry

run-scrapper:
	- @cd src/ && go run *.go scrapper

release: build-functions
	- upx --best bin/function-scrapper
	- ls -lh bin/