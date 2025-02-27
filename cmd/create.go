package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource from stdin",
	Long: `Creates resources in Authentik which are native,
like applications and providers, or abstractions, like tenants,
which are created by this tool.

Examples:
  # Create a tenant
  authentik-cli create tenant example-tenant`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
