.PHONY: all
all: test

.PHONY: test
test:
	go test -v ./...

.PHONY: test-cover
test-cover:
	go test -v -cover ./...

.PHONY: fmt
format:
	go fmt ./...

.PHONY: clean
clean:
	go clean
	rm -f coverage.out