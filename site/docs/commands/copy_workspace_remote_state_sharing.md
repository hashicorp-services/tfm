# tfm copy workspaces --remote-state-sharing

`tfm copy workspaces --remote-state-sharing` is used to set the remote state sharing setting on a workspace.

![copy_ws_remote_state_sharing](../images/copy_ws_remote_state_sharing.png)

This flag is designed for users that want to copy the workspace state sharing setting to the workspaces in the destination.

A workspace can be configured to share the state org wide or to selected workspaces. By default a workspace is set to being shared to selected workspaces, with no workspaces added.
This flag will go through each workspace on the source organization, check the state sharing setting, and on the destination workspace configure it to be shared org wide, shared to selected workspaces with no workspaces added or shared to selected workspaces, with adding workspaces with matching names from the source workspace to the destination workspace.

# tfm copy workspaces --remote-state-sharing --consolidate-global

There is one additional flag that can be used with `--remote-state-sharing` which is `--consolidate-global`

![copy_ws_remote_state_sharing_org_consolidation](../images/copy_ws_remote_state_sharing_org_consolidation.png)

This flag is to be used for companies that are doing an organization consolidation with their migration. For example, a company running TFE may have 10 orgs each with 10 workspaces, but in HCP Terraform they will have 1 orginzation with 100 workspaces. In one of the source orgs, they have 1 workspace that is shared with 9 other workspaces. When tfm migrates this workspace state sharing setting, in the destination org, that workspace would then be shared with 99 workspaces instead of just 9.

The `--consolidate-global` flag is designed to prevent this from happening. For workspaces that are shared with no workspaces or shared with a list of workspaces that functionality is unchanged from above. However if a workspace is shared org wide tfm will behave differently and perform the following step:

1. tfm will first set the destination workspace to being shared with selected workspaces only
2. List all workspaces in the source org
3. List all workspaces in the destination org
4. Filter the two lists to only return workspaces that have the same name as workspaces from the source org
5. Set that returned list of workspaces to be shared on the workspace in the destination org.

This will result in the workspaces that are being shared org wide in the source org to only being shared with selected workspaces in the destination org and prevent the oversharing of the state with workspaces that shouldn't be accessing the state.
