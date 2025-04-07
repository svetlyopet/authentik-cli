package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/core"
)

func DeleteGroupCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "group",
		Short: "Delete a group",
		Long: fmt.Sprintf(`Deletes a group in Authentik.

Examples:
  # Delete a group
  %s delete group example-group`, constants.CmdName),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			err := core.DeleteGroup(name)
			cobra.CheckErr(err)
		},
	}

	return c
}
