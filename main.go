package main

import (
	"flag"
	"fmt"
	"os"
	"take_home_golang/etl"
)

func main() {
	src := flag.String("src", "", "Set .json filename")
	dest := flag.String("dest", "", "Set .csv filename")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-src] [source-file] [-dest] [dest-file]\n", os.Args[0])
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
	etl.Process()
}
