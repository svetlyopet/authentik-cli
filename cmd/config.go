package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/config"
	"github.com/svetlyopet/authentik-cli/internal/constants"
)

func configCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "config",
		Short: fmt.Sprintf("Creates a configuration for the %s", constants.CmdName),
		Long: fmt.Sprintf(`Creates or updates the config file for %s.
Uses the default path which is in your home directory set by $HOME
environment variable`, constants.CmdName),
		Run: func(cmd *cobra.Command, args []string) {
			err := config.Set()
			cobra.CheckErr(err)
		},
	}

	return c
}
