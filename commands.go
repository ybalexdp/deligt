package main

import (
	"github.com/urfave/cli"
	"github.com/urfave/cli/altsrc"
)

const baseURL = "https://kabutan.jp/stock/?code="
const financeURL = "https://kabutan.jp/stock/finance?code="

var commands = []cli.Command{
	commandUpdate,
}

var updateFlags = []cli.Flag{
	cli.BoolFlag{Name: "all", Usage: "update all company"},
	cli.StringFlag{Name: "number", Usage: "Specify a company number"},
	altsrc.NewStringFlag(cli.StringFlag{Name: "path"}),
}

var commandUpdate = cli.Command{
	Name:   "update",
	Usage:  "Update stock data",
	Action: doUpdate,
	Before: setConfig(updateFlags),
	Flags:  updateFlags,
}

func setConfig(flags []cli.Flag) cli.BeforeFunc {
	return altsrc.InitInputSource(flags, func() (altsrc.InputSourceContext, error) {
		return altsrc.NewYamlSourceFromFile("./conf/deligt.yaml")
	})
}
