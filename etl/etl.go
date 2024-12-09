package etl

import (
	"encoding/json"
	"io"
)

/*
	 Conceptually we are:
		- (E)xtracting the json data,
		- (T)ransforming it to CSV Data
		- (L)oading it by writing it to a CSV file
*/

// exporter type
type ETL struct {
	in  io.Reader // a json file here
	out io.Writer // a csv file here
}

// Factory Method for our new ETL object
func NewEtl(reader io.Reader, writer io.Writer) *ETL {
	return &ETL{
		in:  reader,
		out: writer,
	}
}

func (et *ETL) Process() {
	decoder := json.NewDecoder(et.in)
	doc := extract(decoder)
	result := transform(doc)
	load(result, et)
}
