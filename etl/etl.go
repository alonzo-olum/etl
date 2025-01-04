package etl

import (
	"encoding/csv"
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

func (et *ETL) Writer() *csv.Writer {
	return csv.NewWriter(et.out)
}

func (et *ETL) Process(writer *csv.Writer) {
	decoder := json.NewDecoder(et.in)
	doc := extract(decoder)
	result := transform(doc)
	load(result, writer)
}

type headers []string

func (et *ETL) WriteHeaders(writer *csv.Writer) {
	var data = headers{"modified", "publisher.name", "publisher.subOrganizationOf.name", "contactPoint.fn", "keyword"}
	if err := writeCsv(writer, data); err != nil {
		panic(err)
	}
}
