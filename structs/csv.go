package structs

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var (
	// FieldMismatch happens when the csv does not match our expectations
	ErrCSVFieldMismatch = errors.New("line fields mismatch in CSV")
	// UnsupportedType happens when a type is not supported
	ErrCSVUnsupportedType = errors.New("unsupported type")
)

func importCSVLogger() log.FieldLogger {
	return log.WithField("component", "import-csv")
}

// ImportCsv load a csv file for a PersonID
func ImportCsv(person int) BodyMetrics {
	csvFile := "csv/" + strconv.Itoa(person) + ".csv"
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		return nil
	}

	// Load a TXT file.
	f, err := os.Open(csvFile)
	if err != nil {
		return nil
	}

	var result BodyMetrics

	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))

	for {
		var w BodyMetric

		if err := unmarshal(r, &w); err != nil {
			if err == io.EOF {
				break
			}

			importCSVLogger().Fatalf("Error importing file: %s", err)
		}

		result = append(result, w)
	}

	return result
}

func unmarshal(reader *csv.Reader, v interface{}) error {
	record, err := reader.Read()
	if err != nil {
		return err
	}

	s := reflect.ValueOf(v).Elem()
	if s.NumField() != len(record) {
		return fmt.Errorf("%w; expected %d, found %d", ErrCSVFieldMismatch, s.NumField(), len(record))
	}

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)

		switch f.Type().String() {
		case "string":
			f.SetString(record[i])
		case "int":
			ival, err := strconv.ParseInt(record[i], 10, 0)
			if err != nil {
				return err
			}

			f.SetInt(ival)
		case "float32":
			ival, err := strconv.ParseFloat(record[i], 32)
			if err != nil {
				return err
			}

			f.SetFloat(ival)
		default:
			return fmt.Errorf("%w: %s", ErrCSVUnsupportedType, f.Type().String())
		}
	}

	return nil
}
