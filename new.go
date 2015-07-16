package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
	"gopkg.in/libgit2/git2go.v22"
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
			err = filepath.Walk(structureRoot, func(p string, info os.FileInfo, err error) error {
				if p == structureRoot {
					return nil
				}

				dest := appPath + subtractRoot(structureRoot, p)
				if strings.HasSuffix(p, ".tmpl") {
					dest = dest[:len(dest)-len(".tmpl")]
					t := template.New(path.Base(p))
					t, err = t.ParseFiles(p)
					if err != nil {
						errorAndBail(err)
					}
					f, err := os.Create(dest)
					if err != nil {
						errorAndBail(err)
					}
					err = t.Execute(f, view)
					if err != nil {
						errorAndBail(err)
					}
					f.Close()
				} else {
					if info.IsDir() {
						os.Mkdir(dest, 0755)
					} else {
						copyFile(p, dest)
					}
				}
				fmt.Println("-- " + dest)
				return nil
			})
			if err != nil {
				errorAndBail(err)
			}
		},
	}
}

func subtractRoot(root string, p string) string {
	return strings.Replace(p, root, "", -1)
}

func extractNewParameters(ctx *cli.Context) (string, error) {
	if len(ctx.Args()) == 0 {
		return "", errors.New("Name required")
	}
	name := ctx.Args()[0]

	return name, nil
}

func getTemplate(repoPath string) (*Template, error) {
	var parsed Template

	file, err := ioutil.ReadFile(fmt.Sprintf("%s/template.json", repoPath))
	if err != nil {
		return &parsed, err
	}

	err = json.Unmarshal(file, &parsed)
	return &parsed, err
}

func getTemplateRepo(templateURL string) (*git.Repository, string, error) {
	var repo *git.Repository
	var err error

	currentUser, err := user.Current()
	if err != nil {
		return &git.Repository{}, "", err
	}

	parsedTemplate, err := url.Parse(templateURL)
	if err != nil {
		return &git.Repository{}, "", err
	}
	path := fmt.Sprintf("%s/.construct/src/%s%s", currentUser.HomeDir, parsedTemplate.Host, parsedTemplate.Path)

	fmt.Printf("Attempting to clone %s to %s...\n", templateURL, path)
	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("-- Cloning %s...\n", templateURL)
			repo, err = git.Clone(templateURL, path, &git.CloneOptions{})
		} else {
			return &git.Repository{}, "", err
		}
	} else {
		fmt.Printf("-- Found git repository at %s...\n", path)
		repo, err = git.OpenRepository(path)
		if err != nil {
			return repo, path, err
		}
		remote, err := repo.LookupRemote("origin")
		if err != nil {
			return repo, path, err
		}

		refSpecs, err := remote.FetchRefspecs()
		if err != nil {
			return repo, path, err
		}

		err = remote.Fetch(refSpecs, nil, "")
		if err != nil {
			return repo, path, err
		}

		branch, err := repo.LookupBranch("origin/master", git.BranchRemote)
		if err != nil {
			return repo, path, err
		}
		commit, err := repo.LookupCommit(branch.Target())
		if err != nil {
			return repo, path, err
		}
		tree, err := commit.Tree()
		if err != nil {
			return repo, path, err
		}
		err = repo.CheckoutTree(tree, &git.CheckoutOpts{Strategy: git.CheckoutForce})
		if err != nil {
			return repo, path, err
		}
		err = repo.SetHeadDetached(branch.Target(), nil, "")
	}

	return repo, path, err
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

func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}

	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}

	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}

	return d.Close()
}

func errorAndBail(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	os.Exit(1)
}
