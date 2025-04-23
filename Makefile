.DEFAULT_GOAL := build

build:
	go build -o scribe cmd/scribe/main.go

run:
	go run cmd/scribe/main.go

