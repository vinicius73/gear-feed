build:
	- cd src/ && go build -o ../bin/gfeed
	- chmod +x bin/gfeed

run:
	- cd src/ && go run *.go