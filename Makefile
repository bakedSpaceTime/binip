.PHONY: build tidy

build:
	go build .

tidy:
	go mod tidy
