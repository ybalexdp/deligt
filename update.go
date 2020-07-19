package main

import (
	"context"
	"encoding/csv"
	"os"
	"sync"
	"time"

	"github.com/urfave/cli"
)

func doUpdate(c *cli.Context) error {
	var err error
	all := c.Bool("all")
	number := c.String("number")
	csvPath := c.String("path")
	setup(csvPath)
	seat := new(Seat)
	seat.Set(csvPath)
	if all {
		err = seat.updateAll()
		return err
	} else if number != "" {
		err = seat.update(number)
		return err
	}
	return nil
}

func (seat *Seat) updateAll() error {
	file, err := os.Open(seat.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	datalist := []Stock{}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	for _, stock := range seat.stockList {
		errCh := make(chan error, 1)
		wg.Add(1)
		go func(s Stock) {
			errCh <- s.getData(ctx)
			datalist = append(datalist, s)
			defer wg.Done()
		}(stock)
	}
	wg.Wait()

	for _, data := range datalist {
		seat.updateCsv(data)
	}

	return nil
}

func (seat *Seat) update(number string) error {
	stock := Stock{Number: number}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := stock.getData(ctx)
	if err != nil {
		return err
	}

	err = seat.updateCsv(stock)

	return err
}

func (seat *Seat) updateCsv(stock Stock) error {
	var stocklist []Stock
	filecsv, err := os.Open(seat.filePath)
	if err != nil {
		return err
	}
	defer filecsv.Close()

	reader := csv.NewReader(filecsv)
	Unmarshal(reader, &stocklist)

	if err := os.Remove(seat.filePath); err != nil {
		return err
	}
	setup(seat.filePath)

	file, err := os.OpenFile(seat.filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	flag := false
	for _, s := range stocklist {
		if s.Number == "コード" {
			continue
		}
		if s.Number == stock.Number {
			s = stock
			flag = true
		}
		writer.Write([]string{
			s.Number,
			s.Name,
			s.Price,
			s.Per,
			s.Pbr,
			s.Eps,
			s.Bps,
			s.Roe,
		})
		writer.Flush()

	}
	if !flag {
		if err = seat.add(stock); err != nil {
			return err
		}
	}
	return err
}

func (seat *Seat) add(stock Stock) error {
	var err error
	filecsv, err := os.Open(seat.filePath)
	if err != nil {
		panic(err)
	}
	defer filecsv.Close()

	file, err := os.OpenFile(seat.filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	writer.Write([]string{
		stock.Number,
		stock.Name,
		stock.Price,
		stock.Per,
		stock.Pbr,
		stock.Eps,
		stock.Bps,
		stock.Roe,
	})
	writer.Flush()

	return err
}
