SRC=main.go
BIN=main
TEST=etl

test:
	go test -v ./$(TEST)/...

build:
	GOARCH=amd64 GOOS=darwin go build -o bin/$(BIN)-darwin $(SRC)
	GOARCH=amd64 GOOS=linux go build -o bin/$(BIN)-linux $(SRC)
	GOARCH=amd64 GOOS=windows go build -o bin/$(BIN)-windows $(SRC)

clean:
	go clean
	rm -rf bin/$(BIN)-darwin
	rm -rf bin/$(BIN)-linux
	rm -rf bin/$(BIN)-windows

dep:
	go mod download
