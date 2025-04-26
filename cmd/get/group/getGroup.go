package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/core"
	"github.com/svetlyopet/authentik-cli/internal/logger"
	"gopkg.in/yaml.v3"
)

func GetGroupCmd() *cobra.Command {
	var err error
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
			err = GetGroupDetailsWithFormat(name, outputFormat)
			cobra.CheckErr(err)
		},
	}

	c.Flags().StringVarP(&outputFormat, "output", "o", "json", "Output format (json|yaml)")
	if outputFormat != "json" && outputFormat != "yaml" {
		cobra.CheckErr(fmt.Errorf("invalid output format: %s. Supported formats: json, yaml", outputFormat))
	}

	return c
}

func GetGroupDetailsWithFormat(name, outputFormat string) error {
	var groupDetailsMarshaled []byte

	appDetails, err := core.GetGroupDetails(name)
	if err != nil {
		return err
	}

	if outputFormat == "json" {
		groupDetailsMarshaled, err = json.MarshalIndent(appDetails, "", "  ")
		if err != nil {
			return err
		}
	}

	if outputFormat == "yaml" {
		groupDetailsMarshaled, err = yaml.Marshal(appDetails)
		if err != nil {
			return err
		}
	}

	logger.LogObjectDetails(groupDetailsMarshaled)

	return nil
}
