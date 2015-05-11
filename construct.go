package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "construct"
	app.Version = "0.1.0"
	app.Usage = "An application constructor with flexible template support."
	app.Author = "Brian Morton"
	app.Email = "brian@mmm.hm"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "template", Value: "https://github.com/bmorton/go-template", Usage: "a URL to the template to construct from", EnvVar: "CONSTRUCT_TEMPLATE"},
	}
	app.Commands = []cli.Command{
		NewNewCommand(),
	}
	app.Run(os.Args)
}

type Template struct {
	Name string
}
