package user

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/core"
)

var (
	name    string
	surname string
	email   string
	tenant  string
)

var CreateUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Create a user",
	Long: `Creates a local user in Authentik.

Examples:
  # Create a user
  authentik-cli create user example-user --name=example --surname=user --email=example-user@example.com

  # Create a user who will be a tenant admin
  authentik-cli create user example-user --name=example --surname=user --email=example-user@example.com --tenant-admin=example-tenant`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		err := core.CreateUser(username, name, surname, email, tenant)
		cobra.CheckErr(err)
	},
}

func init() {
	var err error

	CreateUserCmd.Flags().StringVar(&name, "name", "", "Given name of the user")
	CreateUserCmd.Flags().StringVar(&surname, "surname", "", "Surname of the user")
	CreateUserCmd.Flags().StringVar(&email, "email", "", "Email associated with the user")
	CreateUserCmd.Flags().StringVar(&tenant, "tenant-admin", "", "Grant admin permissions for a Tenant")

	err = CreateUserCmd.MarkFlagRequired("name")
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	err = CreateUserCmd.MarkFlagRequired("surname")
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	err = CreateUserCmd.MarkFlagRequired("email")
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
}
