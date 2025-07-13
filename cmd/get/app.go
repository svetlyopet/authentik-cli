package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/app"
	"github.com/svetlyopet/authentik-cli/internal/constants"
)

func getAppCmd() *cobra.Command {
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
			appDetails, err := app.Get(name)
			cobra.CheckErr(err)

			err = outputDetailsWithFormat(appDetails, outputFormat)
			cobra.CheckErr(err)
		},
	}

	c.Flags().StringVarP(&outputFormat, "output", "o", "json", "Output format (json|yaml)")
	if outputFormat != "json" && outputFormat != "yaml" {
		cobra.CheckErr(fmt.Errorf("invalid output format: %s. Supported formats: json, yaml", outputFormat))
	}

	return c
}
