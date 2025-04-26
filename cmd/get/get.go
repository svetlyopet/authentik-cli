package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	a "github.com/svetlyopet/authentik-cli/cmd/get/app"
	g "github.com/svetlyopet/authentik-cli/cmd/get/group"
	"github.com/svetlyopet/authentik-cli/internal/constants"
)

func GetCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "get",
		Short: "Get a resource",
		Long: fmt.Sprintf(`Get details about a resource in Authentik.

Examples:
  # Get details about an application
  %s get app example-app`, constants.CmdName),
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Help()
			cobra.CheckErr(err)
		},
	}

	c.AddCommand(a.GetAppCmd())
	c.AddCommand(g.GetGroupCmd())

	return c
}
