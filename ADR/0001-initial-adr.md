# ADR #0001: Initial Architectural Decision Records

|  Architectural Decision | Date  |  Requestor | Approver   | Notes  |
|---|---|---|---|---|
| Use Go to provide cross compiled binaries |  Oct 23 , 2022 | Team  | Team  |  Alternatives were Python scripts OR Bash scripts which could require further dependencies for customer to use  |
| Use Cobra for CLI  |  Oct 23 , 2022 | Team  | Team  |  doormat-cli and tfx use it, it's good enough for us! |
| Use Viper to handle config  | Oct 23 , 2022  | Team  | Team  |   |
| Users should be able to provide configuration via Environment variables OR config file  | Oct 23 , 2022  | Team  | Team  |   |
| Users should be able to re-run the same command to ensure any changes from source is copied/updates in the destination  | Oct 23 , 2022  | Team  | Team  |   |
| Any objects/resources in the destination that already exist, will be skipped.   | Oct 23, 2022  | Team  | Team  | eg If a workspace exists, it will not try to recreate it in the destination  |
| Workspace settings will always be updated in the destination from the source. This allow the tool be re-ran if any changes in the source occurs.  | Nov, 2022  | Team  | Team  |   |
|   |   |   |   |   |