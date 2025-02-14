# Delete

delete sub commands will delete  certain objects from an organization.

!!! warning ""
    DANGER this will delete things!

```
# tfm delete -h

delete objects in an org. DANGER this will delete things!

Usage:
  tfm delete [command]

Available Commands:
  workspace   Workspace command

Flags:
  -h, --help          help for delete
      --side string   Specify source or destination side to process

Global Flags:
      --autoapprove     Auto approve the tfm run. --autoapprove=true . false by default
      --config string   Config file, can be used to store common flags, (default is ./.tfm.hcl).

Use "tfm delete [command] --help" for more information about a command.
```

## delete sub commands

- [`tfm delete workspace`](delete_workspace.md)

## Possible Future list command enhancements

- `tfm delete projects`
- `tfm delete teams`
- `tfm delete variable`

Got an idea for a feature to `tfm`? Submit a [feature request](https://github.com/hashicorp-services/tfm/issues/new?assignees=&labels=&template=feature_request.md&title=)!
