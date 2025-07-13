package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/core"
)

func getUserCmd() *cobra.Command {
	var outputFormat string

	c := &cobra.Command{
		Use:   "user",
		Short: "Get details about a user",
		Long: fmt.Sprintf(`Get details about a user in Authentik.

Examples:
  # Get details about a user
  %s get user example-user`, constants.CmdName),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			userDetails, err := core.GetUserDetails(name)
			cobra.CheckErr(err)

			err = outputDetailsWithFormat(userDetails, outputFormat)
			cobra.CheckErr(err)
		},
	}

	c.Flags().StringVarP(&outputFormat, "output", "o", "json", "Output format (json|yaml)")
	if outputFormat != "json" && outputFormat != "yaml" {
		cobra.CheckErr(fmt.Errorf("invalid output format: %s. Supported formats: json, yaml", outputFormat))
	}

	return c
}
