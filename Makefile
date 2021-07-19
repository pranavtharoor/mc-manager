.DEFAULT_GOAL := build

clean:
	rm -rf bin/*

build: clean
	go build -o bin/mcbot cmd/mcbot/mcbot.go

run:
	go run cmd/mcbot/mcbot.go

fix:
	gofmt -w ./
	go mod tidy