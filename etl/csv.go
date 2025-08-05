package etl

import (
	"encoding/csv"
	"reflect"
)

// write data to file in csv format
func writeCsv(writer *csv.Writer, data interface{}) error {
	switch t := data.(type) {
	case headers:
		write(t, writer)
	case accumulator:
		marshallAndWrite(t, writer)
	default:
	}
	return nil
}

func marshallAndWrite(acc accumulator, writer *csv.Writer) {
	for _, a := range acc {
		rows := make([]string, 0)
		data := reflect.ValueOf(a)
		fields := data.NumField()
		for idx := range fields {
			rows = append(rows, data.Field(idx).String())
		}
		write(rows, writer)
	}
}

func write(data []string, writer *csv.Writer) {
	writer.Write(data)
}
