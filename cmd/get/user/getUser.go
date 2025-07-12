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

func GetUserCmd() *cobra.Command {
	var err error
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
			err = GetUserDetailsWithFormat(name, outputFormat)
			cobra.CheckErr(err)
		},
	}

	c.Flags().StringVarP(&outputFormat, "output", "o", "json", "Output format (json|yaml)")
	if outputFormat != "json" && outputFormat != "yaml" {
		cobra.CheckErr(fmt.Errorf("invalid output format: %s. Supported formats: json, yaml", outputFormat))
	}

	return c
}

func GetUserDetailsWithFormat(name, outputFormat string) error {
	var userDetailsMarshaled []byte

	userDetails, err := core.GetUserDetails(name)
	if err != nil {
		return err
	}

	if outputFormat == "json" {
		userDetailsMarshaled, err = json.MarshalIndent(userDetails, "", "  ")
		if err != nil {
			return err
		}
	}

	if outputFormat == "yaml" {
		userDetailsMarshaled, err = yaml.Marshal(userDetails)
		if err != nil {
			return err
		}
	}

	logger.LogObjectDetails(userDetailsMarshaled)

	return nil
}
