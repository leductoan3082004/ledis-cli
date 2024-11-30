GO = go
GOFMT = gofmt
BINARY_NAME = ledis-cli

build:
	$(GO) build -o $(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

format:
	$(GOFMT) -w .