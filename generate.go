package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/gedex/inflector"
)

type generateViewTemplate struct {
	AppName      string
	ImportPrefix string
	SingularName string
	PluralName   string
}

// NewGenerateCommand returns the CLI command for "generate".
func NewGenerateCommand() cli.Command {
	return cli.Command{
		Name:        "generate",
		ShortName:   "g",
		Usage:       "Generates a set of files for the given type",
		Description: "generate [type] [name]",
		Action: func(ctx *cli.Context) error {
			generatorType, name, err := extractGenerateParameters(ctx)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			repoPath, err := getTemplatePath(ctx.GlobalString("template"))
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			appPath, err := os.Getwd()
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			appName := path.Base(appPath)
			importPrefix := ctx.GlobalString("import-prefix")

			view := &generateViewTemplate{
				AppName:      appName,
				ImportPrefix: importPrefix,
				SingularName: inflector.Singularize(name),
				PluralName:   inflector.Pluralize(name),
			}

			fmt.Println("Creating files...")
			templateRoot := path.Join(repoPath, generatorType)
			replacements := map[string]string{
				"api/resource.go":           fmt.Sprintf("api/%s.go", view.SingularName),
				"api/resources_resource.go": fmt.Sprintf("api/%s_resource.go", view.PluralName),
				"db/resource_record.go":     fmt.Sprintf("db/%s_record.go", view.SingularName),
			}
			err = filepath.Walk(templateRoot, directoryWalker(templateRoot, appPath, view, funcMap, replacements))
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}

			instructionsFile := path.Join(templateRoot, "instructions.tmpl")
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

func extractGenerateParameters(ctx *cli.Context) (string, string, error) {
	if len(ctx.Args()) == 0 {
		return "", "", errors.New("Type required")
	} else if len(ctx.Args()) == 1 {
		return "", "", errors.New("Name required")
	}

	generatorType := ctx.Args()[0]
	name := strings.ToLower(ctx.Args()[1])

	return generatorType, name, nil
}
