# Shell Completion

`tfm` ships with built-in shell completion for **bash**, **zsh**, **fish**, and **PowerShell**,
powered by [Cobra's completion command](https://cobra.dev/docs/how-to-guides/shell-completion/).
Tab-complete subcommands, flags, and argument values to reduce friction during day-to-day use.

```sh
# tfm completion --help

Generate the autocompletion script for tfm for the specified shell.
See each sub-command's help for details on how to use the generated script.

Usage:
  tfm completion [command]

Available Commands:
  bash        Generate the autocompletion script for bash
  fish        Generate the autocompletion script for fish
  powershell  Generate the autocompletion script for powershell
  zsh         Generate the autocompletion script for zsh

Flags:
  -h, --help   help for completion

Global Flags:
      --config string   Config file, can be used to store common flags, (default is ~/.tfm.hcl).

Use "tfm completion [command] --help" for more information about a command.
```

---

## bash

### One-time session (no restart needed)

```bash
source <(tfm completion bash)
```

### Permanent installation

=== "Linux"

    ```bash
    tfm completion bash | sudo tee /etc/bash_completion.d/tfm > /dev/null
    ```

=== "macOS (Homebrew bash)"

    ```bash
    tfm completion bash > $(brew --prefix)/etc/bash_completion.d/tfm
    ```

Then start a new shell or `source ~/.bashrc`.

---

## zsh

### One-time session

```zsh
source <(tfm completion zsh)
```

### Permanent installation

```zsh
# Ensure completions are enabled in your shell (add to ~/.zshrc if not present):
autoload -Uz compinit && compinit

# Write the script to a directory on your $fpath:
tfm completion zsh > "${fpath[1]}/_tfm"
```

Restart your shell or run `exec zsh`.

---

## fish

```fish
# One-time session
tfm completion fish | source

# Permanent installation
tfm completion fish > ~/.config/fish/completions/tfm.fish
```

---

## PowerShell

```powershell
# One-time session
tfm completion powershell | Out-String | Invoke-Expression

# Permanent installation — add to your PowerShell profile:
tfm completion powershell >> $PROFILE
```

---

## Verifying completion is working

After installing, type `tfm ` and press **Tab** to see the list of top-level subcommands.
For example:

```
$ tfm <Tab>
completion  copy  core  delete  generate  list  lock  unlock
```

Flags are also tab-completed:

```
$ tfm list --<Tab>
--autoapprove  --config  --json  --verbose
```
