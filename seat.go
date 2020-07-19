package main

import (
	"encoding/csv"
	"io"
	"os"
	"reflect"
)

type Seat struct {
	stockList []Stock
	filePath  string
	recordNum int
}

func (seat *Seat) Set(csvPath string) error {
	seat.filePath = csvPath
	file, err := os.Open(seat.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var line []string
	i := 0
	for {
		line, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if err.(*csv.ParseError).Err != csv.ErrFieldCount {
				return err
			}
		}
		if line[0] == "コード" {
			continue
		}
		i++
		seat.stockList = append(seat.stockList, Stock{Number: line[0]})
	}
	seat.recordNum = i
	return nil

}

func Unmarshal(r *csv.Reader, v interface{}) error {
	var err error
	rv := reflect.ValueOf(v)
	rt := rv.Type().Elem()

	records, err := r.ReadAll()
	c := len(records)
	if err != nil {
		if err.(*csv.ParseError).Err != csv.ErrFieldCount {
			return err
		}
		c, err = getRecordNum(r)
		if err != nil {
			return err
		}
	}
	slice := reflect.MakeSlice(rt, c, c)

	for i, record := range records {
		for j, column := range record {
			slice.Index(i).Field(j).SetString(column)
		}

	}

	rv.Elem().Set(slice.Slice3(0, c, c))
	return err
}

func getRecordNum(r *csv.Reader) (int, error) {
	i := 0
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if err.(*csv.ParseError).Err != csv.ErrFieldCount {
				return 0, err
			}
		}
		if line[0] == "コード" {
			continue
		}
		i++
	}
	return i, nil
}
