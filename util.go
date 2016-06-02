package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
	"text/template"

	"gopkg.in/libgit2/git2go.v22"
)

func getTemplatePath(templateURL string) (string, error) {
	if _, err := os.Stat(templateURL); err == nil {
		return templateURL, nil
	} else {
		return getTemplateRepo(templateURL)
	}
}

func getTemplateRepo(templateURL string) (string, error) {
	var repo *git.Repository
	var err error

	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	parsedTemplate, err := url.Parse(templateURL)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("%s/.construct/src/%s%s", currentUser.HomeDir, parsedTemplate.Host, parsedTemplate.Path)

	fmt.Printf("Attempting to clone %s to %s...\n", templateURL, path)
	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("-- Cloning %s...\n", templateURL)
			repo, err = git.Clone(templateURL, path, &git.CloneOptions{})
		} else {
			return "", err
		}
	} else {
		fmt.Printf("-- Found git repository at %s...\n", path)
		repo, err = git.OpenRepository(path)
		if err != nil {
			return path, err
		}
		remote, err := repo.LookupRemote("origin")
		if err != nil {
			return path, err
		}

		refSpecs, err := remote.FetchRefspecs()
		if err != nil {
			return path, err
		}

		err = remote.Fetch(refSpecs, nil, "")
		if err != nil {
			return path, err
		}

		branch, err := repo.LookupBranch("origin/master", git.BranchRemote)
		if err != nil {
			return path, err
		}
		commit, err := repo.LookupCommit(branch.Target())
		if err != nil {
			return path, err
		}
		tree, err := commit.Tree()
		if err != nil {
			return path, err
		}
		err = repo.CheckoutTree(tree, &git.CheckoutOpts{Strategy: git.CheckoutForce})
		if err != nil {
			return path, err
		}
		err = repo.SetHeadDetached(branch.Target(), nil, "")
	}

	return path, err
}

func directoryWalker(templateRoot string, appPath string, view interface{}, funcMap template.FuncMap, renameMap map[string]string) func(p string, info os.FileInfo, err error) error {
	return func(p string, info os.FileInfo, err error) error {
		if p == templateRoot {
			return nil
		}
		relativePath := subtractRoot(templateRoot, p)

		if relativePath == "/instructions.tmpl" {
			return nil
		}

		dest := appPath + relativePath
		if strings.HasSuffix(p, ".tmpl") {
			dest = dest[:len(dest)-len(".tmpl")]
			for matcher, replacement := range renameMap {
				r := regexp.MustCompile(matcher)
				dest = r.ReplaceAllString(dest, replacement)
			}
			t := template.New(path.Base(p))
			if funcMap != nil {
				t = t.Funcs(funcMap)
			}
			t, err = t.ParseFiles(p)
			if err != nil {
				return err
			}
			f, err := os.Create(dest)
			if err != nil {
				return err
			}
			err = t.Execute(f, view)
			if err != nil {
				return err
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
	}
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

func subtractRoot(root string, p string) string {
	return strings.Replace(p, root, "", -1)
}

func printInstructions(filename string, view interface{}, funcMap template.FuncMap) error {
	t := template.New(path.Base(filename)).Funcs(funcMap)
	t, err := t.ParseFiles(filename)
	if err != nil {
		return err
	}

	fmt.Println("------------------------------------------------------------------------------")
	err = t.Execute(os.Stdout, view)
	if err != nil {
		return err
	}
	fmt.Println("------------------------------------------------------------------------------")

	return nil
}

var funcMap = template.FuncMap{
	"capitalize": strings.Title,
}
