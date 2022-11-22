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
		Address: "https://" + viper.GetString("sourceHostname"),
		Token:   viper.GetString("sourceToken"),
	}

	sourceClient, err := tfe.NewClient(sourceConfig)
	if err != nil {
		log.Fatal(err)
	}

	destinationConfig := &tfe.Config{
		Address: "https://" + viper.GetString("destinationHostname"),
		Token:   viper.GetString("destinationToken"),
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
		viper.GetString("sourceHostname"),
		viper.GetString("sourceOrganization"),
		viper.GetString("sourceToken"),
		destinationClient,
		destinationCtx,
		viper.GetString("destinationHostname"),
		viper.GetString("destinationOrganization"),
		viper.GetString("destinationToken")}
}


func Foo() string {
	return "Called Foo(), Return with Bar"
}