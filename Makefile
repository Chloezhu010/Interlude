NAME=interlude
BUILD_DIR=./cmd/interlude

.PHONY: build run clean install

build:
	go build -o $(NAME) $(BUILD_DIR)

run:
	go run $(BUILD_DIR)

clean:
	rm -f $(NAME)

install:
	go install $(BUILD_DIR)
