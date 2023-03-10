build:
	go build -o bin/lightbulb .
test:
	cd lightbulb && go test -v .
