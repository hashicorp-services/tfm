# Logging

`tfm` uses structured levelled logging via [go-hclog](https://github.com/hashicorp/go-hclog), following the same convention as the Terraform CLI (`TF_LOG`). Log output is silent by default and does not affect normal command output.

## Log Levels

| Level | Description |
|-------|-------------|
| `TRACE` | Highly detailed internal execution steps |
| `DEBUG` | Useful diagnostic information for developers |
| `INFO` | High-level progress messages |
| `WARN` | Non-fatal issues that may indicate a problem |
| `ERROR` | Failures that caused an operation to stop |
| `OFF` | No log output (default) |

## Environment Variables

### `TFM_LOG`

Set the log level for the current invocation.

```shell
TFM_LOG=DEBUG tfm list workspaces
TFM_LOG=TRACE tfm copy workspaces
TFM_LOG=INFO  tfm copy projects
```

Setting `TFM_LOG=JSON` enables JSON-formatted log output at `TRACE` level, which is suitable for log aggregation pipelines:

```shell
TFM_LOG=JSON tfm copy workspaces 2>tfm.log
```

If an invalid value is supplied, `tfm` prints a warning to stderr and defaults to `OFF`:

```
[WARN] Invalid TFM_LOG value: "BANANA". Defaulting to OFF. Valid levels: TRACE DEBUG INFO WARN ERROR OFF
```

### `TFM_LOG_PATH`

Write log output to a file instead of stderr. The file is created if it does not exist and opened in append mode.

```shell
TFM_LOG=DEBUG TFM_LOG_PATH=/var/log/tfm.log tfm copy workspaces
```

Both variables can be combined:

```shell
export TFM_LOG=DEBUG
export TFM_LOG_PATH=./tfm-debug.log
tfm copy workspaces
```

## CLI Flag

The `--verbose` / `-V` persistent flag enables `INFO`-level logging without setting an environment variable:

```shell
tfm --verbose list workspaces
tfm -V copy projects
```

This is equivalent to `TFM_LOG=INFO`. If `TFM_LOG` is already set to a more verbose level (e.g. `DEBUG`), the flag has no additional effect.

## Level Precedence

```
TFM_LOG (env var)  >  --verbose / -V flag  >  default (OFF)
```

## Examples

### Debugging a failed migration

```shell
TFM_LOG=DEBUG tfm copy workspaces 2>debug.log
```

### CI pipeline with structured logs

```shell
export TFM_LOG=JSON
export TFM_LOG_PATH=/logs/tfm-$(date +%Y%m%d).log
tfm copy workspaces --autoapprove
```

### Verbose output to terminal

```shell
tfm --verbose copy projects
```

## .env File

Log settings can be stored in a `.env` file and loaded with your shell or a tool like [direnv](https://direnv.net/):

```shell
# .env
TFM_LOG=DEBUG
TFM_LOG_PATH=./tfm.log
```

!!! warning
    Never commit `.env` files containing tokens or credentials to version control. See the [Configuration](configuration_file/config_file.md) page for more details.
