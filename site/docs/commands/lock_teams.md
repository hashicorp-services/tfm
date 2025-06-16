# tfm lock teams

## Usage

`tfm lock teams`

```
Set teams - excluding owner - permissions to read in organization

Usage:
  tfm lock teams [flags]

Aliases:
  teams

Flags:
  -h, --help   help for teams

Global Flags:
      --autoapprove     Auto approve the tfm run. --autoapprove=true . false by default
      --config string   Config file, can be used to store common flags, (default is ~/.tfm.hcl).
      --json            Print the output in JSON format
      --side string     Specify source or destination side to process
```

### Purpose

After an organization has been migrated, the first step would be ensuring all the workspaces are locked. 

To prevent any additional runs on the source organization, it is recommended locking down permissions on it only to owners. 

### Example

```
tfm lock teams

Using config file: tfm/.tfm.hcl
Locking teams in:  tfe.source.io
Locking team:  accepted-bat-admin
Locked team:  accepted-bat-admin
Locking team:  accepted-bat-power
Locked team:  accepted-bat-power
Locking team:  accepted-bat-reader
Locked team:  accepted-bat-reader
Skipping team:  owners
```