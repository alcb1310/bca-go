BINARY_NAME=bca

build:
	go build -o /build/${BINARY_NAME} .

run: build:
	./build/${BINARY_NAME}

clean:
	go clean
	rm ./build/${BINARY_NAME}
