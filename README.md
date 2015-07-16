**Warning:** Construct is still under active development and things will be changing rapidly.

# Construct

`construct` allows you to create new applications from any git-based template repository.


## Example

```ShellOutput
$ export CONSTRUCT_TEMPLATE=https://github.com/bmorton/go-template
$ construct new myapp
Attempting to clone https://github.com/bmorton/go-template to /Users/bmorton/.construct/src/github.com/bmorton/go-template...
-- Found git repository at /Users/bmorton/.construct/src/github.com/bmorton/go-template...
Creating directory...
Creating files...
-- ./myapp/README.md.tmpl
-- ./myapp/main.go
```


## Usage

```
$ construct --help
NAME:
   construct - An application constructor with flexible template support.

USAGE:
   construct [global options] command [command options] [arguments...]

VERSION:
   0.1.0

AUTHOR(S):
   Brian Morton <brian@mmm.hm>

COMMANDS:
   new, n Creates a new application using the configured template
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --template "https://github.com/bmorton/go-template"  a URL to the template to construct from [$CONSTRUCT_TEMPLATE]
   --import-prefix                                      the prefix of the import path to use when generating files without a trailing slash (e.g. github.com/bmorton) [$CONSTRUCT_IMPORT_PREFIX]
   --source-path                                        the path where new projects should be created (defaults to using GOPATH) [$GOPATH]
   --help, -h                                           show help
   --version, -v                                        print the version
```


## Configuration via file

Construct will optionally load configuration values from a `~/.construct/config` file if it exists.  The file must be formatted as YAML and the key for each value is equal to the environment variable that would be set.  Here's what an example config would look like:

```yaml
CONSTRUCT_TEMPLATE: https://github.com/bmorton/go-template
CONSTRUCT_IMPORT_PREFIX: github.com/bmorton
```
