package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

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
		Action: func(ctx *cli.Context) {
			generatorType, name, err := extractGenerateParameters(ctx)
			if err != nil {
				errorAndBail(err)
			}

			repoPath, err := getTemplatePath(ctx.GlobalString("template"))
			if err != nil {
				errorAndBail(err)
			}

			appPath, err := os.Getwd()
			if err != nil {
				errorAndBail(err)
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
			funcMap := template.FuncMap{
				"capitalize": strings.Title,
			}
			replacements := map[string]string{
				"api/resource.go":           fmt.Sprintf("api/%s.go", view.SingularName),
				"api/resources_resource.go": fmt.Sprintf("api/%s_resource.go", view.PluralName),
				"db/resource_record.go":     fmt.Sprintf("db/%s_record.go", view.SingularName),
			}
			err = filepath.Walk(templateRoot, directoryWalker(templateRoot, appPath, view, funcMap, replacements))
			if err != nil {
				errorAndBail(err)
			}
			instructionsFile := path.Join(templateRoot, "instructions.tmpl")
			if _, err := os.Stat(instructionsFile); err == nil {
				printInstructions(instructionsFile, view, funcMap)
			}
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

func printInstructions(filename string, view interface{}, funcMap template.FuncMap) {
	t := template.New(path.Base(filename)).Funcs(funcMap)
	t, err := t.ParseFiles(filename)
	if err != nil {
		errorAndBail(err)
	}

	err = t.Execute(os.Stdout, view)
	if err != nil {
		errorAndBail(err)
	}
}
