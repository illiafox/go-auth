BUILD=./cmd/server

all: clean build run

run:
	(cd $(BUILD) && ./bin $(ARGS))

build: clean
	go build -o $(BUILD)/bin $(BUILD)

clean:
	if [ -f $(BUILD)/bin ]; then rm $(BUILD)/bin; fi
