# Copy

![tfm-copy](../images/tfe_tfm_tfc.png)

Copy or Migrate, this sub commands takes X from source organization and copies (or migrates) it to the destination organization. 

```
# tfm copy -h

Copy objects from Source Organization to Destination Organization

Usage:
  tfm copy [command]

Available Commands:
  teams       Copy Teams
  varsets     Copy Variable Sets
  workspaces  Copy Workspaces

Flags:
  -h, --help   help for copy

Global Flags:
      --config string   Config file, can be used to store common flags, (default is ./.tfm.hcl).

Use "tfm copy [command] --help" for more information about a command.
```

## Copy sub commands

- [`tfm copy teams`](copy_teams.md)
- [`tfm copy varsets`](copy_varsets.md)
- [`tfm copy workspaces`](copy_workspaces.md)