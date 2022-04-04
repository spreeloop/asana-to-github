[![coveralls](https://coveralls.io/repos/github/spreeloop/asana-to-github/badge.svg?branch=master)](https://coveralls.io/github/spreeloop/asana-to-github)

# Migrate tasks from Asana to Github

### Setup Go

https://go.dev/doc/install

### Install the CLI

```
go install spreeloop.com/asana-to-github
```

### Usage

```
asana-to-github [flags]

Flags:
      --asana-json string            Relative or absolute path to the exported Asana tasks in JSON format
      --github-organization string   Name of the Github organization
      --github-repository string     Name of the Github repository
      --github-token string          Github personal access token with access to create and modify Projects
  -h, --help                         help for asana-to-github
```
