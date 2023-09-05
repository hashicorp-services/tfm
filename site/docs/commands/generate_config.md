# Generate

Generate a template config `.tfm.hcl` file that you can then go back and configure. 

The template file will be created in the directory in which you run the `tfm generate config` command.

```
# tfm generate config

generate a .tfm.hcl file template

Usage:
  tfm generate [command]

Available Commands:
  config      config command

Flags:
  -h, --help   help for generate

Global Flags:
      --autoapprove     Auto approve the tfm run. --autoapprove=true . false by default
      --config string   Config file, can be used to store common flags, (default is ./.tfm.hcl).
      --json            Print the output in JSON format

Use "tfm generate [command] --help" for more information about a command.
```

Got an idea for a feature to `tfm`? Submit a [feature request](https://github.com/hashicorp-services/tfm/issues/new?assignees=&labels=&template=feature_request.md&title=)! 