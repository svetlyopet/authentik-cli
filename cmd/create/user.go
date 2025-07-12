package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/core"
)

func createUserCmd() *cobra.Command {
	var name string
	var surname string
	var email string
	var tenant string

	var err error

	c := &cobra.Command{
		Use:   "user",
		Short: "Create a user",
		Long: fmt.Sprintf(`Creates a local user in Authentik.

Examples:
  # Create a user
  %s create user example-user --name=example --surname=user --email=example-user@example.com

  # Create a user who will be a tenant admin
  %s create user example-user --name=example --surname=user --email=example-user@example.com --tenant-admin=example-tenant`, constants.CmdName, constants.CmdName),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			err = core.CreateUser(username, name, surname, email, tenant)
			cobra.CheckErr(err)
		},
	}

	c.Flags().StringVar(&name, "name", "", "Given name of the user")
	c.Flags().StringVar(&surname, "surname", "", "Surname of the user")
	c.Flags().StringVar(&email, "email", "", "Email associated with the user")
	c.Flags().StringVar(&tenant, "tenant-admin", "", "Grant admin permissions for a Tenant")

	err = c.MarkFlagRequired("name")
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	err = c.MarkFlagRequired("email")
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	return c
}
