include ./.env

export $(shell sed 's/=.*//' ./.env)

init:
	- cd src/ && go get
	- cd src/ && go mod vendor

build:
	- cd src/ && go build -o ../bin/gfeed
	- chmod +x bin/gfeed

run:
	- @cd src/ && go run *.go --token=${TELEGRAM_TOKEN}

run-scrapper:
	- @cd src/ && go run *.go scrapper