package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/ak"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/core"
)

func createGroupCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "group",
		Short: "Create a group",
		Long: fmt.Sprintf(`Creates a local group in Authentik.

Examples:
  # Create a group
  %s create group example-group`, constants.CmdName),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			_, err := core.CreateGroup(name, []string{}, ak.GroupAttributes{})
			cobra.CheckErr(err)
		},
	}

	return c
}
