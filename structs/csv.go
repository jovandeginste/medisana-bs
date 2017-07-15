package structs

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"reflect"
	"strconv"
)

// ImportCsv load a csv file for a PersonID
func ImportCsv(person int) (result BodyMetrics) {
	csvFile := "csv/" + strconv.Itoa(person) + ".csv"
	if _, err := os.Stat(csvFile); os.IsNotExist(err) {
		return
	}

	// Load a TXT file.
	f, _ := os.Open(csvFile)

	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))
	for {
		var w BodyMetric
		err := unmarshal(r, &w)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		result = append(result, w)
	}
	return
}

func unmarshal(reader *csv.Reader, v interface{}) error {
	record, err := reader.Read()
	if err != nil {
		return err
	}
	s := reflect.ValueOf(v).Elem()
	if s.NumField() != len(record) {
		return &FieldMismatch{s.NumField(), len(record)}
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
			return &UnsupportedType{f.Type().String()}
		}
	}
	return nil
}

// FieldMismatch happens when the csv does not match our expectations
type FieldMismatch struct {
	expected, found int
}

func (e *FieldMismatch) Error() string {
	return "CSV line fields mismatch. Expected " + strconv.Itoa(e.expected) + " found " + strconv.Itoa(e.found)
}

// UnsupportedType happens when a type is not supported
type UnsupportedType struct {
	Type string
}

func (e *UnsupportedType) Error() string {
	return "Unsupported type: " + e.Type
}
