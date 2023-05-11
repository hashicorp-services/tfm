// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package delete

import (
	"github.com/hashicorp-services/tfm/cmd/helper"
	"github.com/hashicorp-services/tfm/tfclient"
)


// Validation function that validates a map is configured correctly in the tfm.hcl file.
// Takes a map's name from the configuration file as a string
func validateMap(c tfclient.ClientContexts, cfgMap string) (bool, map[string]string, error) {
	m, err := helper.ViperStringSliceMap(cfgMap)

	if err != nil {
		o.AddErrorUserProvided3("Error in", cfgMap, "mapping.")
		return false, m, err
	}

	if len(m) <= 0 {
		o.AddErrorUserProvided3("No", cfgMap, "mapping found in configuration file.")
	} else {
		o.AddMessageUserProvided("Using map ", cfgMap)
		o.AddFormattedMessageCalculated("Found %d mappings in the map.", len(m))
		return true, m, nil
	}

	return false, m, nil
}
