# Project Details



## Repository Layout
Information on how this repository is structured.

### Files

```
├── cmd
│   ├── copy
│   ├── helper
│   ├── list
│   └── root.go
├── main.go
├── output
├── site
├── test
├── tfclient
├── tfm
└── version
```

`main.go` 

- Entry point for the CLI, not much code here

`go.mod` and `go.sum` 

- Go package dependencies

`CHANGELOG.md` 

- Repository Changes by release

`version/version.go` 

- Isolated package to create a struct for Version information

`cmds/` 

- Directory for all commands/subcommands

`cmds/copy` 

- each subcommand is placed into it's own directory package

`cmds/list`

- each subcommand is placed into it's own directory package


`cmds/helper` 

- Package that creates an easy set of functions to cleanly retrieve flag values regardless of how they were set


`tfclient` 

- Package for to setup a `go-tfe` source and destination context to interact with the TFC/TFE APIs. 


`output` 

- Package to assist with outputing information for the user. 

`docs` 

- Directory for documentation about the tool. Powered by MKDocs and hosted on GitHub Pages

`.gihub/worksflow`

- Github Action workflows
    - `main.tml` automated builds
    - `release.yaml` binary release build
    - `unit-test.yml` automated features testing pipeline
    - `docs-deploy.yml` deployment of TFM Docs to GitHub pages
