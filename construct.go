package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/joho/godotenv"
	"github.com/mitchellh/go-homedir"
)

func main() {
	configPath, _ := homedir.Expand("~/.construct/config")
	godotenv.Load(configPath)

	app := cli.NewApp()
	app.Name = "construct"
	app.Version = "0.2.0"
	app.Usage = "An application constructor with flexible template support."
	app.Author = "Brian Morton"
	app.Email = "brian@mmm.hm"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "template",
			Value:  "https://github.com/bmorton/go-template",
			Usage:  "a URL to the template to construct from",
			EnvVar: "CONSTRUCT_TEMPLATE",
		},
		cli.StringFlag{
			Name:   "import-prefix",
			Value:  "",
			Usage:  "the prefix of the import path to use when generating files without a trailing slash (e.g. github.com/bmorton)",
			EnvVar: "CONSTRUCT_IMPORT_PREFIX",
		},
		cli.StringFlag{
			Name:   "source-path",
			Value:  "",
			Usage:  "the path where new projects should be created (defaults to using GOPATH)",
			EnvVar: "GOPATH",
		},
	}
	app.Commands = []cli.Command{
		NewNewCommand(),
		NewGenerateCommand(),
	}
	app.Run(os.Args)
}

type Template struct {
	Name string
}
