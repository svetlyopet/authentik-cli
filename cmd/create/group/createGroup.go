package cmd

import (
	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/core"
)

var CreateGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "Create a group",
	Long: `Creates a local group in Authentik.

Examples:
  # Create a group
  authentik-cli create group example-group`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		_, err := core.CreateGroup(name, []string{}, ak.GroupAttributes{})
		cobra.CheckErr(err)
	},
}
