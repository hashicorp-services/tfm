// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tfclient

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/viper"
)

const (
	userAgent = "tfm"
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

// Create the source client and if ther is an error, retry
func createSrcClientWithRetry(sourceConfig *tfe.Config, maxRetries int, initialBackoff time.Duration) (*tfe.Client, error) {
	var SourceClient *tfe.Client
	var err error

	for retry := 0; retry <= maxRetries; retry++ {
		SourceClient, err = tfe.NewClient(sourceConfig)
		if err == nil {
			// Context creation successful, exit retry loop.
			return SourceClient, nil
		}

		// Handle the error (e.g., log it).
		fmt.Printf("Error creating client on attempt %d: %v\n", retry+1, err)

		if retry < maxRetries {
			// Calculate the backoff duration using an exponential strategy.
			backoff := time.Duration(math.Pow(2, float64(retry))) * initialBackoff

			// Sleep for the calculated backoff duration before retrying.
			fmt.Printf("Retrying after sleeping for %s...\n", backoff)
			time.Sleep(backoff)
		}
	}

	return nil, fmt.Errorf("Max retries reached. Last error: %v", err)
}

// Create the destination client and if ther is an error, retry
func createDestClientWithRetry(destinationConfig *tfe.Config, maxRetries int, initialBackoff time.Duration) (*tfe.Client, error) {
	var destinationClient *tfe.Client
	var err error

	for retry := 0; retry <= maxRetries; retry++ {
		destinationClient, err = tfe.NewClient(destinationConfig)
		if err == nil {
			// Context creation successful, exit retry loop.
			return destinationClient, nil
		}

		// Handle the error (e.g., log it).
		fmt.Printf("Error creating client on attempt %d: %v\n", retry+1, err)

		if retry < maxRetries {
			// Calculate the backoff duration using an exponential strategy.
			backoff := time.Duration(math.Pow(2, float64(retry))) * initialBackoff

			// Sleep for the calculated backoff duration before retrying.
			fmt.Printf("Retrying after sleeping for %s...\n", backoff)
			time.Sleep(backoff)
		}
	}

	return nil, fmt.Errorf("Max retries reached. Last error: %v", err)
}

func GetClientContexts() ClientContexts {

	maxRetries := 5                   // Maximum number of retries. Used in instances where API rate limiting or network connectivity is less than ideal.
	initialBackoff := 2 * time.Second // Initial backoff duration. Used in instances where API rate limiting or network connectivity is less than ideal.

	sourceConfig := &tfe.Config{
		Address:           "https://" + viper.GetString("src_tfe_hostname"),
		Token:             viper.GetString("src_tfe_token"),
		RetryServerErrors: true,
		Headers:           make(http.Header),
		RetryLogHook: func(attemptNum int, resp *http.Response) {
		},
	}
	sourceConfig.Headers.Set("User-Agent", userAgent)

	sourceClient, err := createSrcClientWithRetry(sourceConfig, maxRetries, initialBackoff)
	if err != nil {
		println("There was an issue creating the source client connection.")
		log.Fatal(err)
	}

	destinationConfig := &tfe.Config{
		Address:           "https://" + viper.GetString("dst_tfc_hostname"),
		Token:             viper.GetString("dst_tfc_token"),
		RetryServerErrors: true,
		Headers:           make(http.Header),
		RetryLogHook: func(attemptNum int, resp *http.Response) {
		},
	}
	destinationConfig.Headers.Set("User-Agent", userAgent)

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
