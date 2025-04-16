package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/app"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/logger"
	"gopkg.in/yaml.v3"
)

func GetAppCmd() *cobra.Command {
	var err error
	var outputFormat string

	c := &cobra.Command{
		Use:   "app",
		Short: "Get details about an app",
		Long: fmt.Sprintf(`Get details about an application and provider pair in Authentik.

Examples:
  # Get details about an application and pair
  %s get app example-app`, constants.CmdName),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			err = GetAppDetailsWithFormat(name, outputFormat)
			cobra.CheckErr(err)
		},
	}

	c.Flags().StringVarP(&outputFormat, "output", "o", "json", "Output format (json|yaml)")
	if outputFormat != "json" && outputFormat != "yaml" {
		cobra.CheckErr(fmt.Errorf("invalid output format: %s. Supported formats: json, yaml", outputFormat))
	}

	return c
}

func GetAppDetailsWithFormat(name, outputFormat string) error {
	var appDetailsMarshaled []byte

	appDetails, err := app.Get(name)
	if err != nil {
		return err
	}

	if outputFormat == "json" {
		appDetailsMarshaled, err = json.MarshalIndent(appDetails, "", "  ")
		if err != nil {
			return err
		}
	}

	if outputFormat == "yaml" {
		appDetailsMarshaled, err = yaml.Marshal(appDetails)
		if err != nil {
			return err
		}
	}

	logger.LogObjectDetails(appDetailsMarshaled)

	return nil
}
