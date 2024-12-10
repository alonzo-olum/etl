Centripetal Confidential

# Data Transform

### Prerequisites:

The program uses the following files by default if not specified
- src file: `in.json`
- dest file: `out.csv`

### Installing
This program retrofits the installation and build process using Makefile.
To install go modules run: `make dep`

### Testing
To test the application run: `make test`

### Building
To build the executable which has a pre-cursor for portability, run: `make build`

### Running the ETL Job
The executable is embedded in the */bin* directory.
Just:
    ``` sh ./bin/main-<your-os> -src in.json```
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
