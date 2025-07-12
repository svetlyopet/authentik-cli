package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/logger"
	"gopkg.in/yaml.v3"
)

func GetCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "get",
		Short: "Get a resource",
		Long: fmt.Sprintf(`Get details about a resource in Authentik.

Examples:
  # Get details about an application
  %s get app example-app`, constants.CmdName),
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			cobra.CheckErr(err)
		},
	}

	c.AddCommand(getAppCmd())
	c.AddCommand(getUserCmd())
	c.AddCommand(getGroupCmd())
	c.AddCommand(getTenantCmd())

	return c
}

func outputDetailsWithFormat(details any, outputFormat string) error {
	var detailsMarshaled []byte
	var err error

	if outputFormat == "json" {
		detailsMarshaled, err = json.MarshalIndent(details, "", "  ")
		if err != nil {
			return err
		}
	}

	if outputFormat == "yaml" {
		detailsMarshaled, err = yaml.Marshal(details)
		if err != nil {
			return err
		}
	}

	logger.LogObjectDetails(detailsMarshaled)

	return nil
}
