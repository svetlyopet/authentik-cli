package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/constants"
)

func UpdateCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "update",
		Short: "Update a resource",
		Long: fmt.Sprintf(`Applies a partial update to an existing resource in Authentik.

Only the attributes passed as flags are changed; every other attribute
keeps its current value.

Examples:
  # Change the email of a user
  %s update user example-user --email=new@example.com`, constants.CmdName),
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			cobra.CheckErr(err)
		},
	}

	c.AddCommand(updateUserCmd())
	c.AddCommand(updateGroupCmd())
	c.AddCommand(updateAppCmd())

	return c
}

// flagValue returns a pointer to value when the flag was set on the command
// line, and nil otherwise so that the attribute is left unchanged. Checking
// Changed, rather than comparing with the default, lets a flag be set to its
// zero value, as in --active=false.
func flagValue[T any](cmd *cobra.Command, name string, value T) *T {
	if !cmd.Flags().Changed(name) {
		return nil
	}
	return &value
}
