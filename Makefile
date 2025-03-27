BIN_DIR=bin
CLI_BIN=$(BIN_DIR)/llex

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(CLI_BIN): | $(BIN_DIR)
	go build -o $(CLI_BIN) ./cmd

cli: $(CLI_BIN)

clean:
	go clean
	rm -rf $(BIN_DIR)

.PHONY: all cli clean
all: cli
