package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
)

type viewTemplate struct {
	Name         string
	ImportPrefix string
}

// NewNewCommand returns the CLI command for "new".
func NewNewCommand() cli.Command {
	return cli.Command{
		Name:        "new",
		ShortName:   "n",
		Usage:       "Creates a new application using the configured template",
		Description: "new [name]",
		Action: func(ctx *cli.Context) {
			name, err := extractNewParameters(ctx)
			if err != nil {
				errorAndBail(err)
			}

			_, repoPath, err := getTemplateRepo(ctx.GlobalString("template"))
			if err != nil {
				errorAndBail(err)
			}

			importPrefix := ctx.GlobalString("import-prefix")
			sourceRoot := fmt.Sprintf("%s/src/%s", ctx.GlobalString("source-path"), importPrefix)
			appPath, err := makeAppDir(name, sourceRoot)
			if err != nil {
				errorAndBail(err)
			}

			view := &viewTemplate{
				Name:         name,
				ImportPrefix: importPrefix,
			}

			fmt.Println("Creating files...")
			structureRoot := repoPath + "/structure"
			err = filepath.Walk(structureRoot, directoryWalker(structureRoot, appPath, view, nil, map[string]string{}))
			if err != nil {
				errorAndBail(err)
			}
		},
	}
}

func extractNewParameters(ctx *cli.Context) (string, error) {
	if len(ctx.Args()) == 0 {
		return "", errors.New("Name required")
	}
	name := ctx.Args()[0]

	return name, nil
}

func makeAppDir(name string, root string) (string, error) {
	path := filepath.FromSlash(fmt.Sprintf("%s/%s", root, name))

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("Creating directory...\n")
		os.Mkdir(path, 0755)
		return path, nil
	} else {
		return "", errors.New("directory already exists, aborting")
	}
	return "", errors.New("failed to create directory")
}
