package main

import (
	"errors"
	"fmt"
	"os"
	"path"
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
		Action: func(ctx *cli.Context) error {
			name, err := extractNewParameters(ctx)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			repoPath, err := getTemplatePath(ctx.GlobalString("template"))
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			importPrefix := ctx.GlobalString("import-prefix")
			sourceRoot := fmt.Sprintf("%s/src/%s", ctx.GlobalString("source-path"), importPrefix)

			appPath, err := makeAppDir(name, sourceRoot)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			view := &viewTemplate{
				Name:         name,
				ImportPrefix: importPrefix,
			}

			fmt.Println("Creating files...")
			structureRoot := repoPath + "/structure"
			err = filepath.Walk(structureRoot, directoryWalker(structureRoot, appPath, view, funcMap, map[string]string{}))
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			instructionsFile := path.Join(structureRoot, "instructions.tmpl")
			if _, err := os.Stat(instructionsFile); err == nil {
				err = printInstructions(instructionsFile, view, funcMap)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
			}

			return nil
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
		err := os.MkdirAll(path, 0755)
		return path, err
	} else {
		return "", errors.New("directory already exists, aborting")
	}
	return "", errors.New("failed to create directory")
}
