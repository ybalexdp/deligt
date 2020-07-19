package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/mitchellh/colorstring"
)

func setup(csvPath string) error {
	_, ext := os.Stat(csvPath)
	if ext == nil {
		return nil
	}
	file, err := os.OpenFile(csvPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	writer.Write([]string{"コード", "銘柄", "株価", "PER", "PBR", "EPS", "BPS", "ROE"})
	writer.Flush()

	return err
}

func printError(w io.Writer, format string, args ...interface{}) {
	format = fmt.Sprintf("[red]%s[reset]\n", format)
	fmt.Fprintf(w, colorstring.Color(fmt.Sprintf(format, args...)))
}
