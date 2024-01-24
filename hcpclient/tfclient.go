// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfclient

import (
	"context"
	"fmt"

	hcpauth "github.com/hashicorp/hcp-sdk-go/auth"
	//"github.com/hashicorp/hcp-sdk-go/auth/workload"

	cloud_iam "github.com/hashicorp/hcp-sdk-go/clients/cloud-iam/stable/2019-12-10/client"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-iam/stable/2019-12-10/client/iam_service"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-iam/stable/2019-12-10/client/service_principals_service"

	cloud_resource_manager "github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/client"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/client/organization_service"
	"github.com/hashicorp/hcp-sdk-go/clients/cloud-resource-manager/stable/2019-12-10/client/project_service"

	hcpConfig "github.com/hashicorp/hcp-sdk-go/config"
	sdk "github.com/hashicorp/hcp-sdk-go/httpclient"

	"github.com/spf13/viper"
)

type HCPClientContexts struct {
	SourceHCPClient       hcpauth.OauthConfig
	SourceHCPContext      context.Context
	SourceHCPClientID     string
	SourceHCPClientSecret string

	DestinationHCPClient       hcpauth.OauthConfig
	DestinationHCPContext      context.Context
	DestinationHCPClientID     string
	DestinationHCPClientSecret string
}

// ClientConfig specifies configuration for the client that interacts with HCP
type ClientConfig struct {
	ClientID      string
	ClientSecret  string
	SourceChannel string
}

// Client is an HCP client capable of making requests on behalf of a service principal
type Client struct {
	Config ClientConfig

	IAM               iam_service.ClientService
	Organization      organization_service.ClientService
	Project           project_service.ClientService
	ServicePrincipals service_principals_service.ClientService
}

// NewClient creates a new Client that is capable of making HCP requests
func NewSrcClient(config ClientConfig) (*Client, error) {
	// Build the HCP Config options
	opts := hcpConfig.WithClientCredentials(config.ClientID, config.ClientSecret)

	// Create the HCP Config
	hcp, err := hcpConfig.NewHCPConfig(opts)
	if err != nil {
		return nil, fmt.Errorf("invalid HCP config: %w", err)
	}

	// Fetch a token to verify that we have valid credentials
	if _, err := hcp.Token(); err != nil {
		return nil, fmt.Errorf("no valid credentials available: %w", err)
	}

	httpClient, err := sdk.New(sdk.Config{
		HCPConfig:     hcp,
		SourceChannel: config.SourceChannel,
	})
	if err != nil {
		return nil, err
	}

	client := &Client{
		Config:            config,
		IAM:               cloud_iam.New(httpClient, nil).IamService,
		Organization:      cloud_resource_manager.New(httpClient, nil).OrganizationService,
		Project:           cloud_resource_manager.New(httpClient, nil).ProjectService,
		ServicePrincipals: cloud_iam.New(httpClient, nil).ServicePrincipalsService,
	}

	return client, nil
}




func GetHCPClientContexts() HCPClientContexts {

	maxRetries := 5                   // Maximum number of retries. Used in instances where API rate limiting or network connectivity is less than ideal.
	initialBackoff := 2 * time.Second // Initial backoff duration. Used in instances where API rate limiting or network connectivity is less than ideal.

	sourceHCPConfig := hcpConfig &tfe.Config{
		Address:           "https://" + viper.GetString("src_tfe_hostname"),
		Token:             viper.GetString("src_tfe_token"),
		RetryServerErrors: true,
		RetryLogHook: func(attemptNum int, resp *http.Response) {
		},
	}

	sourceClient, err := createSrcClientWithRetry(sourceHCPConfig, maxRetries, initialBackoff)
	if err != nil {
		println("There was an issue creating the source client connection.")
		log.Fatal(err)
	}

	destinationConfig := &tfe.Config{
		Address:           "https://" + viper.GetString("dst_tfc_hostname"),
		Token:             viper.GetString("dst_tfc_token"),
		RetryServerErrors: true,
		RetryLogHook: func(attemptNum int, resp *http.Response) {
		},
	}

	destinationClient, err := createDestClientWithRetry(destinationConfig, maxRetries, initialBackoff)
	if err != nil {
		println("There was an issue creating the destination client connection.")
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