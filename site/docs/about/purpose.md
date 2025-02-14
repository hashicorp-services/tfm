# Why TFM?

If you are asking why does this CLI exist, read on...

HashiCorp Implementation Services have been helping customers get started with Terraform Enterprise in the last 5-6 years. As some of these customers have matured and situations have changed, some customers are asking how they would migrate to Terraform Cloud.

A small group of individuals have already gone through some skunk work engagements for existing customers to migrate from TFE to TFC, however it was very clunky, slow and very custom. At the same time the experience and knowledged gained have insightful on what was needed if a tool was made to assist migrations of this very nature.

There have been multiple examples in the community of a tool or scripts to assist with migration, however we wanted to provide a CLI binary that could be used in our initial migration services offering as well as left with the customer to continue any future migrations.

Our migration services program would :

- teach the customer how to migrate
- plan how they can migrate from TFE to TFC
- assist with a few migrations using our tool
- then allow them to continue migrations in futures for any workspaces that required more time and planning.

We also had aspirations that this tool could be repeatably used by us or a customer in a CI pipeline of some sort to ensure or keep track of migrations of workspaces from TFE to TFC.
