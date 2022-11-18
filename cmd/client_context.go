package cmd

import (
	"context"
	"log"

	tfe "github.com/hashicorp/go-tfe"
)

// type TfeDiscoverClientContextSource struct {
// 	Client              *tfe.Client
// 	Context             context.Context
// 	srcHostname         string
// 	srcOrganizationName string
// 	srcToken            string
// }

// type TfeDiscoverClientContextDestination struct {
// 	Client               *tfe.Client
// 	Context              context.Context
// 	destHostname         string
// 	destOrganizationName string
// 	destToken            string
// }

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
		Address: "https://" + *ViperString("sourceHostname"),
		Token:   *ViperString("sourceToken"),
	}

	sourceClient, err := tfe.NewClient(sourceConfig)
	if err != nil {
		log.Fatal(err)
	}

	destinationConfig := &tfe.Config{
		Address: "https://" + *ViperString("destinationHostname"),
		Token:   *ViperString("destinationToken"),
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
		*ViperString("sourceHostname"),
		*ViperString("sourceOrganization"),
		*ViperString("sourceToken"),
		destinationClient,
		destinationCtx,
		*ViperString("destinationHostname"),
		*ViperString("destinationOrganization"),
		*ViperString("destinationToken")}
}
