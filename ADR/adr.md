# Architectural Decision Records

|  Architectural Decision | Date  |  Requestor | Approver   | Notes  |
|---|---|---|---|---|
| Use Go to provide cross compiled binaries |  Oct 23 , 2022 | Team  | Team  |  Alternatives were Python scripts OR Bash scripts which could require further dependencies for customer to use  |
| Use Cobra for CLI  |  Oct 23 , 2022 | Team  | Team  |  doormat-cli and tfx use it, it's good enough for us! |
| Use Viper to handle config  | Oct 23 , 2022  | Team  | Team  |   |
| Users should be able to provide configuration via Environment variables OR config file  | Oct 23 , 2022  | Team  | Team  |   |
| Users should be able to re-run the same command to ensure any changes from source is copied/updates in the destination  | Oct 23 , 2022  | Team  | Team  |   |
|   |   |   |   |   |
|   |   |   |   |   |
|   |   |   |   |   |