# ETL Data Transform

### Prerequisites:

The program uses the following files by default if not specified
- src file: `in.json`
- dest file: `out.csv`

### Installing
This program retrofits the installation and build process using Makefile.
To install go modules run: `make dep`

### Testing
To test the application run: `make test`. Sample output:
``` go test -v ./etl/...
=== RUN   TestExtract_ValidJson_ReturnDoc
--- PASS: TestExtract_ValidJson_ReturnDoc (0.00s)
=== RUN   TestExtract_InvalidJson_ExceptionThrown
    worker_test.go:175: [TestExtract_InvalidJson_ExceptionThrown]: Found Error: true
--- PASS: TestExtract_InvalidJson_ExceptionThrown (0.00s)
=== RUN   TestExtract_NestedValidJson_ReturnDocList
--- PASS: TestExtract_NestedValidJson_ReturnDocList (0.00s)
=== RUN   TestTransform_DocWithoutSomeFields_ReturnCsvWithEmptyValue
--- PASS: TestTransform_DocWithoutSomeFields_ReturnCsvWithEmptyValue (0.00s)
PASS
ok      take_home_golang/etl    (cached)
```

### Building
To build the executable which has a pre-cursor for portability, run: `make build`. Sample output as:
```GOARCH=amd64 GOOS=darwin go build -o bin/main-darwin main.go
GOARCH=amd64 GOOS=linux go build -o bin/main-linux main.go
GOARCH=amd64 GOOS=windows go build -o bin/main-windows main.go
```

### Running the ETL Job
The executable is embedded in the */bin* directory.

Just: ` sh ./bin/main-<your-os> in.json`

Ideally, `<your-os>` is supposed to facilitate different executables for different machines but for now use main-darwin or whichever executable pleases.

It will write to a csv file by default, if one is not provided.

Use 
``` ./bin/main-darwin -h ``` for help. This is the output to expect.
```
Usage: ./bin/main-darwin [-src] [source-file] [-dest] [dest-file]
  -dest string
        Set .csv filename (default "out.csv")
  -src string
        Set .json filename (default "in.json")
```
### CSV Output

You can find your csv output in root directory as; `out.csv`

### CleanUp

To clean the myriad of executable run `make clean`. Sample output as:
```
go clean
rm -rf bin/main-darwin
rm -rf bin/main-linux
rm -rf bin/main-windows
```
