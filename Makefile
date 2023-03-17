BINARY = shibes

BIN_DIRECTORY = bin


all: $(BINARY)

$(BINARY): clean
	go build -o $(BIN_DIRECTORY)/$(BINARY) cmd/shibesbot/*.go

test:
	go test ./... -v

clean:
	rm -f bin/$(BINARY)


.PHONY: all clean
