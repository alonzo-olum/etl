package etl

import (
	"fmt"
	"io"

	"github.com/gocarina/gocsv"
)

// write data to csv file
func writeCsv(writer io.Writer, data []csvData) error {
	if err := gocsv.Marshal(&data, writer); err != nil {
		return fmt.Errorf("write error: %v occurred!", err)
	}
	return nil
}
