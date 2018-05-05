BINARY = shibes

SRC_DIRECTORY = src/
BIN_DIRECTORY = .
SRC = $(SRC_DIRECTORY)shibes.go

all: $(BINARY)

$(BINARY):
	go build -o $(BIN_DIRECTORY)/$(BINARY) $(SRC)

clean:
	rm -f bin/$(BINARY)

.PHONY: all clean
