NAME=interlude
BUILD_DIR=./cmd/interlude

.PHONY: build start clean install

build:
	go build -o $(NAME) $(BUILD_DIR)

start: build
	./$(NAME) start

clean:
	rm -f $(NAME)

install:
	go install $(BUILD_DIR)
