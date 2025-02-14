# List

list sub commands will list out certain resources from the source organization or destination organization.

```sh
# tfm list -h

List objects in an org

Usage:
  tfm list [command]

Available Commands:
  organization     List Organizations
  projects         Projects command
  ssh              ssh-keys command
  teams            Teams command
  vcs              List VCS Providers
  workspace-filter Filter workspaces
  workspaces       Workspaces command

Flags:
  -h, --help          help for list
      --side string   Specify source or destination side to process

Global Flags:
      --config string   Config file, can be used to store common flags, (default is ./.tfm.hcl).

Use "tfm list [command] --help" for more information about a command.
```

## list sub commands

- [`tfm list organizations`](list_orgs.md)
- [`tfm list ssh`](list_ssh.md)
- [`tfm list teams`](list_teams.md)
- [`tfm list vcs`](list_vcs.md)
- [`tfm list projects`](list_projects.md)
- [`tfm list workspaces`](list_workspaces.md)

## Possible Future list command enhancements

- `tfm list agents`

Got an idea for a feature to `tfm`? Submit a [feature request](https://github.com/hashicorp-services/tfm/issues/new?assignees=&labels=&template=feature_request.md&title=)!
