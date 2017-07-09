package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/gocarina/gocsv"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
)

func CreateCsvDir(file string) {
	path := filepath.Dir(file)
	mode := os.FileMode(0700)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, mode)
	}
}

func ExportCsv(filename string, weights BodyMetrics) {
	fmt.Printf("Writing to file '%s'.\n", filename)
	CreateCsvDir(filename)

	f, err := os.Create(filename)
	if err != nil {
		fmt.Errorf("%#v", err)
	}
	defer f.Close()

	err = gocsv.MarshalWithoutHeaders(&weights, f)

	if err != nil {
		fmt.Errorf("%#v", err)
	}
}

func ImportCsv(filename string) (result BodyMetrics) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return
	}

	// Load a TXT file.
	f, _ := os.Open(filename)

	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))
	for {
		var w BodyMetric
		err := Unmarshal(r, &w)
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

func Unmarshal(reader *csv.Reader, v interface{}) error {
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

type FieldMismatch struct {
	expected, found int
}

func (e *FieldMismatch) Error() string {
	return "CSV line fields mismatch. Expected " + strconv.Itoa(e.expected) + " found " + strconv.Itoa(e.found)
}

type UnsupportedType struct {
	Type string
}

func (e *UnsupportedType) Error() string {
	return "Unsupported type: " + e.Type
}
