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
   --help, -h           show help
   --version, -v          print the version
```
