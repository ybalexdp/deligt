package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Stock struct {
	Number string `csv:"コード"`
	Name   string `csv:"銘柄"`
	Price  string `csv:"株価"`
	Per    string `csv:"PER"`
	Pbr    string `csv:"PBR"`
	Eps    string `csv:"EPS"`
	Bps    string `csv:"BPS"`
	Roe    string `csv:"ROE"`
}

func (d *Stock) getData(ctx context.Context) error {
	doc, err := goquery.NewDocument(baseURL + d.Number)
	if err != nil {
		return err
	}

	// 指定された銘柄があるか確認する
	ex := doc.Find("p").First().Children().Text()
	if ex == "" {
		return fmt.Errorf("specified description does not exist.")
	}

	datalist := []string{}

	doc.Find("div.company_block").Each(func(_ int, s *goquery.Selection) {
		d.Name = s.Children().First().Next().Text()
	})

	doc.Find("span.kabuka").Each(func(_ int, s *goquery.Selection) {
		d.Price = s.Text()
	})

	doc.Find("#stockinfo_i3 > table > tbody > tr > td").Each(func(_ int, s *goquery.Selection) {
		datalist = append(datalist, s.Text())
	})

	d.Per = strings.Split(datalist[0], "倍")[0]
	d.Pbr = strings.Split(datalist[1], "倍")[0]

	doc.Find("span.kubun1").Each(func(_ int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "予") {
			s = s.Parent().Next().Next().Next().Next()
			d.Eps = s.Text()
		}
	})

	doc, err = goquery.NewDocument(financeURL + d.Number)
	if err != nil {
		return err
	}

	doc.Find("td.oc_btn1").Each(func(_ int, s *goquery.Selection) {
		d.Bps = s.Parent().Next().Next().Next().Children().Next().First().Text()
	})

	doc.Find("div.fin_f_t4_d > table > tbody > tr > th > span.kubun1").Each(func(_ int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "予") {
			s = s.Parent().Next().Next().Next().Next()
			d.Roe = s.Text()
		}
	})

	doneCh := make(chan bool, 1)
	doneCh <- true

	select {
	case <-doneCh:
		return nil

	case <-ctx.Done():
		<-doneCh
		return ctx.Err()
	}

}
