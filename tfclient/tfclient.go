// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfclient

import (
	"context"
	"log"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/viper"
)

type ClientContexts struct {
	SourceClient                *tfe.Client
	SourceContext               context.Context
	SourceHostname              string
	SourceOrganizationName      string
	SourceToken                 string
	DestinationClient           *tfe.Client
	DestinationContext          context.Context
	DestinationHostname         string
	DestinationOrganizationName string
	DestinationToken            string
}

func GetClientContexts() ClientContexts {

	sourceConfig := &tfe.Config{
		Address: "https://" + viper.GetString("src_tfe_hostname"),
		Token:   viper.GetString("src_tfe_token"),
	}

	sourceClient, err := tfe.NewClient(sourceConfig)
	if err != nil {
		log.Fatal(err)
	}

	destinationConfig := &tfe.Config{
		Address: "https://" + viper.GetString("dst_tfc_hostname"),
		Token:   viper.GetString("dst_tfc_token"),
	}

	destinationClient, err := tfe.NewClient(destinationConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Create a context
	sourceCtx := context.Background()
	destinationCtx := context.Background()

	return ClientContexts{
		sourceClient,
		sourceCtx,
		viper.GetString("src_tfe_hostname"),
		viper.GetString("src_tfe_org"),
		viper.GetString("src_tfe_token"),
		destinationClient,
		destinationCtx,
		viper.GetString("dst_tfc_hostname"),
		viper.GetString("dst_tfc_org"),
		viper.GetString("dst_tfc_token")}
}

func Foo() string {
	return "Called Foo(), Return with Bar"
}

// GetTfcConfig returns a TFE/TFC config with token if found
// in the terraform local cred file
func GetTfcConfig(hostname string) (tfe.Config, error) {
	config := tfe.Config{
		Address: "https://" + viper.GetString("hostname"),
		Token:   viper.GetString("token"),
	}

	return config, nil
}
