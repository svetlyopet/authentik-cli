package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/app"
	"github.com/svetlyopet/authentik-cli/internal/constants"
)

func deleteAppCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "app",
		Short: "Delete an app",
		Long: fmt.Sprintf(`Deletes an application and provider pair in Authentik.

Examples:
  # Delete an app
  %s delete app example-app`, constants.CmdName),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			err := app.Delete(name)
			cobra.CheckErr(err)
		},
	}

	return c
}
