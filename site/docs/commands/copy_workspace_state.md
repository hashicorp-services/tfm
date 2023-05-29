# tfm copy workspaces --state

`tfm copy workspaces --state` or `tfm copy ws --state` copies a workspaces' states from source to destination org.

!!! note ""
    *NOTE: Currently ALL states will be copied over to the destination.  Future update will be the ability to update 1 or choose X number of the latest states. *


![copy_ws_state](../images/copy_ws_state.png)




# tfm copy workspaces --state --last X

`tfm copy workspaces --state --last X` or `tfm copy ws --state --last X` copies the last X number of workspaces' states from source to destination org.

This flag is designed for users who only want to copy the last X number of states from a workspace. 

!!! WARNING ""
    **WARNING: This operation should not be ran more than once**

![copy_ws_state_last_x](../images/copy_ws_state_last_x.png)




