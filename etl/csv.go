package etl

import (
	"encoding/csv"
	"reflect"
)

// write data to file in csv format
func writeCsv(writer *csv.Writer, data interface{}) error {
	switch t := data.(type) {
	case headers:
		write(writer, t)
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
		write(writer, rows)
	}
}

func write(writer *csv.Writer, data []string) {
	writer.Write(data)
}
