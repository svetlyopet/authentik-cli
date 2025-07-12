package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/core"
)

func getGroupCmd() *cobra.Command {
	var outputFormat string

	c := &cobra.Command{
		Use:   "group",
		Short: "Get details about a group",
		Long: fmt.Sprintf(`Get details about a group in Authentik.

Examples:
  # Get details about a group
  %s get group example-group`, constants.CmdName),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			groupDetails, err := core.GetGroupDetails(name)
			cobra.CheckErr(err)
			outputDetailsWithFormat(groupDetails, outputFormat)
		},
	}

	c.Flags().StringVarP(&outputFormat, "output", "o", "json", "Output format (json|yaml)")
	if outputFormat != "json" && outputFormat != "yaml" {
		cobra.CheckErr(fmt.Errorf("invalid output format: %s. Supported formats: json, yaml", outputFormat))
	}

	return c
}
