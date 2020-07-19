package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	if err := newApp().Run(os.Args); err != nil {
		exitCode := 1
		if excoder, ok := err.(cli.ExitCoder); ok {
			exitCode = excoder.ExitCode()
		}
		printError(os.Stderr, "%s", err)
		os.Exit(exitCode)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = "deligt"
	app.Usage = "collect stock data"
	app.Author = "ybalexdp"
	app.Email = "ybalexdp@gmail.com"
	app.Commands = commands

	return app
}
