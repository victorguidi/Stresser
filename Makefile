build:
	@go build -v -o stresser ./src/...

run: build
	@./stresser

.PHONY: build run
