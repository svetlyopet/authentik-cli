package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/svetlyopet/authentik-cli/internal/constants"
	"github.com/svetlyopet/authentik-cli/internal/core"
)

func updateUserCmd() *cobra.Command {
	var username string
	var name string
	var surname string
	var email string
	var active bool

	c := &cobra.Command{
		Use:   "user",
		Short: "Update a user",
		Long: fmt.Sprintf(`Updates a user in Authentik. Only the attributes passed
as flags are changed.

Examples:
  # Change the email of a user
  %[1]s update user example-user --email=new@example.com

  # Change the display name of a user
  %[1]s update user example-user --name=John --surname=Doe

  # Deactivate a user
  %[1]s update user example-user --active=false

  # Rename a user
  %[1]s update user example-user --username=new-username`, constants.CmdName),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			changes := core.UserChanges{
				Username: flagValue(cmd, "username", username),
				Name:     flagValue(cmd, "name", name),
				Surname:  flagValue(cmd, "surname", surname),
				Email:    flagValue(cmd, "email", email),
				IsActive: flagValue(cmd, "active", active),
			}
			err := core.UpdateUser(args[0], changes)
			cobra.CheckErr(err)
		},
	}

	c.Flags().StringVar(&username, "username", "", "New username of the user")
	c.Flags().StringVar(&name, "name", "", "Given name of the user")
	c.Flags().StringVar(&surname, "surname", "", "Surname of the user (requires --name)")
	c.Flags().StringVar(&email, "email", "", "Email associated with the user")
	c.Flags().BoolVar(&active, "active", true, "Whether the user is active (--active=false deactivates the user)")

	c.MarkFlagsOneRequired("username", "name", "surname", "email", "active")

	return c
}
