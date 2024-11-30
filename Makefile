GO = go
GOFMT = gofmt
BINARY_NAME = ledis-cli

install:
	$(GO) mod tidy

build: install
	$(GO) build -o $(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

format:
	$(GOFMT) -w .