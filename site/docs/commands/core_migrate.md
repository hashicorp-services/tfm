# tfm core migrate

`tfm core migrate` will sequentially run all of the commands required to migrate terraform open source / community edition configurations to TFE/TFC workspace management.

tfm will run the following commands in the following order when the migrate command is used:

`tfm core clone`
`tfm core getstate`
`tfm core create-worksapces`
`tfm core upload-state`
`tfm core link-vcs`

# Flags

`--include remove-backend` will add the `tfm core remove-backend` command to be run last as part of the `tfm core migrate` command.
