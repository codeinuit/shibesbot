BINARY = shibes

SRC_DIRECTORY = src/
BIN_DIRECTORY = .
SRC =   $(SRC_DIRECTORY)shibes.go     \
        $(SRC_DIRECTORY)discord.go     \
        $(SRC_DIRECTORY)dogequests.go     \


all: $(BINARY)

$(BINARY): clean
	go build -o $(BIN_DIRECTORY)/$(BINARY) $(SRC)

clean:
	rm -f bin/$(BINARY)

.PHONY: all clean
