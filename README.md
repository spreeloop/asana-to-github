[![coveralls](https://coveralls.io/repos/github/spreeloop/asana-to-github/badge.svg?branch=master)](https://coveralls.io/github/spreeloop/asana-to-github)

# Migrate tasks from Asana to Github

### Setup Go

https://go.dev/doc/install

### Fetch the code

```
git clone https://github.com/spreeloop/asana-to-github.git
```

### Build

```
go build
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
