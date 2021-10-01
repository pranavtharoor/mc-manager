.DEFAULT_GOAL := build
CONFIG_FILE := config.yml

clean:
	rm -rf bin/*

build: clean
	go build -o bin/mcbot cmd/mcbot/mcbot.go

run:
	go run cmd/mcbot/mcbot.go

fix:
	gofmt -w ./
	go mod tidy

docker_build:
	docker build -t mc-manager .

docker_run:
	[ -f ${PWD}/$(CONFIG_FILE) ] && exit
	docker run -it -v ${PWD}/$(CONFIG_FILE):/app/$(CONFIG_FILE) --entrypoint=./mcbot  mc-manager