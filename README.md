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
      --asana-project-id string       ID of the project whose task we want to export
      --asana-token string            Asana personal access token with access to read project tasks
      --github-organization string    Name of the Github organization
      --github-repository string      Name of the Github repository
      --github-token string           Github personal access token with access to create and modify Projects
      --github-rate-limit-delay int   Delay to apply between github API requests to avoid hitting rate limits (default 10)
      --dry-run                       Print out tasks that will be migrated without actually migrating them (default false)
      --force-update                  Always override the issue in github with the data from asana (default false)
  -h, --help                          help for asana-to-github
```
