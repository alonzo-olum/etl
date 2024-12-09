SRC=main.go
BIN=main

test:
	go test ./...

build:
	GOARCH=amd64 GOOS=darwin go build -o bin/$(BIN)-darwin $(SRC)
	GOARCH=amd64 GOOS=linux go build -o bin/$(BIN)-linux $(SRC)
	GOARCH=amd64 GOOS=windows go build -o bin/$(BIN)-windows $(SRC)

clean:
	go clean
	rm $(BIN)-darwin
	rm $(BIN)-linux
	rm $(BIN)-windows

dep:
	go mod download
