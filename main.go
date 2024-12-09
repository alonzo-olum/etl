package main

import (
	"os"
	"take_home_golang/etl"
)

func main() {
	src, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer src.Close()
	dest, err := os.Open(os.Args[2])
	if err != nil {
		panic(err)
	}
	defer dest.Close()

	etl := etl.NewEtl(src, dest)
	etl.Process()
}
