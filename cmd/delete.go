/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a resource from stdin",
	Long: `Deletes resources in Authentik which are native,
like applications and providers, or abstractions, like tenants,
which are created by this tool.

Examples:
  # Delete a tenant
  authentik-cli delete tenant example-tenant`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
