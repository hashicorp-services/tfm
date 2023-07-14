# TFM Testing
Effective TFM testing includes:
- Migrating objects from a source organization to a destination organization in which resources already exist in the destination organization.
- Migrating objects from a source organization to a destination organization in which resources DO NOT already exist in the destination organization.

### /terraform
The /terraform directory contains terraform code used to create resources in the source and destination TFC orgs

### /configs
The /configs directroy contains the tfm config files used by GitHub actions to perform TFM tests

### /cleanup
The /cleanup directory contains scripts for cleaning up resources

### /state
The /state directory contains the scripts for creating state files