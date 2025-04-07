package cmd

import (
	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/config"
)

func configCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "config",
		Short: "Creates a configuration for the authentik-cli",
		Long: `Creates or updates the config file for authentik-cli.
Uses the default path which is in your home directory set by $HOME
environment variable`,
		Run: func(cmd *cobra.Command, args []string) {
			err := config.Set()
			cobra.CheckErr(err)
		},
	}

	return c
}
