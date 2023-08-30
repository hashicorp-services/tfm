# tfm list projects


`tfm list projects` will list projects by default of the source TFE/TFC instance.

![list_projects](../images/list_projects_src.png)


## `--side` flag
Providing the `--side destination` flag will list projects of the destination TFE/TFC instance.

![list_projects](../images/list_projects_dst.png)

## `--json` flag
Providing the `--json` flag will output the project names and IDs in JSON format to make configuring the tfx configuration file more managable.

![list_projects_json](../images/list_projects_json.png)









!!! warning ""
    Note: Projects is a relatively new feature. If the source or destination TFE API endpoint does not support projects, `tfm` will error out. 
