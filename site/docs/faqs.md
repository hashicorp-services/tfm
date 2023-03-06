# TFM FAQs

## Who is `tfm` developed for?

Engineers/Operators that manage Terraform Enterprise/Cloud organizations that need to perform a migration of workspaces. 


## Can `tfm` perform a TFE to TFE migration?

Yes, we developed `tfm` to utilise the `go-tfe` library which is used for both Terraform Enterprise as well as Terraform Cloud. The following is what is capable

- TFE to TFC (Primary use case)
- TFE to TFE
- TFC to TFC
- TFC to TFE


## What constraints are there with migration?



## Will this work on a very old version of Terraform Enterprise?

In all honesty, we have not tested in anger what versions of `go-tfe` will not work with `tfm`.  Internal HashiCorp engineers do have the ability to spin up an older version of TFE. Let us know if you need help, we have a test-pipeline our github actions/test directory that can help populate TFE. 


## Is `tfm` supported by our HashiCorp Global Support Team?

Currently there is no official support whatsoever for `tfm`. This project was developed purposely built intially to assist Implementation Engineers if a migration project was to occur as we knew a few key customers had been asking for it. 



