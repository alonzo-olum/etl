package main

import (
	"flag"
	"fmt"
	"os"
	"take_home_golang/etl"
)

const (
	DefaultJson string = "in.json" // default json file if you don't give one
	DefaultCsv  string = "out.csv" // default csv file to be automatically created
)

func main() {
	src := flag.String("src", DefaultJson, "Set .json filename")
	dest := flag.String("dest", DefaultCsv, "Set .csv filename")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-src] [source-file] [-dest] [dest-file]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// open .json file as read only
	in, err := os.OpenFile(*src, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer in.Close()
	// open .csv as write only, create if does not exist and overwrite if it does
	out, err := os.OpenFile(*dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)

	if err != nil {
		panic(err)
	}
	defer out.Close()

	etl := etl.NewEtl(in, out)

	csvWriter := etl.Writer()
	defer csvWriter.Flush()

	etl.WriteHeaders(csvWriter)
	etl.Process(csvWriter)
}
