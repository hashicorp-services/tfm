# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

resource "tfe_ssh_key" "source" {
  provider = tfe.source

  name = "tfm-ci-testing-src"
  key  = <<EOF
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACDCgSWWMT5TOfXZnpHB4uAg1e6+gd06QJASShdH7wHbawAAAKiSKdb5kinW
+QAAAAtzc2gtZWQyNTUxOQAAACDCgSWWMT5TOfXZnpHB4uAg1e6+gd06QJASShdH7wHbaw
AAAEAVisUyUHfpsDucm4wBomapQslHyWUwOAnjJcJcGnP5isKBJZYxPlM59dmekcHi4CDV
7r6B3TpAkBJKF0fvAdtrAAAAIGptY2NvbGx1bUBqbWNjb2xsdW0tQzAyRjcwQVVNRDZSAQ
IDBAU=
-----END OPENSSH PRIVATE KEY-----
EOF
}

resource "tfe_ssh_key" "destination" {
  provider = tfe.destination

  name = "tfm-ci-testing-dest"
  key  = <<EOF
-----BEGIN OPENSSH PRIVATE KEY-----
b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAAAMwAAAAtzc2gtZW
QyNTUxOQAAACDCgSWWMT5TOfXZnpHB4uAg1e6+gd06QJASShdH7wHbawAAAKiSKdb5kinW
+QAAAAtzc2gtZWQyNTUxOQAAACDCgSWWMT5TOfXZnpHB4uAg1e6+gd06QJASShdH7wHbaw
AAAEAVisUyUHfpsDucm4wBomapQslHyWUwOAnjJcJcGnP5isKBJZYxPlM59dmekcHi4CDV
7r6B3TpAkBJKF0fvAdtrAAAAIGptY2NvbGx1bUBqbWNjb2xsdW0tQzAyRjcwQVVNRDZSAQ
IDBAU=
-----END OPENSSH PRIVATE KEY-----
EOF
}