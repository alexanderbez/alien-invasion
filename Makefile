TARGET=alien-invasion-sim

all: test build

build:
	@go build -o ./$(TARGET)

test:
	@go test -v ./...
