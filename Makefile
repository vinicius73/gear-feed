include ./.env

export $(shell sed 's/=.*//' ./.env)

build:
	- cd src/ && go build -o ../bin/gfeed
	- chmod +x bin/gfeed

run:
	- @cd src/ && go run *.go --token=${TELEGRAM_TOKEN}

run-scrapper:
	- @cd src/ && go run *.go scrapper